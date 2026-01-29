// Package domain is the domain service implementation.
package domain

import (
	"github.com/aide-family/magicbox/safety"

	"github.com/aide-family/rabbit/pkg/config"
	authv1 "github.com/aide-family/rabbit/pkg/domain/auth/v1"
	namespacev1 "github.com/aide-family/rabbit/pkg/domain/namespace/v1"
)

var globalRegistry = NewRegistry()

func NewRegistry() Registry {
	return &registry{
		namespaceV1: safety.NewSyncMap(make(map[config.DomainConfig_Driver]NamespaceFactoryV1)),
		authV1:      safety.NewSyncMap(make(map[config.DomainConfig_Driver]AuthFactoryV1)),
	}
}

type NamespaceFactoryV1 func(c *config.DomainConfig) (namespacev1.Repository, func() error, error)

type AuthFactoryV1 func(c *config.DomainConfig, jwtConfig *config.JWT) (authv1.Repository, func() error, error)

type Registry interface {
	RegisterNamespaceV1Factory(name config.DomainConfig_Driver, factory NamespaceFactoryV1)
	GetNamespaceV1Factory(name config.DomainConfig_Driver) (NamespaceFactoryV1, bool)
	RegisterAuthV1Factory(name config.DomainConfig_Driver, factory AuthFactoryV1)
	GetAuthV1Factory(name config.DomainConfig_Driver) (AuthFactoryV1, bool)
}

type registry struct {
	namespaceV1 *safety.SyncMap[config.DomainConfig_Driver, NamespaceFactoryV1]
	authV1      *safety.SyncMap[config.DomainConfig_Driver, AuthFactoryV1]
}

func (r *registry) RegisterNamespaceV1Factory(name config.DomainConfig_Driver, factory NamespaceFactoryV1) {
	r.namespaceV1.Set(name, factory)
}

func (r *registry) GetNamespaceV1Factory(name config.DomainConfig_Driver) (NamespaceFactoryV1, bool) {
	return r.namespaceV1.Get(name)
}

func (r *registry) RegisterAuthV1Factory(name config.DomainConfig_Driver, factory AuthFactoryV1) {
	r.authV1.Set(name, factory)
}

func (r *registry) GetAuthV1Factory(name config.DomainConfig_Driver) (AuthFactoryV1, bool) {
	return r.authV1.Get(name)
}

func RegisterNamespaceV1Factory(name config.DomainConfig_Driver, factory NamespaceFactoryV1) {
	globalRegistry.RegisterNamespaceV1Factory(name, factory)
}

func GetNamespaceV1Factory(name config.DomainConfig_Driver) (NamespaceFactoryV1, bool) {
	return globalRegistry.GetNamespaceV1Factory(name)
}

func RegisterAuthV1Factory(name config.DomainConfig_Driver, factory AuthFactoryV1) {
	globalRegistry.RegisterAuthV1Factory(name, factory)
}

func GetAuthV1Factory(name config.DomainConfig_Driver) (AuthFactoryV1, bool) {
	return globalRegistry.GetAuthV1Factory(name)
}
