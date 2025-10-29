package impl

import (
	"context"

	"github.com/aide-family/magicbox/pointer"
	"github.com/aide-family/magicbox/strutil"

	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/do"
	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/internal/biz/vobj"
	"github.com/aide-family/rabbit/internal/data"
	"github.com/aide-family/rabbit/pkg/middler"
)

func NewWebhookTemplateRepository(d *data.Data) repository.WebhookTemplate {
	return &webhookTemplateRepositoryImpl{
		d: d,
	}
}

type webhookTemplateRepositoryImpl struct {
	d *data.Data
}

// SaveWebhookTemplate implements repository.WebhookTemplate.
func (w *webhookTemplateRepositoryImpl) SaveWebhookTemplate(ctx context.Context, req *do.WebhookTemplate) error {
	namespace := middler.GetNamespace(ctx)
	webhookTemplate := w.d.BizQuery(namespace).WebhookTemplate
	req.WithNamespace(namespace)
	wrappers := webhookTemplate.WithContext(ctx)
	if strutil.IsNotEmpty(req.UID) {
		wrappers = wrappers.Where(webhookTemplate.UID.Eq(req.UID), webhookTemplate.Namespace.Eq(namespace))
		_, err := wrappers.Updates(req)
		return err
	}
	return wrappers.Create(req)
}

// UpdateWebhookTemplateStatus implements repository.WebhookTemplate.
func (w *webhookTemplateRepositoryImpl) UpdateWebhookTemplateStatus(ctx context.Context, uid string, status vobj.GlobalStatus) error {
	namespace := middler.GetNamespace(ctx)
	webhookTemplate := w.d.BizQuery(namespace).WebhookTemplate
	wrappers := webhookTemplate.WithContext(ctx).Where(webhookTemplate.Namespace.Eq(namespace), webhookTemplate.UID.Eq(uid))
	_, err := wrappers.Update(webhookTemplate.Status, status)
	return err
}

// DeleteWebhookTemplate implements repository.WebhookTemplate.
func (w *webhookTemplateRepositoryImpl) DeleteWebhookTemplate(ctx context.Context, uid string) error {
	namespace := middler.GetNamespace(ctx)
	webhookTemplate := w.d.BizQuery(namespace).WebhookTemplate
	wrappers := webhookTemplate.WithContext(ctx).Where(webhookTemplate.Namespace.Eq(namespace), webhookTemplate.UID.Eq(uid))
	_, err := wrappers.Delete()
	return err
}

// GetWebhookTemplate implements repository.WebhookTemplate.
func (w *webhookTemplateRepositoryImpl) GetWebhookTemplate(ctx context.Context, uid string) (*do.WebhookTemplate, error) {
	namespace := middler.GetNamespace(ctx)
	webhookTemplate := w.d.BizQuery(namespace).WebhookTemplate
	wrappers := webhookTemplate.WithContext(ctx).Where(webhookTemplate.Namespace.Eq(namespace), webhookTemplate.UID.Eq(uid))
	return wrappers.First()
}

// ListWebhookTemplate implements repository.WebhookTemplate.
func (w *webhookTemplateRepositoryImpl) ListWebhookTemplate(ctx context.Context, req *bo.ListWebhookTemplateBo) (*bo.PageResponseBo[*do.WebhookTemplate], error) {
	namespace := middler.GetNamespace(ctx)
	webhookTemplate := w.d.BizQuery(namespace).WebhookTemplate
	wrappers := webhookTemplate.WithContext(ctx).Where(webhookTemplate.Namespace.Eq(namespace))
	if strutil.IsNotEmpty(req.Keyword) {
		wrappers = wrappers.Where(webhookTemplate.Name.Like("%" + req.Keyword + "%"))
	}
	if req.Status.Exist() && !req.Status.IsUnknown() {
		wrappers = wrappers.Where(webhookTemplate.Status.Eq(req.Status.GetValue()))
	}
	if req.App.Exist() && !req.App.IsUnknown() {
		wrappers = wrappers.Where(webhookTemplate.App.Eq(req.App.GetValue()))
	}
	if pointer.IsNotNil(req.PageRequestBo) {
		total, err := wrappers.Count()
		if err != nil {
			return nil, err
		}
		req.WithTotal(total)
		wrappers = wrappers.Limit(req.Limit()).Offset(req.Offset())
	}
	webhookTemplates, err := wrappers.Order(webhookTemplate.CreatedAt.Desc()).Find()
	if err != nil {
		return nil, err
	}
	return bo.NewPageResponseBo(req.PageRequestBo, webhookTemplates), nil
}
