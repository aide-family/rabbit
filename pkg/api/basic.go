// Package api is the api package for the rabbit service.
package api

import (
	nethttp "net/http"
	"strings"

	"github.com/aide-family/magicbox/pointer"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"

	"github.com/aide-family/rabbit/pkg/middler"
)

type BasicAuthConfig interface {
	GetEnabled() string
	GetUsername() string
	GetPassword() string
}

type HandlerBinding struct {
	Name      string
	Enabled   bool
	BasicAuth BasicAuthConfig
	Handler   nethttp.Handler
	Path      string
}

func BindHandlerWithAuth(httpSrv *http.Server, binding HandlerBinding) {
	if !binding.Enabled {
		klog.Debugf("%s is not enabled", binding.Name)
		return
	}

	handler := binding.Handler
	basicAuth := binding.BasicAuth
	if pointer.IsNotNil(basicAuth) && strings.EqualFold(basicAuth.GetEnabled(), "true") {
		handler = middler.BasicAuthMiddleware(basicAuth.GetUsername(), basicAuth.GetPassword())(handler)
		klog.Debugf("%s route: %s (Basic Auth: %s:%s)", binding.Name, binding.Path, basicAuth.GetUsername(), basicAuth.GetPassword())
	} else {
		klog.Debugf("%s route: %s (No Basic Auth)", binding.Name, binding.Path)
	}

	httpSrv.HandlePrefix(binding.Path, handler)
}
