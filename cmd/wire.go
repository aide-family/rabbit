//go:build wireinject
// +build wireinject

// Package cmd is the root command for the Rabbit service
package cmd

import (
	"github.com/go-kratos/kratos/v2"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	"github.com/aide-family/rabbit/internal/conf"
	"github.com/aide-family/rabbit/internal/server"
	"github.com/aide-family/rabbit/internal/service"
)

func wireApp(bc *conf.Bootstrap, helper *klog.Helper) (*kratos.App, func(), error) {
	panic(wire.Build(
		server.ProviderSetServer,
		service.ProviderSetService,
		newApp,
	))
}
