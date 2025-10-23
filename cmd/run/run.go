// Package run is the run command for the Rabbit service
package run

import (
	"time"

	"github.com/aide-family/magicbox/hello"
	"github.com/aide-family/magicbox/load"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/spf13/cobra"
	clientV3 "go.etcd.io/etcd/client/v3"

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

func newApp(bc *conf.Bootstrap, srvs server.Servers, helper *klog.Helper) (*kratos.App, error) {
	defer hello.Hello()

	etcdConfig := bc.GetEtcd()
	client, err := clientV3.New(clientV3.Config{
		Endpoints:   etcdConfig.GetEndpoints(),
		Username:    etcdConfig.GetUsername(),
		Password:    etcdConfig.GetPassword(),
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		helper.Errorw("msg", "etcd client initialization failed", "error", err)
		return nil, err
	}
	registrar := etcd.New(client)
	opts := []kratos.Option{
		kratos.Logger(helper.Logger()),
		kratos.Server(srvs...),
		kratos.Registrar(registrar),
		kratos.Version(hello.Version()),
		kratos.ID(hello.ID()),
		kratos.Name(hello.Name()),
		kratos.Metadata(hello.Metadata()),
	}
	srvs.BindSwagger(flags.enableSwagger, helper)
	srvs.BindMetrics(flags.enableMetrics, helper)
	return kratos.New(opts...), nil
}
