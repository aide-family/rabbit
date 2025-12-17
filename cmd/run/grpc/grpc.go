// Package grpc is the grpc command for the Rabbit service
package grpc

import (
	"strings"

	"github.com/aide-family/magicbox/hello"
	"github.com/go-kratos/kratos/v2"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cobra"

	"github.com/aide-family/rabbit/cmd"
	"github.com/aide-family/rabbit/cmd/run"
	"github.com/aide-family/rabbit/internal/conf"
	"github.com/aide-family/rabbit/internal/data"
	"github.com/aide-family/rabbit/internal/server"
)

const cmdGRPCLong = `Start the Rabbit gRPC service only, providing high-performance gRPC API interfaces for message delivery and management.

The grpc command starts only the gRPC server component, which provides:
  • Message sending: Send messages through various channels (email, Webhook, SMS, etc.)
  • Configuration management: Manage channel configurations and templates
  • Message query: Query message logs and delivery status
  • Health check: Service health status monitoring

Key Features:
  • gRPC API: High-performance gRPC API interfaces with Protocol Buffers
  • Efficient serialization: Binary Protocol Buffers for efficient data transfer
  • Streaming support: Support for streaming RPC calls for real-time data
  • Multi-tenant support: Namespace-based isolation
  • JWT authentication: Secure API access with JWT tokens

Use Cases:
  • Microservices communication: Deploy gRPC service for inter-service communication in microservices architecture
  • High-performance scenarios: Use gRPC for high-throughput message delivery with low latency
  • Service mesh: Integrate gRPC service into service mesh architecture (Istio, Linkerd, etc.)
  • Internal services: Provide gRPC API for internal service-to-service communication

Note: This command only starts the gRPC service. For asynchronous message processing, you need to
start the job service separately using the "rabbit job" command.

After starting the service, Rabbit gRPC will listen on the configured gRPC port (default: 0.0.0.0:9090,
configurable via --grpc-address) and provide gRPC API interfaces for client access.`

func NewCmd() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "grpc",
		Short: "Run the Rabbit gRPC service only",
		Long:  cmdGRPCLong,
		Annotations: map[string]string{
			"group": cmd.ServiceCommands,
		},
		Run: runGRPCServer,
	}

	flags.addFlags(runCmd)
	return runCmd
}

func runGRPCServer(_ *cobra.Command, _ []string) {
	if err := flags.applyToBootstrap(); err != nil {
		klog.Errorw("msg", "apply to bootstrap failed", "error", err)
		return
	}
	hello.Hello()
	run.StartServer(strings.Join([]string{flags.Name, flags.Server.Name, "grpc"}, "."), WireApp)
}

func newApp(serviceName string, d *data.Data, srvs server.Servers, bc *conf.Bootstrap, helper *klog.Helper) (*kratos.App, error) {
	opts := []kratos.Option{
		kratos.Name(serviceName),
		kratos.ID(hello.ID()),
		kratos.Version(hello.Version()),
		kratos.Metadata(hello.Metadata()),
		kratos.Logger(helper.Logger()),
		kratos.Server(srvs...),
	}

	if registry := d.Registry(); registry != nil {
		opts = append(opts, kratos.Registrar(registry))
	}

	// 生成客户端配置
	if err := run.GenerateClientConfig(bc, srvs, helper); err != nil {
		helper.Warnw("msg", "generate client config failed", "error", err)
	}

	return kratos.New(opts...), nil
}
