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

func NewEmailTemplateRepository(d *data.Data) repository.EmailTemplate {
	return &emailTemplateRepositoryImpl{
		d: d,
	}
}

type emailTemplateRepositoryImpl struct {
	d *data.Data
}

// DeleteEmailTemplate implements repository.EmailTemplate.
func (e *emailTemplateRepositoryImpl) DeleteEmailTemplate(ctx context.Context, uid string) error {
	namespace := middler.GetNamespace(ctx)
	emailTemplate := e.d.BizQuery(namespace).EmailTemplate
	_, err := emailTemplate.WithContext(ctx).Where(emailTemplate.Namespace.Eq(namespace), emailTemplate.UID.Eq(uid)).Delete()
	return err
}

// GetEmailTemplate implements repository.EmailTemplate.
func (e *emailTemplateRepositoryImpl) GetEmailTemplate(ctx context.Context, uid string) (*do.EmailTemplate, error) {
	namespace := middler.GetNamespace(ctx)
	emailTemplate := e.d.BizQuery(namespace).EmailTemplate
	emailTemplateDO, err := emailTemplate.WithContext(ctx).Where(emailTemplate.Namespace.Eq(namespace), emailTemplate.UID.Eq(uid)).First()
	if err != nil {
		return nil, err
	}
	return emailTemplateDO, nil
}

// ListEmailTemplate implements repository.EmailTemplate.
func (e *emailTemplateRepositoryImpl) ListEmailTemplate(ctx context.Context, req *bo.ListEmailTemplateBo) (*bo.PageResponseBo[*do.EmailTemplate], error) {
	namespace := middler.GetNamespace(ctx)
	emailTemplate := e.d.BizQuery(namespace).EmailTemplate
	wrappers := emailTemplate.WithContext(ctx).Where(emailTemplate.Namespace.Eq(namespace))
	if strutil.IsNotEmpty(req.Keyword) {
		wrappers = wrappers.Where(emailTemplate.Name.Like("%" + req.Keyword + "%"))
	}
	if req.Status.Exist() {
		wrappers = wrappers.Where(emailTemplate.Status.Eq(req.Status.GetValue()))
	}
	if pointer.IsNotNil(req.PageRequestBo) {
		total, err := wrappers.Count()
		if err != nil {
			return nil, err
		}
		req.WithTotal(total)
		wrappers = wrappers.Limit(req.Limit()).Offset(req.Offset())
	}
	emailTemplates, err := wrappers.Order(emailTemplate.CreatedAt.Desc()).Find()
	if err != nil {
		return nil, err
	}
	return bo.NewPageResponseBo(req.PageRequestBo, emailTemplates), nil
}

// SaveEmailTemplate implements repository.EmailTemplate.
func (e *emailTemplateRepositoryImpl) SaveEmailTemplate(ctx context.Context, req *do.EmailTemplate) error {
	namespace := middler.GetNamespace(ctx)
	emailTemplate := e.d.BizQuery(namespace).EmailTemplate
	req.WithNamespace(namespace)
	wrappers := emailTemplate.WithContext(ctx)
	if strutil.IsNotEmpty(req.UID) {
		wrappers = wrappers.Where(emailTemplate.UID.Eq(req.UID), emailTemplate.Namespace.Eq(namespace))
	} else {
		req.UID = uuid.New().String()
		req.WithCreator(ctx)
	}
	return wrappers.Save(req)
}

// UpdateEmailTemplateStatus implements repository.EmailTemplate.
func (e *emailTemplateRepositoryImpl) UpdateEmailTemplateStatus(ctx context.Context, uid string, status vobj.GlobalStatus) error {
	namespace := middler.GetNamespace(ctx)
	emailTemplate := e.d.BizQuery(namespace).EmailTemplate
	_, err := emailTemplate.WithContext(ctx).Where(emailTemplate.Namespace.Eq(namespace), emailTemplate.UID.Eq(uid)).Update(emailTemplate.Status, status)
	return err
}
