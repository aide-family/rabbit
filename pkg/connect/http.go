// Package connect is a package for connecting to services.
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
	"github.com/go-kratos/kratos/v2/transport/http"

	"github.com/aide-family/rabbit/pkg/merr"
)

func InitHTTPClient(c InitConfig, opts ...InitOption) (*http.Client, error) {
	cfg, err := NewInitConfig(c, opts...)
	if err != nil {
		return nil, err
	}
	if !strings.EqualFold(strings.ToUpper(cfg.protocol), ProtocolHTTP) {
		return nil, merr.ErrorInternalServer("protocol is not HTTP, got %s", cfg.protocol)
	}
	middlewares := []middleware.Middleware{
		recovery.Recovery(),
		middler.Validate(),
		metadata.Client(),
	}
	if strutil.IsNotEmpty(cfg.secret) {
		middlewares = append(middlewares, getJwtClientMiddleware(cfg.secret, cfg.claim))
	}

	clientOpts := []http.ClientOption{
		http.WithEndpoint(cfg.endpoint),
		http.WithMiddleware(middlewares...),
	}

	if pointer.IsNotNil(cfg.discovery) {
		clientOpts = append(clientOpts, http.WithDiscovery(cfg.discovery), http.WithBlock())
		nodeVersion := strings.TrimSpace(cfg.nodeVersion)
		if strutil.IsNotEmpty(nodeVersion) {
			nodeFilter := filter.Version(nodeVersion)
			clientOpts = append(clientOpts, http.WithNodeFilter(nodeFilter))
		}
	}

	if cfg.timeout > 0 {
		clientOpts = append(clientOpts, http.WithTimeout(cfg.timeout))
	}

	return http.NewClient(context.Background(), clientOpts...)
}
