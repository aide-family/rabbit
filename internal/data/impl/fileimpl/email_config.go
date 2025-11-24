// Package fileimpl is the implementation of the email config repository for file config
package fileimpl

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"

	"github.com/aide-family/magicbox/safety"
	"github.com/aide-family/magicbox/strutil"
	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/do"
	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/internal/biz/vobj"
	"github.com/aide-family/rabbit/internal/conf"
	"github.com/aide-family/rabbit/internal/data"
	"github.com/aide-family/rabbit/pkg/merr"
	"github.com/aide-family/rabbit/pkg/middler"
)

func NewEmailConfigRepository(d *data.Data) repository.EmailConfig {
	e := &emailConfigRepositoryImpl{
		d:            d,
		emailConfigs: d.GetFileConfig().GetEmails(),
	}
	e.initEmailConfigs()
	d.RegisterReloadFunc(data.KeyEmails, func() {
		e.initEmailConfigs()
	})
	return e
}

type emailConfigRepositoryImpl struct {
	d                    *data.Data
	emailConfigs         []*conf.Config_Email
	emailConfigsWithUID  *safety.SyncMap[string, *safety.SyncMap[snowflake.ID, *do.EmailConfig]]
	emailConfigsWithName *safety.SyncMap[string, *safety.SyncMap[string, *do.EmailConfig]]
}

func (e *emailConfigRepositoryImpl) initEmailConfigs() {
	e.emailConfigsWithUID = safety.NewSyncMap(make(map[string]*safety.SyncMap[snowflake.ID, *do.EmailConfig]))
	e.emailConfigsWithName = safety.NewSyncMap(make(map[string]*safety.SyncMap[string, *do.EmailConfig]))
	for _, emailConfig := range e.emailConfigs {
		namespace := emailConfig.GetNamespace()
		uid := snowflake.ParseInt64(emailConfig.GetUid())
		name := emailConfig.GetName()
		if _, ok := e.emailConfigsWithUID.Get(namespace); !ok {
			e.emailConfigsWithUID.Set(namespace, safety.NewSyncMap(map[snowflake.ID]*do.EmailConfig{}))
			e.emailConfigsWithName.Set(namespace, safety.NewSyncMap(map[string]*do.EmailConfig{}))
		}
		item := e.toDoEmailConfig(emailConfig)
		if namespaceEmailConfigsByName, ok := e.emailConfigsWithName.Get(namespace); ok {
			namespaceEmailConfigsByName.Set(name, item)
		}
		if namespaceEmailConfigsByUID, ok := e.emailConfigsWithUID.Get(namespace); ok {
			namespaceEmailConfigsByUID.Set(uid, item)
		}
	}
}

