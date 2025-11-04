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

func NewEmailConfigRepository(d *data.Data) repository.EmailConfig {
	return &emailConfigRepositoryImpl{
		d: d,
	}
}

type emailConfigRepositoryImpl struct {
	d *data.Data
}

// DeleteEmailConfig implements repository.EmailConfig.
func (e *emailConfigRepositoryImpl) DeleteEmailConfig(ctx context.Context, uid snowflake.ID) error {
	namespace := middler.GetNamespace(ctx)
	emailConfig := e.d.BizQuery(ctx, namespace).EmailConfig
	wrappers := emailConfig.WithContext(ctx).Where(emailConfig.Namespace.Eq(namespace), emailConfig.UID.Eq(uid.Int64()))
	_, err := wrappers.Delete()
	return err
}

// GetEmailConfig implements repository.EmailConfig.
func (e *emailConfigRepositoryImpl) GetEmailConfig(ctx context.Context, uid snowflake.ID) (*do.EmailConfig, error) {
	namespace := middler.GetNamespace(ctx)
	emailConfig := e.d.BizQuery(ctx, namespace).EmailConfig
	wrappers := emailConfig.WithContext(ctx).Where(emailConfig.Namespace.Eq(namespace), emailConfig.UID.Eq(uid.Int64()))
	return wrappers.First()
}

// GetEmailConfigByName implements repository.EmailConfig.
func (e *emailConfigRepositoryImpl) GetEmailConfigByName(ctx context.Context, name string) (*do.EmailConfig, error) {
	namespace := middler.GetNamespace(ctx)
	emailConfig := e.d.BizQuery(ctx, namespace).EmailConfig
	wrappers := emailConfig.WithContext(ctx).Where(emailConfig.Namespace.Eq(namespace), emailConfig.Name.Eq(name))
	return wrappers.First()
}

// ListEmailConfig implements repository.EmailConfig.
func (e *emailConfigRepositoryImpl) ListEmailConfig(ctx context.Context, req *bo.ListEmailConfigBo) (*bo.PageResponseBo[*do.EmailConfig], error) {
	namespace := middler.GetNamespace(ctx)
	emailConfig := e.d.BizQuery(ctx, namespace).EmailConfig
	wrappers := emailConfig.WithContext(ctx).Where(emailConfig.Namespace.Eq(namespace))
	if strutil.IsNotEmpty(req.Keyword) {
		wrappers = wrappers.Where(emailConfig.Name.Like("%" + req.Keyword + "%"))
	}
	if req.Status.Exist() && !req.Status.IsUnknown() {
		wrappers = wrappers.Where(emailConfig.Status.Eq(req.Status.GetValue()))
	}
	if pointer.IsNotNil(req.PageRequestBo) {
		total, err := wrappers.Count()
		if err != nil {
			return nil, err
		}
		req.WithTotal(total)
		wrappers = wrappers.Limit(req.Limit()).Offset(req.Offset())
	}
	emailConfigs, err := wrappers.Order(emailConfig.CreatedAt.Desc()).Find()
	if err != nil {
		return nil, err
	}
	return bo.NewPageResponseBo(req.PageRequestBo, emailConfigs), nil
}

// CreateEmailConfig implements repository.EmailConfig.
func (e *emailConfigRepositoryImpl) CreateEmailConfig(ctx context.Context, req *do.EmailConfig) error {
	namespace := middler.GetNamespace(ctx)
	emailConfig := e.d.BizQuery(ctx, namespace).EmailConfig
	wrappers := emailConfig.WithContext(ctx)
	return wrappers.Create(req)
}

// UpdateEmailConfig implements repository.EmailConfig.
func (e *emailConfigRepositoryImpl) UpdateEmailConfig(ctx context.Context, req *do.EmailConfig) error {
	namespace := middler.GetNamespace(ctx)
	emailConfig := e.d.BizQuery(ctx, namespace).EmailConfig
	wrappers := emailConfig.WithContext(ctx).Where(emailConfig.UID.Eq(req.UID.Int64()), emailConfig.Namespace.Eq(namespace))
	_, err := wrappers.Updates(req)
	return err
}

// UpdateEmailConfigStatus implements repository.EmailConfig.
func (e *emailConfigRepositoryImpl) UpdateEmailConfigStatus(ctx context.Context, uid snowflake.ID, status vobj.GlobalStatus) error {
	namespace := middler.GetNamespace(ctx)
	emailConfig := e.d.BizQuery(ctx, namespace).EmailConfig
	wrappers := emailConfig.WithContext(ctx).Where(emailConfig.Namespace.Eq(namespace), emailConfig.UID.Eq(uid.Int64()))
	_, err := wrappers.Update(emailConfig.Status, status)
	return err
}
