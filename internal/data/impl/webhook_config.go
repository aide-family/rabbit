package impl

import (
	"context"

	"github.com/aide-family/magicbox/pointer"
	"github.com/aide-family/magicbox/strutil"
	"github.com/google/uuid"

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
	req.WithNamespace(namespace)
	wrappers := webhookConfig.WithContext(ctx)
	if strutil.IsNotEmpty(req.UID) {
		wrappers = wrappers.Where(webhookConfig.UID.Eq(req.UID), webhookConfig.Namespace.Eq(namespace))
	} else {
		req.UID = uuid.New().String()
		req.WithCreator(ctx)
	}
	return wrappers.Save(req)
}

// UpdateWebhookStatus implements repository.Webhook.
func (w *webhookConfigRepositoryImpl) UpdateWebhookStatus(ctx context.Context, uid string, status vobj.GlobalStatus) error {
	namespace := middler.GetNamespace(ctx)
	webhookConfig := w.d.BizQuery(namespace).WebhookConfig
	_, err := webhookConfig.WithContext(ctx).Where(webhookConfig.Namespace.Eq(namespace), webhookConfig.UID.Eq(uid)).Update(webhookConfig.Status, status)
	return err
}

// DeleteWebhookConfig implements repository.Webhook.
func (w *webhookConfigRepositoryImpl) DeleteWebhookConfig(ctx context.Context, uid string) error {
	namespace := middler.GetNamespace(ctx)
	webhookConfig := w.d.BizQuery(namespace).WebhookConfig
	_, err := webhookConfig.WithContext(ctx).Where(webhookConfig.Namespace.Eq(namespace), webhookConfig.UID.Eq(uid)).Delete()
	return err
}

// GetWebhookConfig implements repository.Webhook.
func (w *webhookConfigRepositoryImpl) GetWebhookConfig(ctx context.Context, uid string) (*do.WebhookConfig, error) {
	namespace := middler.GetNamespace(ctx)
	webhookConfig := w.d.BizQuery(namespace).WebhookConfig
	webhookConfigDO, err := webhookConfig.WithContext(ctx).Where(webhookConfig.Namespace.Eq(namespace), webhookConfig.UID.Eq(uid)).First()
	if err != nil {
		return nil, err
	}
	return webhookConfigDO, nil
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
