package middler

import (
	nethttp "net/http"

	"github.com/aide-family/magicbox/server/middler"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func Cors() http.ServerOption {
	return http.Filter(middler.Cors(&middler.CorsConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{nethttp.MethodGet, nethttp.MethodPost, nethttp.MethodPut, nethttp.MethodDelete, nethttp.MethodOptions},
		MaxAge:       600,
	}))
}
