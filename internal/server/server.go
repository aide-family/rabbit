// Package server is a server package for kratos.
package server

import (
	"embed"
	nethttp "net/http"
	"strings"

	"buf.build/go/protoyaml"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/encoding/json"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.yaml.in/yaml/v2"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/aide-family/rabbit/internal/conf"
	"github.com/aide-family/rabbit/internal/service"
	apiv1 "github.com/aide-family/rabbit/pkg/api/v1"
	"github.com/aide-family/rabbit/pkg/middler"
)

//go:embed swagger
var docFS embed.FS

type protoYAMLCodec struct {
	marshalOptions   protoyaml.MarshalOptions
	unmarshalOptions protoyaml.UnmarshalOptions
}

func newProtoYAMLCodec() *protoYAMLCodec {
	return &protoYAMLCodec{
		marshalOptions: protoyaml.MarshalOptions{
			UseProtoNames:   true,
			EmitUnpopulated: false, // 过滤 0 值和空值
			Indent:          2,
		},
		unmarshalOptions: protoyaml.UnmarshalOptions{
			DiscardUnknown: true,
		},
	}
}

// Marshal implements encoding.Codec.
func (c *protoYAMLCodec) Marshal(v any) ([]byte, error) {
	switch m := v.(type) {
	case protoreflect.ProtoMessage:
		return c.marshalOptions.Marshal(m)
	default:
		return yaml.Marshal(m)
	}
}

// Unmarshal implements encoding.Codec.
func (c *protoYAMLCodec) Unmarshal(data []byte, v any) error {
	switch m := v.(type) {
	case protoreflect.ProtoMessage:
		return c.unmarshalOptions.Unmarshal(data, m)
	default:
		return yaml.Unmarshal(data, m)
	}
}

// Name implements encoding.Codec.
func (c *protoYAMLCodec) Name() string {
	return "yaml"
}

var (
	ProviderSetServerAll  = wire.NewSet(NewHTTPServer, NewGRPCServer, RegisterService, NewJob)
	ProviderSetServerHTTP = wire.NewSet(NewHTTPServer, RegisterHTTPService)
	ProviderSetServerGRPC = wire.NewSet(NewGRPCServer, RegisterGRPCService)
	ProviderSetServerJob  = wire.NewSet(NewJob, RegisterJobService)
)

// init initializes the json.MarshalOptions.
func init() {
	json.MarshalOptions = protojson.MarshalOptions{
		// UseEnumNumbers:  true, // Emit enum values as numbers instead of their string representation (default is string).
		UseProtoNames:   true, // Use the field names defined in the proto file as the output field names.
		EmitUnpopulated: true, // Emit fields even if they are unset or empty.
	}
	encoding.RegisterCodec(newProtoYAMLCodec())
}

type Servers []transport.Server

func BindSwagger(httpSrv *http.Server, bc *conf.Bootstrap, helper *klog.Helper) {
	if strings.EqualFold(bc.GetEnableSwagger(), "true") {
		helper.Debug("swagger is not enabled")
		return
	}

	endpoint, err := httpSrv.Endpoint()
	if err != nil {
		helper.Errorw("msg", "get http server endpoint failed", "error", err)
		return
	}

	// Create file server handler
	authHandler := nethttp.StripPrefix("/doc/", nethttp.FileServer(nethttp.FS(docFS)))
	basicAuth := bc.GetSwaggerBasicAuth()
	if strings.EqualFold(basicAuth.GetEnabled(), "true") {
		authHandler = middler.BasicAuthMiddleware(basicAuth.GetUsername(), basicAuth.GetPassword())(authHandler)
		helper.Debugf("[Swagger] endpoint: %s/doc/swagger (Basic Auth: %s:%s)", endpoint, basicAuth.GetUsername(), basicAuth.GetPassword())
	} else {
		helper.Debugf("[Swagger] endpoint: %s/doc/swagger (No Basic Auth)", endpoint)
	}

	httpSrv.HandlePrefix("/doc/", authHandler)
}

