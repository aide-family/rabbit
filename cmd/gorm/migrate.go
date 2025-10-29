package main

import (
	"github.com/spf13/cobra"
	"gorm.io/gorm"

	"github.com/aide-family/rabbit/internal/biz/do"
)

func newMigrateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "gorm migrate",
		Long:  "gorm migrate",
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
