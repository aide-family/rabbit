package server

import (
	"github.com/aide-family/magicbox/server/middler"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	"github.com/aide-family/rabbit/internal/conf"
	"github.com/aide-family/rabbit/internal/service"
	rabbitMiddler "github.com/aide-family/rabbit/pkg/middler"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(bc *conf.Bootstrap, namespaceService *service.NamespaceService, helper *klog.Helper) *grpc.Server {
	serverConf := bc.GetServer()
	grpcConf := serverConf.GetGrpc()
	jwtConf := bc.GetJwt()

	selectorNamespaceMiddlewares := []middleware.Middleware{
		rabbitMiddler.MustNamespace(),
		rabbitMiddler.MustNamespaceExist(namespaceService.HasNamespace),
	}
	namespaceMiddleware := selector.Server(selectorNamespaceMiddlewares...).Match(middler.AllowListMatcher(namespaceAllowList...)).Build()
	selectorMustAuthMiddlewares := []middleware.Middleware{
		rabbitMiddler.JwtServe(jwtConf.GetSecret()),
		rabbitMiddler.MustLogin(),
		rabbitMiddler.BindJwtToken(),
		namespaceMiddleware,
	}
	authMiddleware := selector.Server(selectorMustAuthMiddlewares...).Match(middler.AllowListMatcher(authAllowList...)).Build()

	grpcMiddlewares := []middleware.Middleware{
		recovery.Recovery(),
		logging.Server(helper.Logger()),
		tracing.Server(),
		metadata.Server(),
		authMiddleware,
		middler.Validate(),
	}
	opts := []grpc.ServerOption{
		grpc.Middleware(grpcMiddlewares...),
	}
	if network := grpcConf.GetNetwork(); network != "" {
		opts = append(opts, grpc.Network(network))
	}
	if address := grpcConf.GetAddress(); address != "" {
		opts = append(opts, grpc.Address(address))
	}
	if timeout := grpcConf.GetTimeout(); timeout != nil {
		opts = append(opts, grpc.Timeout(timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)

	return srv
}
