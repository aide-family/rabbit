package convert

import (
	"context"

	"github.com/aide-family/magicbox/safety"
	"github.com/aide-family/magicbox/strutil"

	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/data/impl/do"
	"github.com/aide-family/rabbit/pkg/contextx"
	"github.com/aide-family/rabbit/pkg/enum"
)

func ToWebhookConfigDO(ctx context.Context, req *bo.CreateWebhookBo) *do.WebhookConfig {
	model := &do.WebhookConfig{
		App:          req.App,
		NamespaceUID: contextx.GetNamespaceUID(ctx),
		Name:         req.Name,
		URL:          req.URL,
		Method:       req.Method,
		Headers:      safety.NewMap(req.Headers),
		Secret:       strutil.EncryptString(req.Secret),
		Status:       enum.GlobalStatus_ENABLED,
	}
	model.WithCreator(contextx.GetUserUID(ctx))
	return model
}

func ToWebhookConfigItemBo(webhookConfigDO *do.WebhookConfig) *bo.WebhookItemBo {
	return &bo.WebhookItemBo{
		UID:       webhookConfigDO.UID,
		App:       webhookConfigDO.App,
		Name:      webhookConfigDO.Name,
		URL:       webhookConfigDO.URL,
		Method:    webhookConfigDO.Method,
		Headers:   webhookConfigDO.Headers.Map(),
		Secret:    string(webhookConfigDO.Secret),
		Status:    webhookConfigDO.Status,
		CreatedAt: webhookConfigDO.CreatedAt,
		UpdatedAt: webhookConfigDO.UpdatedAt,
	}
}

func ToWebhookConfigItemSelectBo(webhookConfigDO *do.WebhookConfig) *bo.WebhookItemSelectBo {
	return &bo.WebhookItemSelectBo{
		UID:      webhookConfigDO.UID,
		Name:     webhookConfigDO.Name,
		Status:   webhookConfigDO.Status,
		Disabled: webhookConfigDO.Status == enum.GlobalStatus_DISABLED || webhookConfigDO.DeletedAt.Valid,
		Tooltip:  webhookConfigDO.Name,
	}
}
