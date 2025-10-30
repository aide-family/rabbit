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

func NewTemplateRepository(d *data.Data) repository.Template {
	return &templateRepositoryImpl{
		d: d,
	}
}

type templateRepositoryImpl struct {
	d *data.Data
}

// SaveTemplate implements repository.Template.
func (t *templateRepositoryImpl) SaveTemplate(ctx context.Context, req *do.Template) error {
	namespace := middler.GetNamespace(ctx)
	template := t.d.BizQuery(namespace).Template
	req.WithNamespace(namespace)
	wrappers := template.WithContext(ctx)
	if strutil.IsNotEmpty(req.UID) {
		wrappers = wrappers.Where(template.UID.Eq(req.UID), template.Namespace.Eq(namespace))
		_, err := wrappers.Updates(req)
		return err
	}
	return wrappers.Create(req)
}

// UpdateTemplateStatus implements repository.Template.
func (t *templateRepositoryImpl) UpdateTemplateStatus(ctx context.Context, uid string, status vobj.GlobalStatus) error {
	namespace := middler.GetNamespace(ctx)
	template := t.d.BizQuery(namespace).Template
	wrappers := template.WithContext(ctx).Where(template.Namespace.Eq(namespace), template.UID.Eq(uid))
	_, err := wrappers.Update(template.Status, status)
	return err
}

// DeleteTemplate implements repository.Template.
func (t *templateRepositoryImpl) DeleteTemplate(ctx context.Context, uid string) error {
	namespace := middler.GetNamespace(ctx)
	template := t.d.BizQuery(namespace).Template
	wrappers := template.WithContext(ctx).Where(template.Namespace.Eq(namespace), template.UID.Eq(uid))
	_, err := wrappers.Delete()
	return err
}

// GetTemplate implements repository.Template.
func (t *templateRepositoryImpl) GetTemplate(ctx context.Context, uid string) (*do.Template, error) {
	namespace := middler.GetNamespace(ctx)
	template := t.d.BizQuery(namespace).Template
	wrappers := template.WithContext(ctx).Where(template.Namespace.Eq(namespace), template.UID.Eq(uid))
	return wrappers.First()
}

// GetTemplateByName implements repository.Template.
func (t *templateRepositoryImpl) GetTemplateByName(ctx context.Context, name string) (*do.Template, error) {
	namespace := middler.GetNamespace(ctx)
	template := t.d.BizQuery(namespace).Template
	wrappers := template.WithContext(ctx).Where(template.Namespace.Eq(namespace), template.Name.Eq(name))
	return wrappers.First()
}

// ListTemplate implements repository.Template.
func (t *templateRepositoryImpl) ListTemplate(ctx context.Context, req *bo.ListTemplateBo) (*bo.PageResponseBo[*do.Template], error) {
	namespace := middler.GetNamespace(ctx)
	template := t.d.BizQuery(namespace).Template
	wrappers := template.WithContext(ctx).Where(template.Namespace.Eq(namespace))
	if strutil.IsNotEmpty(req.Keyword) {
		wrappers = wrappers.Where(template.Name.Like("%" + req.Keyword + "%"))
	}
	if req.Status.Exist() && !req.Status.IsUnknown() {
		wrappers = wrappers.Where(template.Status.Eq(req.Status.GetValue()))
	}
	if req.App.Exist() && !req.App.IsUnknown() {
		wrappers = wrappers.Where(template.App.Eq(req.App.GetValue()))
	}
	if pointer.IsNotNil(req.PageRequestBo) {
		total, err := wrappers.Count()
		if err != nil {
			return nil, err
		}
		req.WithTotal(total)
		wrappers = wrappers.Limit(req.Limit()).Offset(req.Offset())
	}
	templates, err := wrappers.Order(template.CreatedAt.Desc()).Find()
	if err != nil {
		return nil, err
	}
	return bo.NewPageResponseBo(req.PageRequestBo, templates), nil
}
