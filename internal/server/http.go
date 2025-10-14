package server

import (
	"embed"

	"github.com/aide-family/rabbit/internal/conf"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/http"
)

//go:embed swagger
var docFS embed.FS

// NewHTTPServer new an HTTP server.
func NewHTTPServer(bc *conf.Bootstrap, helper *klog.Helper) *http.Server {
	serverConf := bc.GetServer()
	httpConf := serverConf.GetHttp()

	opts := []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			logging.Server(helper.Logger()),
			tracing.Server(),
			metadata.Server(),
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
