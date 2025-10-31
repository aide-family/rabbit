package biz

import (
	"context"
	"errors"

	"github.com/bwmarrin/snowflake"
	klog "github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"

	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/pkg/merr"
)

func NewEmailConfig(
	emailConfigRepo repository.EmailConfig,
	helper *klog.Helper,
) *EmailConfig {
	return &EmailConfig{
		emailConfigRepo: emailConfigRepo,
		helper:          klog.NewHelper(klog.With(helper.Logger(), "biz", "email_config")),
	}
}

type EmailConfig struct {
	helper          *klog.Helper
	emailConfigRepo repository.EmailConfig
}

func (c *EmailConfig) CreateEmailConfig(ctx context.Context, req *bo.CreateEmailConfigBo) error {
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
	doEmailConfig := req.ToDoEmailConfig()
	if err := c.emailConfigRepo.UpdateEmailConfig(ctx, doEmailConfig); err != nil {
		c.helper.Errorw("msg", "update email config failed", "error", err, "name", doEmailConfig.Name)
		return merr.ErrorInternal("update email config %s failed", doEmailConfig.Name)
	}
	return nil
}

func (c *EmailConfig) UpdateEmailConfigStatus(ctx context.Context, req *bo.UpdateEmailConfigStatusBo) error {
	if err := c.emailConfigRepo.UpdateEmailConfigStatus(ctx, req.UID, req.Status); err != nil {
		c.helper.Errorw("msg", "update email config status failed", "error", err, "uid", req.UID)
		return merr.ErrorInternal("update email config status %s failed", req.UID)
	}
	return nil
}

func (c *EmailConfig) DeleteEmailConfig(ctx context.Context, uid snowflake.ID) error {
	if err := c.emailConfigRepo.DeleteEmailConfig(ctx, uid); err != nil {
		c.helper.Errorw("msg", "delete email config failed", "error", err, "uid", uid)
		return merr.ErrorInternal("delete email config %s failed", uid)
	}
	return nil
}

func (c *EmailConfig) GetEmailConfig(ctx context.Context, uid snowflake.ID) (*bo.EmailConfigItemBo, error) {
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

func (c *EmailConfig) ListEmailConfig(ctx context.Context, req *bo.ListEmailConfigBo) (*bo.PageResponseBo[*bo.EmailConfigItemBo], error) {
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
