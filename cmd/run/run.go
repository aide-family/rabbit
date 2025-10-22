// Package run is the run command for the Rabbit service
package run

import (
	"time"

	"github.com/aide-family/magicbox/hello"
	"github.com/aide-family/magicbox/load"
	"github.com/aide-family/magicbox/log"
	"github.com/aide-family/magicbox/log/stdio"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cobra"
	clientv3 "go.etcd.io/etcd/client/v3"

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
	if err := load.Load(flags.configPath, &bc); err != nil {
		panic(err)
	}
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
	metadata := serverConf.GetMetadata()
	metadata["repository"] = "https://github.com/aide-family/rabbit"
	metadata["author"] = "Aide Family"
	metadata["email"] = "1058165620@qq.com"
	envOpts := []hello.Option{
		hello.WithVersion(flags.Version),
		// hello.WithID(flags.Hostname),
		hello.WithID(time.Now().Format("20060102150405")),
		hello.WithName(serverConf.GetName()),
		hello.WithEnv(bc.GetEnvironment().String()),
		hello.WithMetadata(metadata),
	}
	hello.SetEnvWithOption(envOpts...)

	etcdConfig := bc.GetEtcd()
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdConfig.GetEndpoints(),
		Username:    etcdConfig.GetUsername(),
		Password:    etcdConfig.GetPassword(),
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	registrar := etcd.New(client)
	opts := []kratos.Option{
		kratos.Logger(helper.Logger()),
		kratos.Server(srvs...),
		kratos.Registrar(registrar),
		kratos.Version(hello.Version()),
		kratos.ID(hello.ID()),
		kratos.Name(hello.Name()),
		kratos.Metadata(metadata),
	}
	srvs.BindSwagger(flags.enableSwagger, helper)
	srvs.BindMetrics(flags.enableMetrics, helper)
	return kratos.New(opts...)
}
