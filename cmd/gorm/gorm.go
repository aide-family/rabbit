// Package gorm is the gorm package for the Rabbit service
package gorm

import (
	"database/sql"
	"fmt"
	"slices"

	"github.com/aide-family/magicbox/load"
	"github.com/aide-family/magicbox/log/gormlog"
	"github.com/aide-family/magicbox/strutil"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/aide-family/rabbit/cmd"
	"github.com/aide-family/rabbit/internal/conf"
)

func NewCmd() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "gorm",
		Short: "gorm generate and migrate",
		Long:  "gorm generate and migrate, you can use 'rabbit gorm gen' to generate the model and repository, and 'rabbit gorm migrate' to migrate the database",
		Annotations: map[string]string{
			"group": cmd.ServiceCommands,
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	flags.addFlags(runCmd)
	runCmd.AddCommand(
		newGenCmd(),
		newMigrateCmd(),
	)
	return runCmd
}

func initDB() (*gorm.DB, error) {
	flags.Helper = klog.NewHelper(klog.With(flags.Helper.Logger(), "cmd", "gorm"))

	var bc conf.Bootstrap
	if strutil.IsNotEmpty(flags.configPath) {
		flags.Helper.Infow("msg", "load config file", "file", flags.configPath)
		if err := load.Load(flags.configPath, &bc); err != nil {
			flags.Helper.Errorw("msg", "load config file failed", "error", err)
			return nil, err
		}
	}
	flags.applyToBootstrap(&bc)

	// check mysql database is exists
	connectDSN := flags.connectDSN()
	flags.Helper.Infow("msg", "check mysql database is exists", "dsn", connectDSN)
	sqlDB, err := sql.Open("mysql", connectDSN)
	if err != nil {
		flags.Helper.Errorw("msg", "open mysql connection failed", "error", err)
		return nil, err
	}
	defer sqlDB.Close()

	flags.Helper.Infow("msg", "ping mysql")
	if err := sqlDB.Ping(); err != nil {
		flags.Helper.Errorw("msg", "ping mysql failed", "error", err)
		return nil, err
	}
	flags.Helper.Infow("msg", "ping mysql success")
	// show databases
	flags.Helper.Infow("msg", "show databases")
	rows, err := sqlDB.Query("SHOW DATABASES")
	if err != nil {
		flags.Helper.Errorw("msg", "query databases failed", "error", err)
		return nil, err
	}
	defer rows.Close()
	databases := make([]string, 0)
	for rows.Next() {
		var database string
		if err := rows.Scan(&database); err != nil {
			flags.Helper.Errorw("msg", "scan database failed", "error", err)
			return nil, err
		}
		databases = append(databases, database)
	}
	flags.Helper.Infow("msg", "show databases success", "databases", databases)
	if !slices.Contains(databases, flags.database) {
		// create database
		flags.Helper.Warnw("msg", "database not exists", "database", flags.database)
		flags.Helper.Infow("msg", "create database", "database", flags.database)
		_, err := sqlDB.Exec(fmt.Sprintf("CREATE DATABASE %s", flags.database))
		if err != nil {
			flags.Helper.Errorw("msg", "create database failed", "error", err, "database", flags.database)
			return nil, err
		}
		flags.Helper.Infow("msg", "create database success", "database", flags.database)
	}

	flags.Helper.Infow("msg", "open mysql connection", "dsn", flags.databaseDSN())
	db, err := gorm.Open(mysql.Open(flags.databaseDSN()), &gorm.Config{
		Logger: gormlog.New(flags.Helper.Logger()),
	})
	if err != nil {
		flags.Helper.Errorw("msg", "open mysql connection failed", "error", err)
		return nil, err
	}
	flags.Helper.Infow("msg", "open mysql connection success")
	return db.Debug(), nil
}
