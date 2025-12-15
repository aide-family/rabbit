// Package all is the all command for the Rabbit service
package all

import (
	"github.com/aide-family/magicbox/hello"
	"github.com/go-kratos/kratos/v2"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/spf13/cobra"

	"github.com/aide-family/rabbit/cmd"
	"github.com/aide-family/rabbit/cmd/run"
	"github.com/aide-family/rabbit/internal/conf"
	"github.com/aide-family/rabbit/internal/data"
	"github.com/aide-family/rabbit/internal/server"
)

const cmdAllLong = `Start the Rabbit messaging service with all services (HTTP, gRPC, and Job).

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

func NewCmd() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "all",
		Short: "Run the Rabbit all services",
		Long:  cmdAllLong,
		Annotations: map[string]string{
			"group": cmd.ServiceCommands,
		},
		Run: runAll,
	}

	flags.addFlags(runCmd)
	return runCmd
}

func runAll(_ *cobra.Command, _ []string) {
	flags.applyToBootstrap()

	run.StartServer("all", wireApp)
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

	for _, srv := range srvs {
		if httpSrv, ok := srv.(*http.Server); ok {
			server.BindSwagger(httpSrv, bc, helper)
			server.BindMetrics(httpSrv, bc, helper)
		}
	}

	// 生成客户端配置
	if err := run.GenerateClientConfig(bc, srvs, helper); err != nil {
		helper.Warnw("msg", "generate client config failed", "error", err)
	}

	return kratos.New(opts...), nil
}
