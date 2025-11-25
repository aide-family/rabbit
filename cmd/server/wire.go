//go:build wireinject
// +build wireinject

// Package server is the server command for the Rabbit service
package server

import (
	"github.com/go-kratos/kratos/v2"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	"github.com/aide-family/rabbit/internal/biz"
	"github.com/aide-family/rabbit/internal/conf"
	"github.com/aide-family/rabbit/internal/data"
	"github.com/aide-family/rabbit/internal/data/impl"
	"github.com/aide-family/rabbit/internal/server"
	"github.com/aide-family/rabbit/internal/service"
)

func wireApp(bc *conf.Bootstrap, helper *klog.Helper) (*kratos.App, func(), error) {
	panic(wire.Build(
		server.ProviderSetServer,
		service.ProviderSetService,
		biz.ProviderSetBiz,
		impl.ProviderSetImpl,
		data.ProviderSetData,
		newApp,
	))
}
