package server

import (
	"embed"

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
	"github.com/aide-family/rabbit/internal/service"
	rabbitMiddler "github.com/aide-family/rabbit/pkg/middler"
)

//go:embed swagger
var docFS embed.FS

// NewHTTPServer new an HTTP server.
func NewHTTPServer(bc *conf.Bootstrap, namespaceService *service.NamespaceService, helper *klog.Helper) *http.Server {
	serverConf := bc.GetServer()
	httpConf := serverConf.GetHttp()
	jwtConf := bc.GetJwt()

	selectorNamespaceMiddlewares := []middleware.Middleware{
		rabbitMiddler.MustNamespace(),
		rabbitMiddler.MustNamespaceExist(namespaceService.HasNamespace),
	}
	namespaceMiddleware := selector.Server(selectorNamespaceMiddlewares...).Match(middler.AllowListMatcher(bc.GetNamespaceAllowList()...)).Build()
	selectorMustAuthMiddlewares := []middleware.Middleware{
		rabbitMiddler.JwtServe(jwtConf.GetSecret()),
		rabbitMiddler.MustLogin(),
		rabbitMiddler.BindJwtToken(),
		namespaceMiddleware,
	}
	authMiddleware := selector.Server(selectorMustAuthMiddlewares...).Match(middler.AllowListMatcher(jwtConf.GetAllowList()...)).Build()

	httpMiddlewares := []middleware.Middleware{
		recovery.Recovery(),
		logging.Server(helper.Logger()),
		tracing.Server(),
		metadata.Server(),
		authMiddleware,
		middler.Validate(),
	}

	opts := []http.ServerOption{
		rabbitMiddler.Cors(),
		http.Middleware(httpMiddlewares...),
	}
	if network := httpConf.GetNetwork(); network != "" {
		opts = append(opts, http.Network(network))
	}
	if address := httpConf.GetAddress(); address != "" {
		opts = append(opts, http.Address(address))
	}
	if timeout := httpConf.GetTimeout(); timeout != nil {
		opts = append(opts, http.Timeout(timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)

	return srv
}
