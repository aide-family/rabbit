package server

import (
	"embed"
	nethttp "net/http"

	"github.com/aide-family/magicbox/server/middler"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/http"

	"github.com/aide-family/rabbit/internal/conf"
	rabbitMiddler "github.com/aide-family/rabbit/pkg/middler"
)

//go:embed swagger
var docFS embed.FS

// NewHTTPServer new an HTTP server.
func NewHTTPServer(bc *conf.Bootstrap, helper *klog.Helper) *http.Server {
	serverConf := bc.GetServer()
	httpConf := serverConf.GetHttp()
	jwtConf := bc.GetJwt()

	selectorMiddlewares := []middleware.Middleware{
		rabbitMiddler.JwtServe(jwtConf.GetSecret()),
		rabbitMiddler.MustLogin(),
		rabbitMiddler.BindJwtToken(),
		rabbitMiddler.BindNamespace(),
	}
	authMiddleware := selector.Server(selectorMiddlewares...).Match(middler.AllowListMatcher(jwtConf.GetAllowList()...)).Build()

	opts := []http.ServerOption{
		http.Filter(middler.Cors(&middler.CorsConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{nethttp.MethodGet, nethttp.MethodPost, nethttp.MethodPut, nethttp.MethodDelete, nethttp.MethodOptions},
			MaxAge:       600,
		})),
		http.Middleware(
			recovery.Recovery(),
			logging.Server(helper.Logger()),
			tracing.Server(),
			metadata.Server(),
			authMiddleware,
			middler.Validate(),
		),
	}
	if httpConf.GetNetwork() != "" {
		opts = append(opts, http.Network(httpConf.GetNetwork()))
	}
	if httpConf.GetAddress() != "" {
		opts = append(opts, http.Address(httpConf.GetAddress()))
	}
	if httpConf.GetTimeout() != nil {
		opts = append(opts, http.Timeout(httpConf.GetTimeout().AsDuration()))
	}
	srv := http.NewServer(opts...)

	return srv
}
