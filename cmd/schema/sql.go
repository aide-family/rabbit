package schema

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	authmodel "github.com/aide-family/magicbox/domain/auth/v1/gormimpl/model"
	namespacemodel "github.com/aide-family/magicbox/domain/namespace/v1/gormimpl/model"
	"github.com/glebarez/sqlite"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/aide-family/rabbit/internal/data/impl/do"
	schematool "github.com/aide-family/rabbit/internal/tool/schema"
)

var (
	onlyDB bool
	debug  bool
	output = "deploy/sql/"

	gormConfig = &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	}
)

// models returns all models in migration order: namespace and auth first (referenced by do models), then do models.
func models() []any {
	return append(append(namespacemodel.Models(), authmodel.Models()...), do.Models()...)
}

func newSQLCmd() *cobra.Command {
	sqlCmd := &cobra.Command{
		Use:   "sql",
		Short: "Dump table schema to a SQL file",
		Long:  "Dump DDL from GORM models. Use subcommand: sqlite (in-memory) or mysql (with --host/--user/--pass).",
	}
	sqlCmd.PersistentFlags().BoolVar(&onlyDB, "only-db", false, "Only generate schema for database")
	sqlCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Debug mode")
	sqlCmd.PersistentFlags().StringVarP(&output, "output", "o", output, "Output SQL file path")
	sqlCmd.AddCommand(newSQLiteCmd(), newMySQLCmd())
	return sqlCmd
}

func newSQLiteCmd() *cobra.Command {
	dsn := "file::memory:?cache=shared"
	cmd := &cobra.Command{
		Use:   "sqlite",
		Short: "Dump schema using in-memory SQLite",
		RunE: func(c *cobra.Command, args []string) error {
			db, err := gorm.Open(sqlite.Open(dsn), gormConfig)
			if err != nil {
				return fmt.Errorf("open sqlite: %w", err)
			}
			if debug {
				db = db.Debug()
			}
			if !onlyDB {
				if err := db.AutoMigrate(models()...); err != nil {
					return fmt.Errorf("migrate models: %w", err)
				}
			}
			filename := "schema-sqlite.sql"
			filePath := filepath.Join(output, filename)
			if err := schematool.DumpSchema(schematool.NewSQLiteSchema(db), filePath); err != nil {
				return err
			}
			klog.Infof("Schema written to %s", filePath)
			return nil
		},
	}
	cmd.Flags().StringVar(&dsn, "dsn", dsn, "SQLite DSN")
	return cmd
}

func newMySQLCmd() *cobra.Command {
	var host, user, pass, port, database string
	cmd := &cobra.Command{
		Use:   "mysql",
		Short: "Dump schema from MySQL (create temp DB and migrate from models; for existing DB use: schema from-db mysql)",
		RunE: func(c *cobra.Command, args []string) error {
			if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
				return fmt.Errorf("create output dir: %w", err)
			}
			var dsn string
			if database == "" {
				dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
					user, url.QueryEscape(pass), host, port)
				conn, err := gorm.Open(mysql.Open(dsn), gormConfig)
				if err != nil {
					return fmt.Errorf("open mysql: %w", err)
				}
				database = "rabbit_schema_tmp"
				if err := conn.Exec("CREATE DATABASE IF NOT EXISTS `" + database + "`").Error; err != nil {
					return fmt.Errorf("create database: %w", err)
				}
				rawDB, err := conn.DB()
				if err != nil {
					return fmt.Errorf("get db: %w", err)
				}
				if err := rawDB.Close(); err != nil {
					return fmt.Errorf("close db: %w", err)
				}
				dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
					user, url.QueryEscape(pass), host, port, database)
			}

			db, err := gorm.Open(mysql.Open(dsn), gormConfig)
			if err != nil {
				return fmt.Errorf("open mysql: %w", err)
			}
			if debug {
				db = db.Debug()
			}
			if !onlyDB {
				if err := db.AutoMigrate(models()...); err != nil {
					return fmt.Errorf("migrate models: %w", err)
				}
			}
			filename := "schema-mysql.sql"
			filePath := filepath.Join(output, filename)
			if err := schematool.DumpSchema(schematool.NewMySQLSchema(db), filePath); err != nil {
				return err
			}
			klog.Infof("Schema written to %s", filePath)
			return nil
		},
	}
	cmd.Flags().StringVar(&host, "host", "localhost", "MySQL host")
	cmd.Flags().StringVar(&user, "user", "root", "MySQL user")
	cmd.Flags().StringVar(&pass, "pass", "123456", "MySQL password")
	cmd.Flags().StringVar(&port, "port", "3306", "MySQL port")
	cmd.Flags().StringVar(&database, "database", "", "MySQL database name")
	return cmd
}
