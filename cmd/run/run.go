// Package run is the run command for the Rabbit service
package run

import (
	"github.com/aide-family/magicbox/hello"
	"github.com/aide-family/magicbox/load"
	"github.com/aide-family/magicbox/strutil"
	"github.com/go-kratos/kratos/v2"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/spf13/cobra"

	"github.com/aide-family/rabbit/cmd"
	"github.com/aide-family/rabbit/internal/conf"
	"github.com/aide-family/rabbit/internal/data"
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

func runServer(_ *cobra.Command, _ []string) {
	flags.GlobalFlags = cmd.GetGlobalFlags()
	var bc conf.Bootstrap
	if err := load.Load(flags.configPath, &bc); err != nil {
		flags.Helper.Errorw("msg", "load config failed", "error", err)
		return
	}
	flags.applyToBootstrap(&bc)
	serverConf := bc.GetServer()
	metadata := serverConf.GetMetadata()
	metadata["repository"] = "https://github.com/aide-family/rabbit"
	metadata["author"] = "Aide Family"
	metadata["email"] = "1058165620@qq.com"
	envOpts := []hello.Option{
		hello.WithVersion(flags.Version),
		hello.WithID(flags.Hostname),
		hello.WithName(serverConf.GetName()),
		hello.WithEnv(bc.GetEnvironment().String()),
		hello.WithMetadata(metadata),
	}
	if serverConf.GetUseRandomID() {
		envOpts = append(envOpts, hello.WithID(strutil.RandomID()))
	}
	hello.SetEnvWithOption(envOpts...)

	helper := klog.NewHelper(klog.With(flags.Helper.Logger(),
		"cmd", "run",
		"service.name", hello.Name(),
		"service.id", hello.ID(),
		"caller", klog.DefaultCaller,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID()),
	)

	app, cleanup, err := wireApp(&bc, helper)
	if err != nil {
		flags.Helper.Errorw("msg", "wireApp failed", "error", err)
		return
	}
	defer cleanup()
	if err := app.Run(); err != nil {
		flags.Helper.Errorw("msg", "app run failed", "error", err)
		return
	}
}

func newApp(d *data.Data, srvs server.Servers, helper *klog.Helper) (*kratos.App, error) {
	defer hello.Hello()
	opts := []kratos.Option{
		kratos.Logger(helper.Logger()),
		kratos.Server(srvs...),
		kratos.Version(hello.Version()),
		kratos.ID(hello.ID()),
		kratos.Name(hello.Name()),
		kratos.Metadata(hello.Metadata()),
	}

	if registry := d.Registry(); registry != nil {
		opts = append(opts, kratos.Registrar(registry))
	}

	srvs.BindSwagger(flags.enableSwagger, helper)
	srvs.BindMetrics(flags.enableMetrics, helper)
	return kratos.New(opts...), nil
}
