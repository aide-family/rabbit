//go:build wireinject
// +build wireinject

// Package job is the job command for the Rabbit service
package job

import (
	"github.com/go-kratos/kratos/v2"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	"github.com/aide-family/rabbit/internal/conf"
	"github.com/aide-family/rabbit/internal/data"
	"github.com/aide-family/rabbit/internal/data/impl"
	"github.com/aide-family/rabbit/internal/server"
)

func wireApp(bc *conf.Bootstrap, helper *klog.Helper) (*kratos.App, func(), error) {
	panic(wire.Build(
		server.ProviderSetJob,
		impl.ProviderSetImpl,
		data.ProviderSetData,
		newApp,
	))
}
