package connect

import (
	"context"
	"strings"

	"github.com/aide-family/magicbox/pointer"
	"github.com/aide-family/magicbox/server/middler"
	"github.com/aide-family/magicbox/strutil"
	"github.com/aide-family/rabbit/pkg/merr"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/selector/filter"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	jwtv5 "github.com/golang-jwt/jwt/v5"
	ggrpc "google.golang.org/grpc"
)

func InitGRPCClient(c InitConfig, opts ...InitOption) (*ggrpc.ClientConn, error) {
	cfg := NewInitConfig(c, opts...)
	if strings.EqualFold(strings.ToUpper(cfg.protocol), ProtocolGRPC) {
		return nil, merr.ErrorInternalServer("protocol is not GRPC")
	}
	middlewares := []middleware.Middleware{
		recovery.Recovery(),
		middler.Validate(),
		metadata.Client(),
	}
	if strutil.IsNotEmpty(cfg.secret) {
		middlewares = append(middlewares, jwt.Client(func(token *jwtv5.Token) (interface{}, error) {
			return []byte(cfg.secret), nil
		}))
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
