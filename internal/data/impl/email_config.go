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

// SelectEmailConfig implements repository.EmailConfig.
func (e *emailConfigRepositoryImpl) SelectEmailConfig(ctx context.Context, req *bo.SelectEmailConfigBo) (*bo.SelectEmailConfigResult, error) {
	namespace := middler.GetNamespace(ctx)
	emailConfig := e.d.BizQuery(ctx, namespace).EmailConfig
	wrappers := emailConfig.WithContext(ctx).Where(emailConfig.Namespace.Eq(namespace))

	if strutil.IsNotEmpty(req.Keyword) {
		wrappers = wrappers.Where(emailConfig.Name.Like("%" + req.Keyword + "%"))
	}
	if req.Status.Exist() && !req.Status.IsUnknown() {
		wrappers = wrappers.Where(emailConfig.Status.Eq(req.Status.GetValue()))
	}

	// 获取总数
	total, err := wrappers.Count()
	if err != nil {
		return nil, err
	}

	// 游标分页：如果提供了lastUID，则查询UID小于lastUID的记录
	if req.LastUID > 0 {
		wrappers = wrappers.Where(emailConfig.UID.Lt(req.LastUID.Int64()))
	}

	// 限制返回数量
	wrappers = wrappers.Limit(int(req.Limit))

	// 按UID倒序排列（snowflake ID按时间生成，与CreatedAt一致）
	emailConfigs, err := wrappers.Order(emailConfig.UID.Desc()).Find()
	if err != nil {
		return nil, err
	}

	// 获取最后一个UID，用于下次分页
	var lastUID snowflake.ID
	if len(emailConfigs) > 0 {
		lastUID = emailConfigs[len(emailConfigs)-1].UID
	}

	return &bo.SelectEmailConfigResult{
		Items:   emailConfigs,
		Total:   total,
		LastUID: lastUID,
	}, nil
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
