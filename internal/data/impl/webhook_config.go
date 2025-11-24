package impl

import (
	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/internal/data"
	"github.com/aide-family/rabbit/internal/data/impl/dbimpl"
	"github.com/aide-family/rabbit/internal/data/impl/fileimpl"
)

func NewWebhookConfigRepository(d *data.Data) repository.WebhookConfig {
	newRepo := fileimpl.NewWebhookConfigRepository
	if d.UseDatabase() {
		newRepo = dbimpl.NewWebhookConfigRepository
	}
	return &webhookConfigRepositoryImpl{WebhookConfig: newRepo(d)}
}

type webhookConfigRepositoryImpl struct {
	repository.WebhookConfig
}
