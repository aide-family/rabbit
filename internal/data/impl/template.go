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

func NewTemplateRepository(d *data.Data) repository.Template {
	return &templateRepositoryImpl{
		d: d,
	}
}

type templateRepositoryImpl struct {
	d *data.Data
}

// CreateTemplate implements repository.Template.
func (t *templateRepositoryImpl) CreateTemplate(ctx context.Context, req *do.Template) error {
	namespace := middler.GetNamespace(ctx)
	template := t.d.BizQuery(ctx, namespace).Template
	return template.WithContext(ctx).Create(req)
}

func (t *templateRepositoryImpl) UpdateTemplate(ctx context.Context, req *do.Template) error {
	namespace := middler.GetNamespace(ctx)
	template := t.d.BizQuery(ctx, namespace).Template
	wrappers := template.WithContext(ctx).Where(template.Namespace.Eq(namespace), template.UID.Eq(req.UID.Int64()))
	_, err := wrappers.Updates(req)
	return err
}

// UpdateTemplateStatus implements repository.Template.
func (t *templateRepositoryImpl) UpdateTemplateStatus(ctx context.Context, uid snowflake.ID, status vobj.GlobalStatus) error {
	namespace := middler.GetNamespace(ctx)
	template := t.d.BizQuery(ctx, namespace).Template
	wrappers := template.WithContext(ctx).Where(template.Namespace.Eq(namespace), template.UID.Eq(uid.Int64()))
	_, err := wrappers.Update(template.Status, status)
	return err
}

// DeleteTemplate implements repository.Template.
func (t *templateRepositoryImpl) DeleteTemplate(ctx context.Context, uid snowflake.ID) error {
	namespace := middler.GetNamespace(ctx)
	template := t.d.BizQuery(ctx, namespace).Template
	wrappers := template.WithContext(ctx).Where(template.Namespace.Eq(namespace), template.UID.Eq(uid.Int64()))
	_, err := wrappers.Delete()
	return err
}

// GetTemplate implements repository.Template.
func (t *templateRepositoryImpl) GetTemplate(ctx context.Context, uid snowflake.ID) (*do.Template, error) {
	namespace := middler.GetNamespace(ctx)
	template := t.d.BizQuery(ctx, namespace).Template
	wrappers := template.WithContext(ctx).Where(template.Namespace.Eq(namespace), template.UID.Eq(uid.Int64()))
	return wrappers.First()
}

// GetTemplateByName implements repository.Template.
func (t *templateRepositoryImpl) GetTemplateByName(ctx context.Context, name string) (*do.Template, error) {
	namespace := middler.GetNamespace(ctx)
	template := t.d.BizQuery(ctx, namespace).Template
	wrappers := template.WithContext(ctx).Where(template.Namespace.Eq(namespace), template.Name.Eq(name))
	return wrappers.First()
}

// ListTemplate implements repository.Template.
func (t *templateRepositoryImpl) ListTemplate(ctx context.Context, req *bo.ListTemplateBo) (*bo.PageResponseBo[*do.Template], error) {
	namespace := middler.GetNamespace(ctx)
	template := t.d.BizQuery(ctx, namespace).Template
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

// SelectTemplate implements repository.Template.
func (t *templateRepositoryImpl) SelectTemplate(ctx context.Context, req *bo.SelectTemplateBo) (*bo.SelectTemplateResult, error) {
	namespace := middler.GetNamespace(ctx)
	template := t.d.BizQuery(ctx, namespace).Template
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

	// 获取总数
	total, err := wrappers.Count()
	if err != nil {
		return nil, err
	}

	// 游标分页：如果提供了lastUID，则查询UID小于lastUID的记录
	if req.LastUID > 0 {
		wrappers = wrappers.Where(template.UID.Lt(req.LastUID.Int64()))
	}

	// 限制返回数量
	wrappers = wrappers.Limit(int(req.Limit))

	// 按UID倒序排列（snowflake ID按时间生成，与CreatedAt一致）
	templates, err := wrappers.Order(template.UID.Desc()).Find()
	if err != nil {
		return nil, err
	}

	// 获取最后一个UID，用于下次分页
	var lastUID snowflake.ID
	if len(templates) > 0 {
		lastUID = templates[len(templates)-1].UID
	}

	return &bo.SelectTemplateResult{
		Items:   templates,
		Total:   total,
		LastUID: lastUID,
	}, nil
}
