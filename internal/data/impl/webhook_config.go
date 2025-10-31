package impl

import (
	"context"

	"github.com/aide-family/magicbox/pointer"
	"github.com/aide-family/magicbox/strutil"
	"github.com/bwmarrin/snowflake"

	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/do"
	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/internal/biz/vobj"
	"github.com/aide-family/rabbit/internal/data"
	"github.com/aide-family/rabbit/pkg/middler"
)

func NewWebhookConfigRepository(d *data.Data) repository.WebhookConfig {
	return &webhookConfigRepositoryImpl{
		d: d,
	}
}

type webhookConfigRepositoryImpl struct {
	d *data.Data
}

// SaveWebhookConfig implements repository.Webhook.
func (w *webhookConfigRepositoryImpl) SaveWebhookConfig(ctx context.Context, req *do.WebhookConfig) error {
	namespace := middler.GetNamespace(ctx)
	webhookConfig := w.d.BizQuery(namespace).WebhookConfig
	wrappers := webhookConfig.WithContext(ctx)
	if strutil.IsNotEmpty(req.UID.String()) {
		wrappers = wrappers.Where(webhookConfig.UID.Eq(req.UID.Int64()), webhookConfig.Namespace.Eq(namespace))
		_, err := wrappers.Updates(req)
		return err
	}
	return wrappers.Create(req)
}

// UpdateWebhookStatus implements repository.Webhook.
func (w *webhookConfigRepositoryImpl) UpdateWebhookStatus(ctx context.Context, uid snowflake.ID, status vobj.GlobalStatus) error {
	namespace := middler.GetNamespace(ctx)
	webhookConfig := w.d.BizQuery(namespace).WebhookConfig
	wrappers := webhookConfig.WithContext(ctx).Where(webhookConfig.Namespace.Eq(namespace), webhookConfig.UID.Eq(uid.Int64()))
	_, err := wrappers.Update(webhookConfig.Status, status)
	return err
}

// DeleteWebhookConfig implements repository.Webhook.
func (w *webhookConfigRepositoryImpl) DeleteWebhookConfig(ctx context.Context, uid snowflake.ID) error {
	namespace := middler.GetNamespace(ctx)
	webhookConfig := w.d.BizQuery(namespace).WebhookConfig
	wrappers := webhookConfig.WithContext(ctx).Where(webhookConfig.Namespace.Eq(namespace), webhookConfig.UID.Eq(uid.Int64()))
	_, err := wrappers.Delete()
	return err
}

// GetWebhookConfig implements repository.Webhook.
func (w *webhookConfigRepositoryImpl) GetWebhookConfig(ctx context.Context, uid snowflake.ID) (*do.WebhookConfig, error) {
	namespace := middler.GetNamespace(ctx)
	webhookConfig := w.d.BizQuery(namespace).WebhookConfig
	wrappers := webhookConfig.WithContext(ctx).Where(webhookConfig.Namespace.Eq(namespace), webhookConfig.UID.Eq(uid.Int64()))
	return wrappers.First()
}

// GetWebhookConfigByName implements repository.WebhookConfig.
func (w *webhookConfigRepositoryImpl) GetWebhookConfigByName(ctx context.Context, name string) (*do.WebhookConfig, error) {
	namespace := middler.GetNamespace(ctx)
	webhookConfig := w.d.BizQuery(namespace).WebhookConfig
	wrappers := webhookConfig.WithContext(ctx).Where(webhookConfig.Namespace.Eq(namespace), webhookConfig.Name.Eq(name))
	return wrappers.First()
}

// ListWebhookConfig implements repository.Webhook.
func (w *webhookConfigRepositoryImpl) ListWebhookConfig(ctx context.Context, req *bo.ListWebhookBo) (*bo.PageResponseBo[*do.WebhookConfig], error) {
	namespace := middler.GetNamespace(ctx)
	webhookConfig := w.d.BizQuery(namespace).WebhookConfig
	wrappers := webhookConfig.WithContext(ctx).Where(webhookConfig.Namespace.Eq(namespace))
	if req.App.Exist() && !req.App.IsUnknown() {
		wrappers = wrappers.Where(webhookConfig.App.Eq(req.App.GetValue()))
	}
	if strutil.IsNotEmpty(req.Keyword) {
		wrappers = wrappers.Where(webhookConfig.Name.Like("%" + req.Keyword + "%"))
	}
	if pointer.IsNotNil(req.PageRequestBo) {
		total, err := wrappers.Count()
		if err != nil {
			return nil, err
		}
		req.WithTotal(total)
		wrappers = wrappers.Limit(req.Limit()).Offset(req.Offset())
	}
	webhookConfigs, err := wrappers.Order(webhookConfig.CreatedAt.Desc()).Find()
	if err != nil {
		return nil, err
	}
	return bo.NewPageResponseBo(req.PageRequestBo, webhookConfigs), nil
}
