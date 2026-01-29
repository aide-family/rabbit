package connect

import (
	"fmt"
	"net/url"
	"strings"

	klog "github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/aide-family/magicbox/log/gormlog"
	"github.com/aide-family/magicbox/pointer"
	"github.com/aide-family/rabbit/pkg/config"
	"github.com/aide-family/rabbit/pkg/merr"
)

func init() {
	globalRegistry.RegisterORMFactory(config.ORMConfig_MYSQL, buildORMFromMySQL)
	globalRegistry.RegisterORMFactory(config.ORMConfig_SQLITE, buildORMFromSQLite)
}

// NewDB creates a new database connection.
// If the orm factory is not registered, it will return an error.
// The database connection is not closed, you need to call the returned function to close the connection.
func NewDB(c *config.ORMConfig) (*gorm.DB, func() error, error) {
	factory, ok := globalRegistry.GetORMFactory(c.GetDialector())
	if !ok {
		return nil, nil, merr.ErrorInternalServer("orm factory not registered")
	}
	db, err := factory(c)
	if err != nil {
		return nil, nil, err
	}
	return db, func() error { return closeDBConnection(db) }, nil
}

func closeDBConnection(db *gorm.DB) error {
	mdb, err := db.DB()
	if err != nil {
		return err
	}
	return mdb.Close()
}

func buildORMFromMySQL(c *config.ORMConfig) (*gorm.DB, error) {
	mysqlConf := &config.MySQLOptions{}
	if pointer.IsNotNil(c.GetOptions()) {
		if err := anypb.UnmarshalTo(c.GetOptions(), mysqlConf, proto.UnmarshalOptions{Merge: true}); err != nil {
			return nil, merr.ErrorInternalServer("unmarshal mysql config failed: %v", err)
		}
	}
	params := url.Values{}
	for key, value := range mysqlConf.Parameters {
		params.Add(key, value)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", mysqlConf.Username, mysqlConf.Password, mysqlConf.Host, mysqlConf.Port, mysqlConf.Database, params.Encode())
	ormConfig := &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	}
	if strings.EqualFold(c.GetUseSystemLogger(), "true") {
		ormConfig.Logger = gormlog.New(klog.GetLogger())
	}

	b := &gormBuilder{
		Dialector: mysql.Open(dsn),
		Config:    ormConfig,
		IsDebug:   strings.EqualFold(c.GetDebug(), "true"),
	}
	return b.build()
}

func buildORMFromSQLite(c *config.ORMConfig) (*gorm.DB, error) {
	sqliteConf := &config.SQLiteOptions{}
	if pointer.IsNotNil(c.GetOptions()) {
		if err := anypb.UnmarshalTo(c.GetOptions(), sqliteConf, proto.UnmarshalOptions{Merge: true}); err != nil {
			return nil, merr.ErrorInternalServer("unmarshal sqlite config failed: %v", err)
		}
	}
	ormConfig := &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	}
	if strings.EqualFold(c.GetUseSystemLogger(), "true") {
		ormConfig.Logger = gormlog.New(klog.GetLogger())
	}
	b := &gormBuilder{
		Dialector: sqlite.Open(sqliteConf.Dsn),
		Config:    ormConfig,
		IsDebug:   strings.EqualFold(c.GetDebug(), "true"),
	}
	return b.build()
}

type gormBuilder struct {
	Dialector gorm.Dialector
	Config    *gorm.Config
	IsDebug   bool
}

func (c *gormBuilder) build() (*gorm.DB, error) {
	db, err := gorm.Open(c.Dialector, c.Config)
	if err != nil {
		return nil, err
	}

	if c.IsDebug {
		return db.Debug(), nil
	}

	return db, nil
}
