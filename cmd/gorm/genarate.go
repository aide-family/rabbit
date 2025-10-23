package gorm

import (
	"github.com/spf13/cobra"
	"gorm.io/gen"
	"gorm.io/gorm"

	"github.com/aide-family/rabbit/internal/biz/do"
)

var config = gen.Config{
	OutPath: "./internal/biz/do/query",
	Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	// If you want to generate pointer type properties for nullable fields, set FieldNullable to true
	// FieldNullable: true,
	// If you want to assign default values to fields in the `Create` API, set FieldCoverable to true, see: https://gorm.io/docs/create.html#Default-Values
	FieldCoverable: true,
	// If you want to generate unsigned integer type fields, set FieldSignable to true
	FieldSignable: true,
	// If you want to generate index tags from the database, set FieldWithIndexTag to true
	FieldWithIndexTag: true,
	// If you want to generate type tags from the database, set FieldWithTypeTag to true
	FieldWithTypeTag: true,
	// If you need unit tests for query code, set WithUnitTest to true
	// WithUnitTest: true,
}

func newGenCmd() *cobra.Command {
	genCmd := &cobra.Command{
		Use:   "gen",
		Short: "gorm gen",
		Long:  "gorm generate",
		Run: func(cmd *cobra.Command, args []string) {
			db, err := initDB()
			if err != nil {
				flags.Helper.Errorw("msg", "init db failed", "error", err)
				return
			}
			generate(db)
		},
	}
	genCmd.Flags().StringVarP(&config.OutPath, "out", "o", "./internal/biz/do/query", "output directory")
	return genCmd
}

func generate(db *gorm.DB) {
	g := gen.NewGenerator(config)
	flags.Helper.Infow("msg", "generate code start")
	g.ApplyBasic(do.Models()...)
	g.Execute()
	flags.Helper.Infow("msg", "generate code success")
}
