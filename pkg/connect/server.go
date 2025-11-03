package connect

import "github.com/go-kratos/kratos/v2/registry"

type Registry interface {
	registry.Registrar
	registry.Discovery
}
