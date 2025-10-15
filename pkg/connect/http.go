// Package connect is a package for connecting to services.
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
	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/selector/filter"
	"github.com/go-kratos/kratos/v2/selector/wrr"
	"github.com/go-kratos/kratos/v2/transport/http"
	jwtv5 "github.com/golang-jwt/jwt/v5"
)

func InitHTTPClient(c InitConfig, opts ...InitOption) (*http.Client, error) {
	cfg := NewInitConfig(c, opts...)
	if strings.EqualFold(strings.ToUpper(cfg.protocol), ProtocolHTTP) {
		return nil, merr.ErrorInternalServer("network is not http")
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

	clientOpts := []http.ClientOption{
		http.WithEndpoint(cfg.endpoint),
		http.WithMiddleware(middlewares...),
	}

	if pointer.IsNotNil(cfg.discovery) {
		clientOpts = append(clientOpts, http.WithDiscovery(cfg.discovery))
		nodeVersion := strings.TrimSpace(cfg.nodeVersion)
		if strutil.IsNotEmpty(nodeVersion) {
			nodeFilter := filter.Version(nodeVersion)
			selector.SetGlobalSelector(wrr.NewBuilder())
			clientOpts = append(clientOpts, http.WithNodeFilter(nodeFilter))
		}
	}

	if cfg.timeout > 0 {
		clientOpts = append(clientOpts, http.WithTimeout(cfg.timeout))
	}

	return http.NewClient(context.Background(), clientOpts...)
}
