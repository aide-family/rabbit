// Package all is the all command for the Rabbit service
package all

import (
	"strings"
	"sync"

	"github.com/aide-family/magicbox/hello"
	"github.com/spf13/cobra"

	"github.com/aide-family/rabbit/cmd"
	"github.com/aide-family/rabbit/cmd/run"
	"github.com/aide-family/rabbit/cmd/run/grpc"
	"github.com/aide-family/rabbit/cmd/run/http"
	"github.com/aide-family/rabbit/cmd/run/job"
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
	hello.Hello()
	wg := new(sync.WaitGroup)
	wg.Go(func() {
		run.StartServer(strings.Join([]string{flags.Name, flags.Server.Name, "http"}, "."), http.WireApp)
	})
	wg.Go(func() {
		run.StartServer(strings.Join([]string{flags.Name, flags.Server.Name, "grpc"}, "."), grpc.WireApp)
	})
	wg.Go(func() {
		run.StartServer(strings.Join([]string{flags.Name, flags.Server.Name, "job"}, "."), job.WireApp)
	})
	wg.Wait()
}
