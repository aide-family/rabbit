// Package data is the data package for the Rabbit service.
package data

import (
	"strings"

	"github.com/aide-family/magicbox/safety"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/gorm"

	"github.com/aide-family/rabbit/internal/biz/do/query"
	"github.com/aide-family/rabbit/internal/conf"
	"github.com/aide-family/rabbit/pkg/connect"
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
	helper *klog.Helper
	c      *conf.Bootstrap
	dbs    *safety.SyncMap[string, *gorm.DB]
	mainDB *gorm.DB

	closes map[string]func() error
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

func (d *Data) BizDB(namespace string) *gorm.DB {
	db, ok := d.dbs.Get(namespace)
	if ok {
		return db
	}
	return d.mainDB
}
