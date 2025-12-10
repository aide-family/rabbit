// Package job is the job command for the Rabbit service
package job

import (
	"strings"

	"github.com/aide-family/magicbox/hello"
	"github.com/aide-family/magicbox/load"
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

func runJobServer(_ *cobra.Command, _ []string) {
	flags.GlobalFlags = cmd.GetGlobalFlags()
	flags.applyToBootstrap()
	var bc conf.Bootstrap
	sourceOpts := make([]config.Source, 0, len(flags.configPaths))
	if flags.useEnv {
		sourceOpts = append(sourceOpts, env.NewSource())
	}
	if len(flags.configPaths) > 0 {
		for _, configPath := range flags.configPaths {
			if strutil.IsNotEmpty(configPath) {
				sourceOpts = append(sourceOpts, file.NewSource(load.ExpandHomeDir(strings.TrimSpace(configPath))))
			}
		}
	}
	if len(sourceOpts) > 0 {
		if err := conf.Load(&bc, sourceOpts...); err != nil {
			klog.Errorw("msg", "load config failed", "error", err)
			return
		}
		flags.Bootstrap = &bc
	}

	serverConf := flags.GetServer()
	envOpts := []hello.Option{
		hello.WithVersion(flags.Version),
		hello.WithID(flags.Hostname),
		hello.WithName(serverConf.GetName()),
		hello.WithEnv(flags.Environment.String()),
		hello.WithMetadata(serverConf.GetMetadata()),
	}
	if serverConf.GetUseRandomID() == "true" {
		envOpts = append(envOpts, hello.WithID(strutil.RandomID()))
	}
	hello.SetEnvWithOption(envOpts...)

	helper := klog.NewHelper(klog.With(klog.GetLogger(),
		"cmd", "job",
		"service.name", hello.Name(),
		"service.id", hello.ID(),
		"caller", klog.DefaultCaller,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID()),
	)

	app, cleanup, err := wireApp(flags.Bootstrap, helper)
	if err != nil {
		klog.Errorw("msg", "wireApp failed", "error", err)
		return
	}
	defer cleanup()
	if err := app.Run(); err != nil {
		klog.Errorw("msg", "app run failed", "error", err)
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

	return kratos.New(opts...), nil
}

