package biz

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
	klog "github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"

	"github.com/aide-family/magicbox/safety"
	"github.com/aide-family/magicbox/strutil"
	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/internal/biz/vobj"
	"github.com/aide-family/rabbit/internal/conf"
	"github.com/aide-family/rabbit/pkg/merr"
	"github.com/aide-family/rabbit/pkg/middler"
)

func NewEmailConfig(
	bc *conf.Bootstrap,
	emailConfigRepo repository.EmailConfig,
	helper *klog.Helper,
) *EmailConfig {
	emailConfigs := safety.NewSyncMap(make(map[string]*safety.SyncMap[snowflake.ID, *bo.EmailConfigItemBo]))
	for _, emailConfig := range conf.GetFileConfig().GetEmails() {
		namespace := emailConfig.GetNamespace()
		uid := snowflake.ParseInt64(emailConfig.GetUid())
		if _, ok := emailConfigs.Get(namespace); !ok {
			emailConfigs.Set(namespace, safety.NewSyncMap(map[snowflake.ID]*bo.EmailConfigItemBo{}))
		}
		namespaceEmailConfigs, ok := emailConfigs.Get(namespace)
		if !ok {
			continue
		}
		createdAt, _ := time.Parse(time.DateTime, emailConfig.GetCreatedAt())
		updatedAt, _ := time.Parse(time.DateTime, emailConfig.GetUpdatedAt())
		namespaceEmailConfigs.Set(uid, &bo.EmailConfigItemBo{
			UID:       uid,
			Name:      emailConfig.GetName(),
			Host:      emailConfig.GetHost(),
			Port:      emailConfig.GetPort(),
			Username:  emailConfig.GetUsername(),
			Password:  emailConfig.GetPassword(),
			Status:    vobj.GlobalStatus(emailConfig.GetStatus()),
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
	}
	return &EmailConfig{
		useDatabase:     bc.GetUseDatabase() == "true",
		emailConfigRepo: emailConfigRepo,
		emailConfigs:    emailConfigs,
		helper:          klog.NewHelper(klog.With(helper.Logger(), "biz", "email_config")),
	}
}

type EmailConfig struct {
	helper          *klog.Helper
	useDatabase     bool
	emailConfigRepo repository.EmailConfig
	emailConfigs    *safety.SyncMap[string, *safety.SyncMap[snowflake.ID, *bo.EmailConfigItemBo]]
}

func (c *EmailConfig) CreateEmailConfig(ctx context.Context, req *bo.CreateEmailConfigBo) error {
	if !c.useDatabase {
		return merr.ErrorParamsNotSupportFileConfig()
	}
	doEmailConfig := req.ToDoEmailConfig()
	if _, err := c.emailConfigRepo.GetEmailConfigByName(ctx, doEmailConfig.Name); err == nil {
		return merr.ErrorParams("email config %s already exists", doEmailConfig.Name)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.helper.Errorw("msg", "check email config exists failed", "error", err, "name", doEmailConfig.Name)
		return merr.ErrorInternal("create email config %s failed", doEmailConfig.Name)
	}
	if err := c.emailConfigRepo.CreateEmailConfig(ctx, doEmailConfig); err != nil {
		c.helper.Errorw("msg", "create email config failed", "error", err, "name", doEmailConfig.Name)
		return merr.ErrorInternal("create email config %s failed", doEmailConfig.Name)
	}
	return nil
}

func (c *EmailConfig) UpdateEmailConfig(ctx context.Context, req *bo.UpdateEmailConfigBo) error {
	if !c.useDatabase {
		return merr.ErrorParamsNotSupportFileConfig()
	}
	doEmailConfig := req.ToDoEmailConfig()
	if err := c.emailConfigRepo.UpdateEmailConfig(ctx, doEmailConfig); err != nil {
		c.helper.Errorw("msg", "update email config failed", "error", err, "name", doEmailConfig.Name)
		return merr.ErrorInternal("update email config %s failed", doEmailConfig.Name)
	}
	return nil
}

func (c *EmailConfig) UpdateEmailConfigStatus(ctx context.Context, req *bo.UpdateEmailConfigStatusBo) error {
	if !c.useDatabase {
		return merr.ErrorParamsNotSupportFileConfig()
	}
	if err := c.emailConfigRepo.UpdateEmailConfigStatus(ctx, req.UID, req.Status); err != nil {
		c.helper.Errorw("msg", "update email config status failed", "error", err, "uid", req.UID)
		return merr.ErrorInternal("update email config status %s failed", req.UID)
	}
	return nil
}

func (c *EmailConfig) DeleteEmailConfig(ctx context.Context, uid snowflake.ID) error {
	if !c.useDatabase {
		return merr.ErrorParamsNotSupportFileConfig()
	}
	if err := c.emailConfigRepo.DeleteEmailConfig(ctx, uid); err != nil {
		c.helper.Errorw("msg", "delete email config failed", "error", err, "uid", uid)
		return merr.ErrorInternal("delete email config %s failed", uid)
	}
	return nil
}

func (c *EmailConfig) getEmailConfigByFileConfigWithUID(ctx context.Context, uid snowflake.ID) (*bo.EmailConfigItemBo, error) {
	namespaceEmailConfigs, ok := c.emailConfigs.Get(middler.GetNamespace(ctx))
	if !ok {
		return nil, merr.ErrorNotFound("email config %s not found", uid)
	}
	emailConfig, ok := namespaceEmailConfigs.Get(uid)
	if !ok {
		return nil, merr.ErrorNotFound("email config %s not found", uid)
	}
	return emailConfig, nil
}

func (c *EmailConfig) GetEmailConfig(ctx context.Context, uid snowflake.ID) (*bo.EmailConfigItemBo, error) {
	if !c.useDatabase {
		return c.getEmailConfigByFileConfigWithUID(ctx, uid)
	}
	doEmailConfig, err := c.emailConfigRepo.GetEmailConfig(ctx, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, merr.ErrorNotFound("email config %s not found", uid)
		}
		c.helper.Errorw("msg", "get email config failed", "error", err, "uid", uid)
		return nil, merr.ErrorInternal("get email config %s failed", uid)
	}
	return bo.NewEmailConfigItemBo(doEmailConfig), nil
}

func (c *EmailConfig) getEmailConfigByFileConfigWithNamespace(ctx context.Context, req *bo.ListEmailConfigBo) ([]*bo.EmailConfigItemBo, error) {
	namespace := middler.GetNamespace(ctx)
	namespaceEmailConfigs, ok := c.emailConfigs.Get(namespace)
	if !ok {
		return nil, merr.ErrorNotFound("email config %s not found", namespace)
	}
	emailConfigs := make([]*bo.EmailConfigItemBo, 0, namespaceEmailConfigs.Len())
	for _, emailConfig := range namespaceEmailConfigs.Values() {
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
	return emailConfigs, nil
}

func (c *EmailConfig) ListEmailConfig(ctx context.Context, req *bo.ListEmailConfigBo) (*bo.PageResponseBo[*bo.EmailConfigItemBo], error) {
	if !c.useDatabase {
		emailConfigs, err := c.getEmailConfigByFileConfigWithNamespace(ctx, req)
		if err != nil {
			return nil, err
		}
		return bo.NewPageResponseBo(req.PageRequestBo, emailConfigs), nil
	}
	pageResponseBo, err := c.emailConfigRepo.ListEmailConfig(ctx, req)
	if err != nil {
		c.helper.Errorw("msg", "list email config failed", "error", err, "req", req)
		return nil, merr.ErrorInternal("list email config failed")
	}
	items := make([]*bo.EmailConfigItemBo, 0, len(pageResponseBo.GetItems()))
	for _, item := range pageResponseBo.GetItems() {
		items = append(items, bo.NewEmailConfigItemBo(item))
	}
	return bo.NewPageResponseBo(pageResponseBo.PageRequestBo, items), nil
}

func (c *EmailConfig) getSelectEmailConfigByFileConfig(ctx context.Context, req *bo.SelectEmailConfigBo) ([]*bo.EmailConfigItemSelectBo, error) {
	namespace := middler.GetNamespace(ctx)
	namespaceEmailConfigs, ok := c.emailConfigs.Get(namespace)
	if !ok {
		return nil, merr.ErrorNotFound("email config %s not found", namespace)
	}
	emailConfigs := make([]*bo.EmailConfigItemSelectBo, 0, namespaceEmailConfigs.Len())
	for _, emailConfig := range namespaceEmailConfigs.Values() {
		if strutil.IsNotEmpty(req.Keyword) && !strings.Contains(emailConfig.Name, req.Keyword) {
			continue
		}
		if req.Status.Exist() && !req.Status.IsUnknown() && emailConfig.Status != req.Status {
			continue
		}
		emailConfigs = append(emailConfigs, &bo.EmailConfigItemSelectBo{
			UID:      emailConfig.UID,
			Name:     emailConfig.Name,
			Status:   emailConfig.Status,
			Disabled: emailConfig.Status != vobj.GlobalStatusEnabled,
			Tooltip:  "",
		})
	}
	total := int64(len(emailConfigs))
	req.Limit = int32(total)
	req.LastUID = 0
	return emailConfigs, nil
}

func (c *EmailConfig) SelectEmailConfig(ctx context.Context, req *bo.SelectEmailConfigBo) (*bo.SelectEmailConfigBoResult, error) {
	if !c.useDatabase {
		emailConfigs, err := c.getSelectEmailConfigByFileConfig(ctx, req)
		if err != nil {
			return nil, err
		}
		return &bo.SelectEmailConfigBoResult{
			Items:   emailConfigs,
			Total:   int64(len(emailConfigs)),
			LastUID: 0,
		}, nil
	}
	result, err := c.emailConfigRepo.SelectEmailConfig(ctx, req)
	if err != nil {
		c.helper.Errorw("msg", "select email config failed", "error", err, "req", req)
		return nil, merr.ErrorInternal("select email config failed")
	}
	items := make([]*bo.EmailConfigItemSelectBo, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, bo.NewEmailConfigItemSelectBo(item))
	}
	return &bo.SelectEmailConfigBoResult{
		Items:   items,
		Total:   result.Total,
		LastUID: result.LastUID,
	}, nil
}
