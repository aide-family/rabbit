// Package http is the http command for the Rabbit service
package http

import (
	"strings"

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

func NewCmd() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "http",
		Short: "Run the Rabbit HTTP service only",
		Long:  cmdHTTPLong,
		Annotations: map[string]string{
			"group": cmd.ServiceCommands,
		},
		Run: runHTTPServer,
	}

	flags.addFlags(runCmd)
	return runCmd
}

func runHTTPServer(_ *cobra.Command, _ []string) {
	flags.applyToBootstrap()
	hello.Hello()
	run.StartServer(strings.Join([]string{flags.Name, flags.Server.Name, "http"}, "."), WireApp)
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
