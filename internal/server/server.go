// Package server is a server package for kratos.
package server

import (
	"embed"
	nethttp "net/http"

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

var ProviderSetServer = wire.NewSet(NewHTTPServer, NewGRPCServer, RegisterService, NewEventBus)

// ProviderSetServerOnly provides only HTTP and gRPC servers without EventBus
var ProviderSetServerOnly = wire.NewSet(NewHTTPServer, NewGRPCServer, RegisterServiceOnly)

// ProviderSetJob provides only EventBus for job workers
var ProviderSetJob = wire.NewSet(NewEventBus)

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

func (s Servers) BindSwagger(bc *conf.Bootstrap, helper *klog.Helper) {
	if bc.GetEnableSwagger() != "true" {
		helper.Debugw("msg", "swagger is not enabled", "enableSwagger", bc.GetEnableSwagger())
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

	// Create file server handler
	authHandler := nethttp.StripPrefix("/doc/", nethttp.FileServer(nethttp.FS(docFS)))
	basicAuth := bc.GetSwaggerBasicAuth()
	if basicAuth.GetEnabled() == "true" {
		authHandler = middler.BasicAuthMiddleware(basicAuth.GetUsername(), basicAuth.GetPassword())(authHandler)
		helper.Debugf("[Swagger] endpoint: %s/doc/swagger (Basic Auth: %s:%s)", endpoint, basicAuth.GetUsername(), basicAuth.GetPassword())
	} else {
		helper.Debugf("[Swagger] endpoint: %s/doc/swagger (No Basic Auth)", endpoint)
	}

	httSrv.HandlePrefix("/doc/", authHandler)
}

func (s Servers) BindMetrics(bc *conf.Bootstrap, helper *klog.Helper) {
	if bc.GetEnableMetrics() != "true" {
		helper.Debugw("msg", "metrics is not enabled", "enableMetrics", bc.GetEnableMetrics())
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
	basicAuth := bc.GetMetricsBasicAuth()
	authHandler := promhttp.Handler()
	if basicAuth != nil && basicAuth.GetEnabled() == "true" {
		authHandler = middler.BasicAuthMiddleware(basicAuth.GetUsername(), basicAuth.GetPassword())(authHandler)
		helper.Debugf("[Metrics] endpoint: %s/metrics (Basic Auth: %s:%s)", endpoint, basicAuth.GetUsername(), basicAuth.GetPassword())
	} else {
		helper.Debugf("[Metrics] endpoint: %s/metrics (No Basic Auth)", endpoint)
	}
	httSrv.Handle("/metrics", authHandler)
}

// RegisterService registers the service.
func RegisterService(
	c *conf.Bootstrap,
	httpSrv *http.Server,
	grpcSrv *grpc.Server,
	eventBusSrv *EventBus,
	healthService *service.HealthService,
	emailService *service.EmailService,
	webhookService *service.WebhookService,
	senderService *service.SenderService,
	namespaceService *service.NamespaceService,
	messageLogService *service.MessageLogService,
	templateService *service.TemplateService,
	eventBusService *service.EventBusService,
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

	eventBusSrv.RegisterHandler(eventBusService)
	return Servers{httpSrv, grpcSrv, eventBusSrv}
}

// RegisterServiceOnly registers only HTTP and gRPC services without EventBus.
func RegisterServiceOnly(
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
