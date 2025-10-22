package connect

import (
	"context"
	"strings"

	"github.com/aide-family/magicbox/pointer"
	"github.com/aide-family/magicbox/server/middler"
	"github.com/aide-family/magicbox/strutil"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/selector/filter"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	ggrpc "google.golang.org/grpc"

	"github.com/aide-family/rabbit/pkg/merr"
)

func InitGRPCClient(c InitConfig, opts ...InitOption) (*ggrpc.ClientConn, error) {
	cfg, err := NewInitConfig(c, opts...)
	if err != nil {
		return nil, err
	}
	if !strings.EqualFold(strings.ToUpper(cfg.protocol), ProtocolGRPC) {
		return nil, merr.ErrorInternalServer("protocol is not GRPC, got %s", cfg.protocol)
	}
	middlewares := []middleware.Middleware{
		recovery.Recovery(),
		middler.Validate(),
		metadata.Client(),
	}
	if strutil.IsNotEmpty(cfg.secret) {
		middlewares = append(middlewares, getJwtClientMiddleware(cfg.secret, cfg.claim))
	}

	clientOpts := []grpc.ClientOption{
		grpc.WithEndpoint(cfg.endpoint),
		grpc.WithMiddleware(middlewares...),
	}

	if pointer.IsNotNil(cfg.discovery) {
		clientOpts = append(clientOpts, grpc.WithDiscovery(cfg.discovery))
		nodeVersion := strings.TrimSpace(cfg.nodeVersion)
		if nodeVersion != "" {
			nodeFilter := filter.Version(nodeVersion)
			clientOpts = append(clientOpts, grpc.WithNodeFilter(nodeFilter))
		}
	}

	if cfg.timeout > 0 {
		clientOpts = append(clientOpts, grpc.WithTimeout(cfg.timeout))
	}

	return grpc.DialInsecure(context.Background(), clientOpts...)
}
