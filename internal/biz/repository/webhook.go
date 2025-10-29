package repository

import (
	"context"

	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/do"
	"github.com/aide-family/rabbit/internal/biz/vobj"
)

type WebhookConfig interface {
	SaveWebhookConfig(ctx context.Context, req *do.WebhookConfig) error
	UpdateWebhookStatus(ctx context.Context, uid string, status vobj.GlobalStatus) error
	DeleteWebhookConfig(ctx context.Context, uid string) error
	GetWebhookConfig(ctx context.Context, uid string) (*do.WebhookConfig, error)
	ListWebhookConfig(ctx context.Context, req *bo.ListWebhookBo) (*bo.PageResponseBo[*do.WebhookConfig], error)
}

type WebhookTemplate interface {
	SaveWebhookTemplate(ctx context.Context, req *do.WebhookTemplate) error
	UpdateWebhookTemplateStatus(ctx context.Context, uid string, status vobj.GlobalStatus) error
	DeleteWebhookTemplate(ctx context.Context, uid string) error
	GetWebhookTemplate(ctx context.Context, uid string) (*do.WebhookTemplate, error)
	ListWebhookTemplate(ctx context.Context, req *bo.ListWebhookTemplateBo) (*bo.PageResponseBo[*do.WebhookTemplate], error)
}
