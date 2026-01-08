// Package gorm is the gorm package for the Rabbit service
package main

import (
	"github.com/aide-family/magicbox/log"
	"github.com/aide-family/magicbox/log/stdio"
	"github.com/aide-family/magicbox/strutil"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/config/file"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cobra"
	"gorm.io/gorm"

	"github.com/aide-family/rabbit/cmd"
	"github.com/aide-family/rabbit/internal/conf"
	"github.com/aide-family/rabbit/pkg/connect"
	"github.com/aide-family/rabbit/pkg/merr"
)

const cmdLong = `GORM code generation and database migration tools for Rabbit service.

The gorm command provides database-related utilities for the Rabbit service, including
automatic code generation for data models and repositories, as well as database schema
migration capabilities.

Key Features:
  • Code generation: Automatically generate GORM query code, models, and repository interfaces
  • Database migration: Automatically migrate database schemas based on model definitions
  • Database management: Support for database creation and connection management
  • Type-safe queries: Generate type-safe query methods for database operations

Subcommands:
  • gen      Generate GORM query code for models and repositories
  • migrate  Migrate database tables based on model definitions

Use Cases:
  • Initial setup: Generate database models and perform initial schema migration
  • Schema updates: Migrate database when model definitions change
  • Code generation: Automatically generate type-safe query code from database models
  • Development workflow: Streamline database operations during development

Use 'rabbit gorm gen' to generate model and repository code, and 'rabbit gorm migrate'
to migrate the database schema.`

func init() {
	logger, err := log.NewLogger(stdio.LoggerDriver())
	if err != nil {
		panic(merr.ErrorInternal("new logger failed with error: %v", err).WithCause(err))
	}
	logger = klog.With(logger,
		"ts", klog.DefaultTimestamp,
	)
	filterLogger := klog.NewFilter(logger, klog.FilterLevel(klog.LevelInfo))
	helper := klog.NewHelper(filterLogger)
	klog.SetLogger(helper.Logger())
}

func NewCmd() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "gorm",
		Short: "GORM code generation and database migration tools",
		Long:  cmdLong,
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

func initDB() (*gorm.DB, func() error, error) {
	flags.GlobalFlags = cmd.GetGlobalFlags()

	var bc conf.Bootstrap
	if strutil.IsNotEmpty(flags.configPath) {
		klog.Debugw("msg", "load config file", "file", flags.configPath)
		c := config.New(config.WithSource(
			env.NewSource(),
			file.NewSource(flags.configPath),
		))
		if err := c.Load(); err != nil {
			klog.Errorw("msg", "load config failed", "error", err)
			return nil, nil, err
		}

		if err := c.Scan(&bc); err != nil {
			klog.Errorw("msg", "scan config failed", "error", err)
			return nil, nil, err
		}
	}
	flags.applyToBootstrap(&bc)

	return connect.NewDB(flags.ormConf, klog.NewHelper(klog.GetLogger()))
}

func main() {
	rootCmd := cmd.NewCmd()
	cmd.Execute(rootCmd, NewCmd())
}
