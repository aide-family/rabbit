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

func NewEmailConfigRepository(d *data.Data) repository.EmailConfig {
	return &emailConfigRepositoryImpl{
		d: d,
	}
}

type emailConfigRepositoryImpl struct {
	d *data.Data
}

// DeleteEmailConfig implements repository.EmailConfig.
func (e *emailConfigRepositoryImpl) DeleteEmailConfig(ctx context.Context, uid string) error {
	namespace := middler.GetNamespace(ctx)
	emailConfig := e.d.BizDB(namespace).EmailConfig
	_, err := emailConfig.WithContext(ctx).Where(emailConfig.Namespace.Eq(namespace), emailConfig.UID.Eq(uid)).Delete()
	return err
}

// GetEmailConfig implements repository.EmailConfig.
func (e *emailConfigRepositoryImpl) GetEmailConfig(ctx context.Context, uid string) (*do.EmailConfig, error) {
	namespace := middler.GetNamespace(ctx)
	emailConfig := e.d.BizDB(namespace).EmailConfig
	emailConfigDO, err := emailConfig.WithContext(ctx).Where(emailConfig.Namespace.Eq(namespace), emailConfig.UID.Eq(uid)).First()
	if err != nil {
		return nil, err
	}
	return emailConfigDO, nil
}

// ListEmailConfig implements repository.EmailConfig.
func (e *emailConfigRepositoryImpl) ListEmailConfig(ctx context.Context, req *bo.ListEmailConfigBo) (*bo.PageResponseBo[*do.EmailConfig], error) {
	namespace := middler.GetNamespace(ctx)
	emailConfig := e.d.BizDB(namespace).EmailConfig
	wrappers := emailConfig.WithContext(ctx).Where(emailConfig.Namespace.Eq(namespace))
	if strutil.IsNotEmpty(req.Keyword) {
		wrappers = wrappers.Where(emailConfig.Name.Like("%" + req.Keyword + "%"))
	}
	if req.Status.Exist() {
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

// SaveEmailConfig implements repository.EmailConfig.
func (e *emailConfigRepositoryImpl) SaveEmailConfig(ctx context.Context, req *do.EmailConfig) error {
	namespace := middler.GetNamespace(ctx)
	emailConfig := e.d.BizDB(namespace).EmailConfig
	req.WithNamespace(namespace)
	wrappers := emailConfig.WithContext(ctx)
	if strutil.IsNotEmpty(req.UID) {
		wrappers = wrappers.Where(emailConfig.UID.Eq(req.UID), emailConfig.Namespace.Eq(namespace))
	} else {
		req.UID = uuid.New().String()
		req.WithCreator(ctx)
	}
	return wrappers.Save(req)
}

// UpdateEmailConfigStatus implements repository.EmailConfig.
func (e *emailConfigRepositoryImpl) UpdateEmailConfigStatus(ctx context.Context, uid string, status vobj.GlobalStatus) error {
	namespace := middler.GetNamespace(ctx)
	emailConfig := e.d.BizDB(namespace).EmailConfig
	_, err := emailConfig.WithContext(ctx).Where(emailConfig.Namespace.Eq(namespace), emailConfig.UID.Eq(uid)).Update(emailConfig.Status, status)
	return err
}
