package gorm

import (
	"github.com/aide-family/rabbit/internal/biz/do"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

func newMigrateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "gorm migrate",
		Long:  "gorm migrate",
		Run: func(cmd *cobra.Command, args []string) {
			db, err := initDB()
			if err != nil {
				log.Errorw("msg", "init db failed", "error", err)
				return
			}
			migrate(db)
		},
	}
}

func migrate(db *gorm.DB) {
	tables := do.Models()
	log.Infow("msg", "migrate database", "tables", tables)
	if err := db.Migrator().AutoMigrate(tables...); err != nil {
		log.Errorw("msg", "migrate database failed", "error", err, "tables", tables)
		return
	}
	log.Infow("msg", "migrate database success")
}
