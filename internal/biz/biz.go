// Package biz is the business logic for the Rabbit service.
package biz

import "github.com/google/wire"

var ProviderSetBiz = wire.NewSet(
	NewEmail,
	NewHealth,
	NewEmailConfig,
	NewNamespace,
	NewMessageLog,
)
