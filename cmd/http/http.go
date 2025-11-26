// Package http provides the http-only command for the Rabbit service.
package http

import (
	"github.com/aide-family/magicbox/hello"
	"github.com/aide-family/magicbox/strutil"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/config/file"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/spf13/cobra"

	"github.com/aide-family/rabbit/cmd"
	"github.com/aide-family/rabbit/internal/conf"
	"github.com/aide-family/rabbit/internal/data"
	"github.com/aide-family/rabbit/internal/server"
)

func NewCmd() *cobra.Command {
	httpCmd := &cobra.Command{
		Use:   "http",
		Short: "Run the Rabbit HTTP server only",
		Long:  `Run the Rabbit service with HTTP server only (without gRPC and job workers).`,
		Annotations: map[string]string{
			"group": cmd.ServiceCommands,
		},
		Run: runHTTP,
	}
	flags.addFlags(httpCmd)
	return httpCmd
}

func runHTTP(_ *cobra.Command, _ []string) {
	flags.GlobalFlags = cmd.GetGlobalFlags()
	var bc conf.Bootstrap
	c := config.New(config.WithSource(
		env.NewSource(),
		file.NewSource(flags.configPath),
	))
	if err := c.Load(); err != nil {
		flags.Helper.Errorw("msg", "load config failed", "error", err)
		return
	}

	if err := c.Scan(&bc); err != nil {
		flags.Helper.Errorw("msg", "scan config failed", "error", err)
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
		hello.WithName(serverConf.GetName() + "-http"),
		hello.WithEnv(bc.GetEnvironment().String()),
		hello.WithMetadata(metadata),
	}
	if serverConf.GetUseRandomID() == "true" {
		envOpts = append(envOpts, hello.WithID(strutil.RandomID()))
	}
	hello.SetEnvWithOption(envOpts...)

	helper := klog.NewHelper(klog.With(flags.Helper.Logger(),
		"cmd", "http",
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

func newApp(d *data.Data, srvs server.Servers, bc *conf.Bootstrap, helper *klog.Helper) (*kratos.App, error) {
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

	srvs.BindSwagger(bc, helper)
	srvs.BindMetrics(bc, helper)

	return kratos.New(opts...), nil
}
