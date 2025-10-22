// Package gorm is the gorm package for the Rabbit service
package gorm

import (
	"database/sql"
	"fmt"
	"slices"

	"github.com/aide-family/magicbox/load"
	"github.com/aide-family/magicbox/strutil"
	"github.com/go-kratos/kratos/v2/log"
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
	var bc conf.Bootstrap
	if strutil.IsNotEmpty(flags.configPath) {
		log.Infow("msg", "load config file", "file", flags.configPath)
		if err := load.Load(flags.configPath, &bc); err != nil {
			log.Errorw("msg", "load config file failed", "error", err)
			return nil, err
		}
	}
	flags.applyToBootstrap(&bc)

	// check mysql database is exists
	connectDSN := flags.connectDSN()
	log.Infow("msg", "check mysql database is exists", "dsn", connectDSN)
	sqlDB, err := sql.Open("mysql", connectDSN)
	if err != nil {
		log.Errorw("msg", "open mysql connection failed", "error", err)
		return nil, err
	}
	defer sqlDB.Close()

	log.Infow("msg", "ping mysql")
	if err := sqlDB.Ping(); err != nil {
		log.Errorw("msg", "ping mysql failed", "error", err)
		return nil, err
	}
	log.Infow("msg", "ping mysql success")
	// show databases
	log.Infow("msg", "show databases")
	rows, err := sqlDB.Query("SHOW DATABASES")
	if err != nil {
		log.Errorw("msg", "query databases failed", "error", err)
		return nil, err
	}
	defer rows.Close()
	databases := make([]string, 0)
	for rows.Next() {
		var database string
		if err := rows.Scan(&database); err != nil {
			log.Errorw("msg", "scan database failed", "error", err)
			return nil, err
		}
		databases = append(databases, database)
	}
	log.Infow("msg", "show databases success", "databases", databases)
	if !slices.Contains(databases, flags.database) {
		// create database
		log.Warnw("msg", "database not exists", "database", flags.database)
		log.Infow("msg", "create database", "database", flags.database)
		_, err := sqlDB.Exec(fmt.Sprintf("CREATE DATABASE %s", flags.database))
		if err != nil {
			log.Errorw("msg", "create database failed", "error", err, "database", flags.database)
			return nil, err
		}
		log.Infow("msg", "create database success", "database", flags.database)
	}

	log.Infow("msg", "open mysql connection", "dsn", flags.databaseDSN())
	db, err := gorm.Open(mysql.Open(flags.databaseDSN()), &gorm.Config{})
	if err != nil {
		log.Errorw("msg", "open mysql connection failed", "error", err)
		return nil, err
	}
	log.Infow("msg", "open mysql connection success")
	return db, nil
}
