package cmd

import (
	"github.com/aide-family/magicbox/hello"
	"github.com/aide-family/magicbox/load"
	"github.com/aide-family/magicbox/log"
	"github.com/aide-family/magicbox/log/stdio"
	"github.com/go-kratos/kratos/v2"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cobra"

	"github.com/aide-family/rabbit/internal/conf"
	"github.com/aide-family/rabbit/internal/server"
)

var (
	configPath    string
	enableSwagger bool
	enableMetrics bool
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the Rabbit service",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	Run: func(cmd *cobra.Command, args []string) {
		var bc conf.Bootstrap
		load.Load(configPath, &bc)

		logger, err := log.NewLogger(stdio.LoggerDriver())
		if err != nil {
			panic(err)
		}
		helper := klog.NewHelper(logger)
		app, cleanup, err := wireApp(&bc, helper)
		if err != nil {
			panic(err)
		}
		defer cleanup()
		if err := app.Run(); err != nil {
			panic(err)
		}
	},
}

func newApp(bc *conf.Bootstrap, srvs server.Servers, helper *klog.Helper) *kratos.App {
	defer hello.Hello()
	serverConf := bc.GetServer()
	envOpts := []hello.Option{
		hello.WithVersion(Version),
		hello.WithID(id),
		hello.WithName(serverConf.GetName()),
		hello.WithEnv(bc.GetEnvironment().String()),
		hello.WithMetadata(serverConf.GetMetadata()),
	}
	hello.SetEnvWithOption(envOpts...)

	opts := []kratos.Option{
		kratos.Logger(helper.Logger()),
		kratos.Server(srvs...),
	}
	srvs.BindSwagger(enableSwagger, helper)
	srvs.BindMetrics(enableMetrics, helper)
	return kratos.New(opts...)
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.
	runCmd.Flags().StringVar(&configPath, "config", "./config", "config file (default is $HOME/config)")
	runCmd.Flags().BoolVar(&enableSwagger, "swagger", false, "enable swagger")
	runCmd.Flags().BoolVar(&enableMetrics, "metrics", false, "enable metrics")
}
