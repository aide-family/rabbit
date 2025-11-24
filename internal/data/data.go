// Package data is the data package for the Rabbit service.
package data

import (
	"context"
	"time"

	"github.com/aide-family/magicbox/plugin/cache"
	"github.com/aide-family/magicbox/plugin/cache/mem"
	"github.com/aide-family/magicbox/pointer"
	"github.com/aide-family/magicbox/safety"
	"github.com/aide-family/magicbox/strutil"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	kuberegistry "github.com/go-kratos/kratos/contrib/registry/kubernetes/v2"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	clientV3 "go.etcd.io/etcd/client/v3"
	"gorm.io/gorm"

	"github.com/aide-family/rabbit/internal/biz/do/query"
	"github.com/aide-family/rabbit/internal/conf"
	"github.com/aide-family/rabbit/pkg/config"
	"github.com/aide-family/rabbit/pkg/connect"
	"github.com/aide-family/rabbit/pkg/merr"
)

// ProviderSetData is a set of data providers.
var ProviderSetData = wire.NewSet(New)

// New a data and returns.
func New(c *conf.Bootstrap, helper *klog.Helper) (*Data, func(), error) {
	d := &Data{
		helper:      helper,
		c:           c,
		dbs:         safety.NewSyncMap(make(map[string]*gorm.DB)),
		closes:      make(map[string]func() error),
		useDatabase: strutil.IsNotEmpty(c.GetUseDatabase()) && c.GetUseDatabase() == "true",
		reloadFuncs: safety.NewSyncMap(make(map[string]func())),
	}

	if err := d.initRegistry(); err != nil {
		return nil, d.close, err
	}
	if d.useDatabase {
		mainDB, err := connect.NewGorm(d.c.GetMain(), d.helper)
		if err != nil {
			return nil, d.close, err
		}
		d.mainDB = mainDB
		d.closes["mainDB"] = func() error { return connect.CloseDB(mainDB) }

		for namespace, biz := range d.c.GetBiz() {
			db, err := connect.NewGorm(biz, d.helper)
			if err != nil {
				return nil, d.close, err
			}

			for _, ns := range strutil.SplitSkipEmpty(namespace, ",") {
				d.dbs.Set(ns, db)
			}

			d.closes["bizDB.["+namespace+"]"] = func() error { return connect.CloseDB(db) }
		}
	} else {
		if err := d.LoadFileConfig(d.c, d.helper); err != nil {
			return nil, d.close, err
		}
		d.helper.Infow("msg", "file config mode enabled, database connections disabled")
	}

	cacheDriver := mem.CacheDriver()
	cache, err := cache.New(context.Background(), cacheDriver)
	if err != nil {
		return nil, d.close, err
	}
	d.cache = cache
	d.closes["cache"] = func() error { return cache.Close() }

	return d, d.close, nil
}

type Data struct {
	helper      *klog.Helper
	c           *conf.Bootstrap
	dbs         *safety.SyncMap[string, *gorm.DB]
	mainDB      *gorm.DB
	registry    connect.Registry
	cache       cache.Interface
	closes      map[string]func() error
	useDatabase bool
	fileConfig  conf.Config
	reloadFuncs *safety.SyncMap[string, func()]
}

func (d *Data) AppendClose(name string, close func() error) {
	d.closes[name] = close
}

func (d *Data) close() {
	for name, close := range d.closes {
		if err := close(); err != nil {
			d.helper.Errorw("msg", "close db failed", "name", name, "error", err)
			continue
		}
		d.helper.Infow("msg", "close success", "name", name)
	}
}

func (d *Data) UseDatabase() bool {
	return d.useDatabase
}

func (d *Data) MainDB(ctx context.Context) *gorm.DB {
	if !d.UseDatabase() {
		panic(merr.ErrorInternalServer("database mode is not enabled, please check the config"))
	}
	if tx, ok := GetMainTransaction(ctx); ok {
		return tx.DB.WithContext(ctx)
	}
	return d.mainDB.WithContext(ctx)
}

func (d *Data) MainQuery(ctx context.Context) *query.Query {
	return query.Use(d.MainDB(ctx))
}

func (d *Data) BizQuery(ctx context.Context, namespace string) *query.Query {
	return query.Use(d.BizDB(ctx, namespace))
}

func (d *Data) BizQueryWithTable(ctx context.Context, namespace string, tableName string, args ...any) *query.Query {
	return query.Use(d.BizDB(ctx, namespace).Table(tableName, args...))
}

func (d *Data) BizDB(ctx context.Context, namespace string) *gorm.DB {
	if !d.UseDatabase() {
		panic(merr.ErrorInternalServer("database mode is not enabled, please check the config"))
	}
	if tx, ok := GetBizTransaction(ctx, namespace); ok {
		return tx.DB.WithContext(ctx)
	}
	db, ok := d.dbs.Get(namespace)
	if ok {
		return db.WithContext(ctx)
	}
	return d.mainDB.WithContext(ctx)
}

func (d *Data) Registry() connect.Registry {
	return d.registry
}

func (d *Data) initRegistry() error {
	switch registryType := d.c.GetRegistryType(); registryType {
	case config.RegistryType_KUBERNETES:
		kubeConfig := d.c.GetKubernetes()
		if pointer.IsNil(kubeConfig) {
			return merr.ErrorInternalServer("kubernetes config is not found")
		}
		kubeClient, err := connect.NewKubernetesClientSet(kubeConfig.GetKubeConfig())
		if err != nil {
			d.helper.Errorw("msg", "kubernetes client initialization failed", "error", err)
			return err
		}
		registrar := kuberegistry.NewRegistry(kubeClient, kubeConfig.GetNamespace())
		d.registry = registrar
	case config.RegistryType_ETCD:
		etcdConfig := d.c.GetEtcd()
		if pointer.IsNil(etcdConfig) {
			return merr.ErrorInternalServer("etcd config is not found")
		}
		client, err := clientV3.New(clientV3.Config{
			Endpoints:   strutil.SplitSkipEmpty(etcdConfig.GetEndpoints(), ","),
			Username:    etcdConfig.GetUsername(),
			Password:    etcdConfig.GetPassword(),
			DialTimeout: 5 * time.Second,
		})
		if err != nil {
			d.helper.Errorw("msg", "etcd client initialization failed", "error", err)
			return err
		}
		registrar := etcd.New(client)
		d.registry = registrar
		d.closes["etcdClient"] = func() error { return client.Close() }
	}
	return nil
}
