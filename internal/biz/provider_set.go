// Package biz is the business logic for the rabbit service.
package biz

import "github.com/google/wire"

var ProviderSetBiz = wire.NewSet(
	NewHealth,
	NewNamespace,
	NewLoginBiz,
	NewJob,
	NewEmailConfig,
	NewEmail,
	NewWebhookConfig,
	NewWebhook,
	NewMessageLog,
	NewMessage,
	NewTemplate,
)
