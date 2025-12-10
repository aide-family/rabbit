// Package http is the http command for the Rabbit service
package http

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cobra"

	"github.com/aide-family/rabbit/cmd"
	"github.com/aide-family/rabbit/internal/conf"
)

const cmdHTTPLong = `Start the Rabbit HTTP service only, providing RESTful API interfaces for message delivery and management.

The http command starts only the HTTP server component, which provides:
  • Message sending: Send messages through various channels (email, Webhook, SMS, etc.)
  • Configuration management: Manage channel configurations and templates
  • Message query: Query message logs and delivery status
  • Health check: Service health status monitoring

Key Features:
  • RESTful API: Standard HTTP/REST API interfaces
  • Swagger documentation: Automatic API documentation (if enabled via --enable-swagger)
  • Metrics endpoint: Prometheus metrics endpoint at /metrics (if enabled via --enable-metrics)
  • Multi-tenant support: Namespace-based isolation
  • JWT authentication: Secure API access with JWT tokens

Use Cases:
  • API gateway: Deploy HTTP service separately as an API gateway for external clients
  • Load balancing: Deploy multiple HTTP instances behind a load balancer for horizontal scaling
  • Microservices: Integrate HTTP service into microservices architecture
  • Web applications: Provide HTTP API for web frontend applications

Note: This command only starts the HTTP service. For asynchronous message processing, you need to
start the job service separately using the "rabbit job" command.

After starting the service, Rabbit HTTP will listen on the configured HTTP port (default: 0.0.0.0:8080,
configurable via --http-address) and provide RESTful API interfaces for client access.`

func NewCmd(defaultServerConfigBytes []byte) *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "http",
		Short: "Run the Rabbit HTTP service only",
		Long:  cmdHTTPLong,
		Annotations: map[string]string{
			"group": cmd.ServiceCommands,
		},
		Run: runHTTPServer,
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