func (e *emailConfigRepositoryImpl) toDoEmailConfig(emailConfig *conf.Config_Email) *do.EmailConfig {
	createdAt, _ := time.Parse(time.DateTime, emailConfig.GetCreatedAt())
	updatedAt, _ := time.Parse(time.DateTime, emailConfig.GetUpdatedAt())
	return &do.EmailConfig{
		NamespaceModel: do.NamespaceModel{
			Namespace: emailConfig.GetNamespace(),
			BaseModel: do.BaseModel{
				ID:        emailConfig.GetId(),
				UID:       snowflake.ParseInt64(emailConfig.GetUid()),
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
		},
		Name:     emailConfig.GetName(),
		Host:     emailConfig.GetHost(),
		Port:     emailConfig.GetPort(),
		Username: emailConfig.GetUsername(),
		Password: strutil.EncryptString(emailConfig.GetPassword()),
		Status:   vobj.GlobalStatus(emailConfig.GetStatus()),
	}
}

// CreateEmailConfig implements repository.EmailConfig.
func (e *emailConfigRepositoryImpl) CreateEmailConfig(ctx context.Context, req *do.EmailConfig) error {
	return merr.ErrorParamsNotSupportFileConfig()
}

// DeleteEmailConfig implements repository.EmailConfig.
func (e *emailConfigRepositoryImpl) DeleteEmailConfig(ctx context.Context, uid snowflake.ID) error {
	return merr.ErrorParamsNotSupportFileConfig()
}

// GetEmailConfig implements repository.EmailConfig.
func (e *emailConfigRepositoryImpl) GetEmailConfig(ctx context.Context, uid snowflake.ID) (*do.EmailConfig, error) {
	namespace := middler.GetNamespace(ctx)
	emailConfigWithUID, ok := e.emailConfigsWithUID.Get(namespace)
	if !ok {
		return nil, merr.ErrorNotFound("email config not found")
	}
	emailConfig, ok := emailConfigWithUID.Get(uid)
	if !ok {
		return nil, merr.ErrorNotFound("email config not found")
	}
	return emailConfig, nil
}

// GetEmailConfigByName implements repository.EmailConfig.
func (e *emailConfigRepositoryImpl) GetEmailConfigByName(ctx context.Context, name string) (*do.EmailConfig, error) {
	namespace := middler.GetNamespace(ctx)
	emailConfigWithName, ok := e.emailConfigsWithName.Get(namespace)
	if !ok {
		return nil, merr.ErrorNotFound("email config not found")
	}
	emailConfig, ok := emailConfigWithName.Get(name)
	if !ok {
		return nil, merr.ErrorNotFound("email config not found")
	}
	return emailConfig, nil
}

// ListEmailConfig implements repository.EmailConfig.
func (e *emailConfigRepositoryImpl) ListEmailConfig(ctx context.Context, req *bo.ListEmailConfigBo) (*bo.PageResponseBo[*do.EmailConfig], error) {
	namespace := middler.GetNamespace(ctx)
	emailConfigWithUID, ok := e.emailConfigsWithUID.Get(namespace)
	if !ok {
		pageRequestBo := bo.NewPageRequestBo(req.Page, req.PageSize)
		pageRequestBo.WithTotal(0)
		req.PageRequestBo = pageRequestBo
		return bo.NewPageResponseBo(req.PageRequestBo, []*do.EmailConfig{}), nil
	}
	emailConfigs := make([]*do.EmailConfig, 0, emailConfigWithUID.Len())
	for _, emailConfig := range emailConfigWithUID.Values() {
		if strutil.IsNotEmpty(req.Keyword) && !strings.Contains(emailConfig.Name, req.Keyword) {
			continue
		}
		if req.Status.Exist() && !req.Status.IsUnknown() && emailConfig.Status != req.Status {
			continue
		}
		emailConfigs = append(emailConfigs, emailConfig)
	}
	total := int64(len(emailConfigs))
	pageRequestBo := bo.NewPageRequestBo(req.Page, req.PageSize)
	pageRequestBo.WithTotal(total)
	req.PageRequestBo = pageRequestBo
	sort.Slice(emailConfigs, func(i, j int) bool {
		return emailConfigs[i].UID < emailConfigs[j].UID
	})
	return bo.NewPageResponseBo(req.PageRequestBo, emailConfigs), nil
}

// SelectEmailConfig implements repository.EmailConfig.
func (e *emailConfigRepositoryImpl) SelectEmailConfig(ctx context.Context, req *bo.SelectEmailConfigBo) (*bo.SelectEmailConfigResult, error) {
	namespace := middler.GetNamespace(ctx)
	emailConfigWithUID, ok := e.emailConfigsWithUID.Get(namespace)
	if !ok {
		return &bo.SelectEmailConfigResult{
			Items:   []*do.EmailConfig{},
			Total:   0,
			LastUID: 0,
		}, nil
	}
	emailConfigs := make([]*do.EmailConfig, 0, emailConfigWithUID.Len())
	for _, emailConfig := range emailConfigWithUID.Values() {
		if strutil.IsNotEmpty(req.Keyword) && !strings.Contains(emailConfig.Name, req.Keyword) {
			continue
		}
		if req.Status.Exist() && !req.Status.IsUnknown() && emailConfig.Status != req.Status {
			continue
		}
		if req.LastUID > 0 && emailConfig.UID >= req.LastUID {
			continue
		}
		emailConfigs = append(emailConfigs, emailConfig)
	}
	total := int64(len(emailConfigs))
	sort.Slice(emailConfigs, func(i, j int) bool {
		return emailConfigs[i].UID > emailConfigs[j].UID
	})
	if int32(len(emailConfigs)) > req.Limit {
		emailConfigs = emailConfigs[:req.Limit]
	}
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

// UpdateEmailConfig implements repository.EmailConfig.
func (e *emailConfigRepositoryImpl) UpdateEmailConfig(ctx context.Context, req *do.EmailConfig) error {
	return merr.ErrorParamsNotSupportFileConfig()
}

// UpdateEmailConfigStatus implements repository.EmailConfig.
func (e *emailConfigRepositoryImpl) UpdateEmailConfigStatus(ctx context.Context, uid snowflake.ID, status vobj.GlobalStatus) error {
	return merr.ErrorParamsNotSupportFileConfig()
}
