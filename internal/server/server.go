// Package server is a server package for kratos.
package server

import (
	nethttp "net/http"

	"github.com/go-kratos/kratos/v2/encoding/json"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/aide-family/rabbit/internal/conf"
	"github.com/aide-family/rabbit/internal/service"
	apiv1 "github.com/aide-family/rabbit/pkg/api/v1"
)

var ProviderSetServer = wire.NewSet(NewHTTPServer, NewGRPCServer, RegisterService)

// init initializes the json.MarshalOptions.
func init() {
	json.MarshalOptions = protojson.MarshalOptions{
		UseEnumNumbers:  true, // Emit enum values as numbers instead of their string representation (default is string).
		UseProtoNames:   true, // Use the field names defined in the proto file as the output field names.
		EmitUnpopulated: true, // Emit fields even if they are unset or empty.
	}
}

type Servers []transport.Server

func (s Servers) Append(servers ...transport.Server) Servers {
	return append(s, servers...)
}

func (s Servers) BindSwagger(enableSwagger bool, helper *klog.Helper) {
	if !enableSwagger {
		return
	}

	httSrv, ok := s[0].(*http.Server)
	if !ok {
		return
	}

	endpoint, err := httSrv.Endpoint()
	if err != nil {
		return
	}
	httSrv.HandlePrefix("/doc/", nethttp.StripPrefix("/doc/", nethttp.FileServer(nethttp.FS(docFS))))
	helper.Infof("[Swagger] endpoint: %s/doc/swagger", endpoint)
}

func (s Servers) BindMetrics(enableMetrics bool, helper *klog.Helper) {
	if !enableMetrics {
		return
	}
	httSrv, ok := s[0].(*http.Server)
	if !ok {
		return
	}

	endpoint, err := httSrv.Endpoint()
	if err != nil {
		return
	}
	httSrv.Handle("/metrics", promhttp.Handler())
	helper.Infof("[Metrics] endpoint: %s/metrics", endpoint)
}

// RegisterService registers the service.
func RegisterService(
	c *conf.Bootstrap,
	httpSrv *http.Server,
	grpcSrv *grpc.Server,
	healthService *service.HealthService,
	emailService *service.EmailService,
	webhookService *service.WebhookService,
	senderService *service.SenderService,
	namespaceService *service.NamespaceService,
	messageLogService *service.MessageLogService,
	templateService *service.TemplateService,
) Servers {
	apiv1.RegisterHealthServer(grpcSrv, healthService)
	apiv1.RegisterEmailServer(grpcSrv, emailService)
	apiv1.RegisterWebhookServer(grpcSrv, webhookService)
	apiv1.RegisterSenderServer(grpcSrv, senderService)
	apiv1.RegisterNamespaceServer(grpcSrv, namespaceService)
	apiv1.RegisterMessageLogServer(grpcSrv, messageLogService)
	apiv1.RegisterTemplateServer(grpcSrv, templateService)

	apiv1.RegisterHealthHTTPServer(httpSrv, healthService)
	apiv1.RegisterEmailHTTPServer(httpSrv, emailService)
	apiv1.RegisterWebhookHTTPServer(httpSrv, webhookService)
	apiv1.RegisterSenderHTTPServer(httpSrv, senderService)
	apiv1.RegisterNamespaceHTTPServer(httpSrv, namespaceService)
	apiv1.RegisterMessageLogHTTPServer(httpSrv, messageLogService)
	apiv1.RegisterTemplateHTTPServer(httpSrv, templateService)
	return Servers{httpSrv, grpcSrv}
}
