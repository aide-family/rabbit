package repository

import (
	"context"

	"github.com/bwmarrin/snowflake"

	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/do"
	"github.com/aide-family/rabbit/internal/biz/vobj"
)

type WebhookConfig interface {
	CreateWebhookConfig(ctx context.Context, req *do.WebhookConfig) error
	UpdateWebhookConfig(ctx context.Context, req *do.WebhookConfig) error
	UpdateWebhookStatus(ctx context.Context, uid snowflake.ID, status vobj.GlobalStatus) error
	DeleteWebhookConfig(ctx context.Context, uid snowflake.ID) error
	GetWebhookConfig(ctx context.Context, uid snowflake.ID) (*do.WebhookConfig, error)
	GetWebhookConfigByName(ctx context.Context, name string) (*do.WebhookConfig, error)
	ListWebhookConfig(ctx context.Context, req *bo.ListWebhookBo) (*bo.PageResponseBo[*do.WebhookConfig], error)
	SelectWebhookConfig(ctx context.Context, req *bo.SelectWebhookBo) (*bo.SelectWebhookResult, error)
}
