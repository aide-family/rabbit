package main

import (
	"os"

	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cobra"
	"gorm.io/gen"

	"github.com/aide-family/rabbit/cmd"
	"github.com/aide-family/rabbit/internal/biz/do"
)

var genConfig = gen.Config{
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

const cmdGenLong = `Generate GORM query code for models and repositories.

The gen command automatically generates type-safe query code, repository interfaces,
and helper methods based on GORM model definitions. This eliminates the need for
manual query code writing and ensures type safety.

Key Features:
  • Type-safe queries: Generate type-safe query methods for all model operations
  • Repository generation: Automatically generate repository interfaces and implementations
  • Query builder: Generate query builder methods for complex database operations
  • Code customization: Support for custom query methods and field configurations

Configuration:
  • Output path: Specify the output directory for generated code (default: ./internal/biz/do/query)
  • Generation mode: Configure generation options (WithDefaultQuery, WithQueryInterface, etc.)
  • Field options: Control nullable fields, default values, and type tags

Use Cases:
  • Initial code generation: Generate query code for all models at once
  • Model updates: Regenerate code when models are modified
  • Development efficiency: Reduce boilerplate code and improve development speed

The generated code includes query methods, repository interfaces, and type-safe
database operations that can be used throughout the application.`

func newGenCmd() *cobra.Command {
	genCmd := &cobra.Command{
		Use:   "gen",
		Short: "Generate GORM query code for models and repositories",
		Long:  cmdGenLong,
		Annotations: map[string]string{
			"group": cmd.CodeCommands,
		},
		Run: func(cmd *cobra.Command, args []string) {
			generate()
		},
	}
	genCmd.Flags().StringVarP(&genConfig.OutPath, "out", "o", "./internal/biz/do/query", "output directory")
	return genCmd
}

func generate() {
	if flags.forceGen {
		klog.Debugw("msg", "remove all files")
		os.RemoveAll(genConfig.OutPath)
		klog.Debugw("msg", "remove all files success", "path", genConfig.OutPath)
	}
	g := gen.NewGenerator(genConfig)
	g.SetLogger(&genLogger{helper: klog.NewHelper(klog.GetLogger())})
	klog.Debugw("msg", "generate code start")
	g.ApplyBasic(do.Models()...)
	g.Execute()
	klog.Debugw("msg", "generate code success")
}

type genLogger struct {
	helper *klog.Helper
}

func (g *genLogger) Println(v ...any) {
	g.helper.Debug(v...)
}