func BindMetrics(httpSrv *http.Server, bc *conf.Bootstrap, helper *klog.Helper) {
	if strings.EqualFold(bc.GetEnableMetrics(), "true") {
		helper.Debug("metrics is not enabled")
		return
	}

	endpoint, err := httpSrv.Endpoint()
	if err != nil {
		helper.Errorw("msg", "get http server endpoint failed", "error", err)
		return
	}

	basicAuth := bc.GetMetricsBasicAuth()
	authHandler := promhttp.Handler()
	if strings.EqualFold(basicAuth.GetEnabled(), "true") {
		authHandler = middler.BasicAuthMiddleware(basicAuth.GetUsername(), basicAuth.GetPassword())(authHandler)
		helper.Debugf("[Metrics] endpoint: %s/metrics (Basic Auth: %s:%s)", endpoint, basicAuth.GetUsername(), basicAuth.GetPassword())
	} else {
		helper.Debugf("[Metrics] endpoint: %s/metrics (No Basic Auth)", endpoint)
	}
	httpSrv.Handle("/metrics", authHandler)
}

// RegisterService registers the service.
func RegisterService(
	c *conf.Bootstrap,
	httpSrv *http.Server,
	grpcSrv *grpc.Server,
	jobSrv *Job,
	healthService *service.HealthService,
	emailService *service.EmailService,
	webhookService *service.WebhookService,
	senderService *service.SenderService,
	namespaceService *service.NamespaceService,
	messageLogService *service.MessageLogService,
	templateService *service.TemplateService,
	jobService *service.JobService,
) Servers {
	var srvs Servers

	srvs = append(srvs, RegisterHTTPService(c, httpSrv,
		healthService,
		emailService,
		webhookService,
		senderService,
		namespaceService,
		messageLogService,
		templateService,
	)...)
	srvs = append(srvs, RegisterGRPCService(c, grpcSrv,
		healthService,
		emailService,
		webhookService,
		senderService,
		namespaceService,
		messageLogService,
		templateService,
	)...)
	srvs = append(srvs, RegisterJobService(c, jobSrv,
		jobService,
	)...)
	return srvs
}

// RegisterHTTPService registers only HTTP service.
func RegisterHTTPService(
	c *conf.Bootstrap,
	httpSrv *http.Server,
	healthService *service.HealthService,
	emailService *service.EmailService,
	webhookService *service.WebhookService,
	senderService *service.SenderService,
	namespaceService *service.NamespaceService,
	messageLogService *service.MessageLogService,
	templateService *service.TemplateService,
) Servers {
	apiv1.RegisterHealthHTTPServer(httpSrv, healthService)
	apiv1.RegisterEmailHTTPServer(httpSrv, emailService)
	apiv1.RegisterWebhookHTTPServer(httpSrv, webhookService)
	apiv1.RegisterSenderHTTPServer(httpSrv, senderService)
	apiv1.RegisterNamespaceHTTPServer(httpSrv, namespaceService)
	apiv1.RegisterMessageLogHTTPServer(httpSrv, messageLogService)
	apiv1.RegisterTemplateHTTPServer(httpSrv, templateService)
	return Servers{httpSrv}
}

// RegisterGRPCService registers only gRPC service.
func RegisterGRPCService(
	c *conf.Bootstrap,
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
	return Servers{grpcSrv}
}

// RegisterJobService registers only Job service.
func RegisterJobService(
	c *conf.Bootstrap,
	jobSrv *Job,
	jobService *service.JobService,
) Servers {
	jobSrv.RegisterHandler(jobService)
	return Servers{jobSrv}
}

var namespaceAllowList = []string{
	apiv1.OperationNamespaceCreateNamespace,
	apiv1.OperationNamespaceUpdateNamespace,
	apiv1.OperationNamespaceUpdateNamespaceStatus,
	apiv1.OperationNamespaceDeleteNamespace,
	apiv1.OperationNamespaceGetNamespace,
	apiv1.OperationNamespaceListNamespace,
	apiv1.OperationHealthHealthCheck,
}

var authAllowList = []string{
	apiv1.OperationHealthHealthCheck,
}
