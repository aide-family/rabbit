// Package grpc is the grpc command for the Rabbit service
package grpc

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cobra"

	"github.com/aide-family/rabbit/cmd"
	"github.com/aide-family/rabbit/internal/conf"
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

func NewCmd(defaultServerConfigBytes []byte) *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "grpc",
		Short: "Run the Rabbit gRPC service only",
		Long:  cmdGRPCLong,
		Annotations: map[string]string{
			"group": cmd.ServiceCommands,
		},
		Run: runGRPCServer,
	}
	var bc conf.Bootstrap
	c := config.New(config.WithSource(
		env.NewSource(),
		conf.NewBytesSource(defaultServerConfigBytes),
	))
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

