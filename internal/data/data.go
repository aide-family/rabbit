// Package data is the data package for the Rabbit service.
package data

import (
	"context"

	"github.com/aide-family/magicbox/plugin/cache"
	"github.com/aide-family/magicbox/plugin/cache/mem"
	"github.com/aide-family/magicbox/safety"
	"github.com/aide-family/magicbox/strutil"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
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
		closes:      safety.NewSyncMap(make(map[string]func() error)),
		useDatabase: c.GetMain().GetDialector() == config.ORMConfig_TYPE_UNKNOWN,
		reloadFuncs: safety.NewSyncMap(make(map[string]func())),
	}

	if err := d.initRegistry(); err != nil {
		return nil, d.close, err
	}
	if d.useDatabase {
		mainDB, closeDB, err := connect.NewDB(d.c.GetMain(), d.helper)
		if err != nil {
			return nil, d.close, err
		}
		d.mainDB = mainDB
		d.closes.Set("mainDB", closeDB)

		for namespace, biz := range d.c.GetBiz() {
			db, closeDB, err := connect.NewDB(biz, d.helper)
			if err != nil {
				return nil, d.close, err
			}

			for _, ns := range strutil.SplitSkipEmpty(namespace, ",") {
				d.dbs.Set(ns, db)
			}

			// 使用局部变量避免闭包捕获问题
			namespaceKey := "bizDB.[" + namespace + "]"
			d.closes.Set(namespaceKey, closeDB)
		}
	} else {
		if err := d.LoadFileConfig(d.c, d.helper); err != nil {
			return nil, d.close, err
		}
		d.helper.Debugw("msg", "file config mode enabled, database connections disabled")
	}

	cacheDriver := mem.CacheDriver()
	cache, err := cache.New(context.Background(), cacheDriver)
	if err != nil {
		return nil, d.close, err
	}
	d.cache = cache
	d.closes.Set("cache", func() error { return cache.Close() })

	return d, d.close, nil
}

type Data struct {
	helper      *klog.Helper
	c           *conf.Bootstrap
	dbs         *safety.SyncMap[string, *gorm.DB]
	mainDB      *gorm.DB
	registry    connect.Report
	cache       cache.Interface
	closes      *safety.SyncMap[string, func() error] // 使用SyncMap保证并发安全
	useDatabase bool
	fileConfig  conf.Config
	reloadFuncs *safety.SyncMap[string, func()]
}

func (d *Data) AppendClose(name string, close func() error) {
	d.closes.Set(name, close)
}

func (d *Data) close() {
	d.closes.Range(func(name string, close func() error) bool {
		if err := close(); err != nil {
			d.helper.Errorw("msg", "close db failed", "name", name, "error", err)
			return true // 继续遍历
		}
		d.helper.Debugw("msg", "close success", "name", name)
		return true // 继续遍历
	})
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

func (d *Data) Registry() connect.Report {
	return d.registry
}

func (d *Data) initRegistry() error {
	registry, closeRegistry, err := connect.NewReport(d.c.GetReport(), d.helper)
	if err != nil {
		return err
	}
	d.registry = registry
	d.closes.Set("registry", closeRegistry)
	return nil
}
