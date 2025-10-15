// Package run is the run command for the Rabbit service
package run

import (
	"github.com/aide-family/magicbox/hello"
	"github.com/aide-family/magicbox/load"
	"github.com/aide-family/magicbox/log"
	"github.com/aide-family/magicbox/log/stdio"
	"github.com/go-kratos/kratos/v2"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cobra"

	"github.com/aide-family/rabbit/cmd"
	"github.com/aide-family/rabbit/internal/conf"
	"github.com/aide-family/rabbit/internal/server"
)

func NewCmd() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run the Rabbit service",
		Long:  `A longer description that spans multiple lines and likely contains examples and usage of using your command.`,
		Annotations: map[string]string{
			"group": cmd.ServiceCommands,
		},
		Run: runServer,
	}
	flags.addFlags(runCmd)
	return runCmd
}

func runServer(cmd *cobra.Command, args []string) {
	var bc conf.Bootstrap
	load.Load(flags.configPath, &bc)
	flags.applyToBootstrap(&bc)

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
}

func newApp(bc *conf.Bootstrap, srvs server.Servers, helper *klog.Helper) *kratos.App {
	defer hello.Hello()
	serverConf := bc.GetServer()
	envOpts := []hello.Option{
		hello.WithVersion(flags.Version),
		hello.WithID(flags.Hostname),
		hello.WithName(serverConf.GetName()),
		hello.WithEnv(bc.GetEnvironment().String()),
		hello.WithMetadata(serverConf.GetMetadata()),
	}
	hello.SetEnvWithOption(envOpts...)

	opts := []kratos.Option{
		kratos.Logger(helper.Logger()),
		kratos.Server(srvs...),
	}
	srvs.BindSwagger(flags.enableSwagger, helper)
	srvs.BindMetrics(flags.enableMetrics, helper)
	return kratos.New(opts...)
}
