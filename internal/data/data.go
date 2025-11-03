// Package data is the data package for the Rabbit service.
package data

import (
	"strings"
	"time"

	"github.com/aide-family/magicbox/pointer"
	"github.com/aide-family/magicbox/safety"
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
		helper: helper,
		c:      c,
		dbs:    safety.NewSyncMap(make(map[string]*gorm.DB)),
		closes: make(map[string]func() error),
	}
	if err := d.initRegistry(); err != nil {
		return nil, d.close, err
	}
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

		for _, namespace := range strings.Split(namespace, ",") {
			d.dbs.Set(namespace, db)
		}

		d.closes["bizDB."+namespace] = func() error { return connect.CloseDB(db) }
	}

	return d, d.close, nil
}

type Data struct {
	helper   *klog.Helper
	c        *conf.Bootstrap
	dbs      *safety.SyncMap[string, *gorm.DB]
	mainDB   *gorm.DB
	registry connect.Registry

	closes map[string]func() error
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

func (d *Data) MainDB() *gorm.DB {
	return d.mainDB
}

func (d *Data) MainQuery() *query.Query {
	return query.Use(d.MainDB())
}

func (d *Data) BizQuery(namespace string) *query.Query {
	return query.Use(d.BizDB(namespace))
}

func (d *Data) BizQueryWithTable(namespace string, tableName string, args ...any) *query.Query {
	return query.Use(d.BizDB(namespace).Table(tableName, args...))
}

func (d *Data) BizDB(namespace string) *gorm.DB {
	db, ok := d.dbs.Get(namespace)
	if ok {
		return db
	}
	return d.mainDB
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
			Endpoints:   etcdConfig.GetEndpoints(),
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
