// Package data is the data package for the Rabbit service.
package data

import (
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	"github.com/aide-family/rabbit/internal/conf"
)

// ProviderSetData is a set of data providers.
var ProviderSetData = wire.NewSet(New)

type Data struct{}

// New a data and returns.
func New(c *conf.Bootstrap, logger *klog.Helper) (*Data, func(), error) {
	return nil, nil, nil
}
