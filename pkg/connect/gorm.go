package connect

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	klog "github.com/go-kratos/kratos/v2/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/aide-family/magicbox/log/gormlog"
	"github.com/aide-family/rabbit/pkg/config"
)

func NewGorm(mysqlConf *config.MySQLConfig, logger *klog.Helper) (*gorm.DB, error) {
	params := url.Values{}
	for key, value := range mysqlConf.Parameters {
		params.Add(key, value)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", mysqlConf.Username, mysqlConf.Password, mysqlConf.Host, mysqlConf.Port, mysqlConf.Database, params.Encode())
	gormConfig := &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	}
	if strings.EqualFold(mysqlConf.UseSystemLogger, "true") {
		gormConfig.Logger = gormlog.New(logger.Logger())
	}
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("open mysql connection failed: %w, dsn: %s", err, dsn)
	}

	// 配置连接池参数
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get underlying sql.DB failed: %w", err)
	}
	sqlDB.SetMaxIdleConns(20)                  // 最大空闲连接数
	sqlDB.SetConnMaxLifetime(time.Hour)        // 连接最大生存时间
	sqlDB.SetConnMaxIdleTime(10 * time.Minute) // 空闲连接最大生存时间

	if strings.EqualFold(mysqlConf.Debug, "true") {
		db = db.Debug()
	}

	return db, nil
}

func CloseDB(db *gorm.DB) error {
	mdb, err := db.DB()
	if err != nil {
		return fmt.Errorf("get db connection failed: %w", err)
	}
	return mdb.Close()
}
