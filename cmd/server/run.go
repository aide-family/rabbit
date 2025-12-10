// Package server is the server command for the Rabbit service
package server

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

const cmdRunLong = `Start the Rabbit messaging service with all services (HTTP, gRPC, and Job).

The server command starts all services together:
  • HTTP service: Provides RESTful API interfaces for message delivery and management
  • gRPC service: Provides high-performance gRPC API interfaces for inter-service communication
  • Job service: Provides asynchronous message processing capabilities via EventBus

Rabbit is a distributed messaging platform built on the Kratos framework, supporting unified
management and delivery of multiple message channels (email, Webhook, SMS, etc.). It implements
multi-tenant isolation through namespaces and supports both file-based and database storage modes
to meet different deployment requirements.

Key Features:
  • Multi-channel messaging: Unified management of email, Webhook, SMS, and other message channels
  • Template-based delivery: Support for message template configuration with dynamic content rendering and reuse
  • Asynchronous processing: Queue-based asynchronous message delivery for improved throughput and reliability
  • Configuration management: Centralized management of channel configurations (email servers, Webhook endpoints, etc.)
  • Multi-tenant isolation: Namespace-based isolation of configurations and data for different businesses or tenants

Use Cases:
  • All-in-one deployment: Deploy all services together for simple deployment scenarios
  • Development and testing: Quick start for development and testing environments
  • Small to medium deployments: Suitable for deployments that don't require service separation

Note: For production environments requiring service separation, consider using the http, grpc, or job
commands to start services independently for better scalability and resource management.

After starting the service, Rabbit will listen on the configured ports:
  • HTTP: Default 0.0.0.0:8080 (configurable via --http-address)
  • gRPC: Default 0.0.0.0:9090 (configurable via --grpc-address)
  • Job: Default 0.0.0.0:9091 (configurable via --job-address)`

func NewCmd(defaultServerConfigBytes []byte) *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "server",
		Short: "Run the Rabbit service",
		Long:  cmdRunLong,
		Annotations: map[string]string{
			"group": cmd.ServiceCommands,
		},
		Run: runServer,
	}
	var bc conf.Bootstrap
	c := config.New(config.WithSource(
		env.NewSource(),
		conf.NewBytesSource(defaultServerConfigBytes),
	) /*config.WithPrintLoadedDebugLog(false)*/)
	if err := c.Load(); err != nil {
		klog.Errorw("msg", "load config failed", "error", err)
		panic(err)
	}

	if err := c.Scan(&bc); err != nil {
		klog.Errorw("msg", "scan config failed", "error", err)
		panic(err)
	}

	flags.addFlags(runCmd, &bc)
	return runCmd
}

func runServer(_ *cobra.Command, _ []string) {
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
		"cmd", "run",
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

	srvs.BindSwagger(bc, helper)
	srvs.BindMetrics(bc, helper)

	// 生成客户端配置
	if err := generateClientConfig(bc, srvs, helper); err != nil {
		helper.Warnw("msg", "generate client config failed", "error", err)
	}

	return kratos.New(opts...), nil
}
