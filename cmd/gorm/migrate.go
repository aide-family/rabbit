package main

import (
	"github.com/spf13/cobra"
	"gorm.io/gorm"

	"github.com/aide-family/rabbit/cmd"
	"github.com/aide-family/rabbit/internal/biz/do"
)

const cmdMigrateLong = `Migrate database tables based on GORM model definitions.

The migrate command automatically creates or updates database tables according to the
model definitions in the codebase. It uses GORM's AutoMigrate feature to ensure the
database schema matches the current model structure.

Key Features:
  • Automatic migration: Automatically create or update tables based on model definitions
  • Schema synchronization: Keep database schema in sync with code models
  • Safe operations: Only adds missing columns and indexes, does not delete existing data
  • Multi-table support: Migrate all models defined in the application

Use Cases:
  • Initial database setup: Create all required tables for the first time
  • Schema updates: Update database schema when models are modified
  • Environment synchronization: Ensure database schema consistency across environments

The migration process will create the database if it doesn't exist and migrate all
tables defined in the model registry.`

func newMigrateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Migrate database tables based on GORM model definitions",
		Long:  cmdMigrateLong,
		Annotations: map[string]string{
			"group": cmd.DatabaseCommands,
		},
		Run: func(cmd *cobra.Command, args []string) {
			db, err := initDB()
			if err != nil {
				flags.Helper.Errorw("msg", "init db failed", "error", err)
				return
			}
			migrate(db)
		},
	}
}

func migrate(db *gorm.DB) {
	tables := do.Models()
	flags.Helper.Infow("msg", "migrate database", "tables", tables)
	if err := db.Migrator().AutoMigrate(tables...); err != nil {
		flags.Helper.Errorw("msg", "migrate database failed", "error", err, "tables", tables)
		return
	}
	flags.Helper.Infow("msg", "migrate database success")
}
