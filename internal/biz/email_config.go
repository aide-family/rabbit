package biz

import (
	"context"
	"errors"

	klog "github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"

	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/pkg/merr"
)

func NewEmailConfig(
	emailConfigRepo repository.EmailConfig,
	emailTemplateRepo repository.EmailTemplate,
	namespace string,
	helper *klog.Helper,
) *EmailConfig {
	return &EmailConfig{
		emailConfigRepo:   emailConfigRepo,
		emailTemplateRepo: emailTemplateRepo,
		helper:            klog.NewHelper(klog.With(helper.Logger(), "biz", "email_config")),
	}
}

type EmailConfig struct {
	helper            *klog.Helper
	emailConfigRepo   repository.EmailConfig
	emailTemplateRepo repository.EmailTemplate
}

func (c *EmailConfig) CreateEmailConfig(ctx context.Context, req *bo.CreateEmailConfigBo) error {
	doEmailConfig := req.ToDoEmailConfig()
	if err := c.emailConfigRepo.SaveEmailConfig(ctx, doEmailConfig); err != nil {
		c.helper.Errorw("msg", "create email config failed", "error", err, "name", doEmailConfig.Name)
		return merr.ErrorInternal("create email config %s failed", doEmailConfig.Name)
	}
	return nil
}

func (c *EmailConfig) UpdateEmailConfig(ctx context.Context, req *bo.UpdateEmailConfigBo) error {
	doEmailConfig := req.ToDoEmailConfig()
	if err := c.emailConfigRepo.SaveEmailConfig(ctx, doEmailConfig); err != nil {
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

func (c *EmailConfig) DeleteEmailConfig(ctx context.Context, uid string) error {
	if err := c.emailConfigRepo.DeleteEmailConfig(ctx, uid); err != nil {
		c.helper.Errorw("msg", "delete email config failed", "error", err, "uid", uid)
		return merr.ErrorInternal("delete email config %s failed", uid)
	}
	return nil
}

func (c *EmailConfig) GetEmailConfig(ctx context.Context, uid string) (*bo.EmailConfigItemBo, error) {
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

func (c *EmailConfig) CreateEmailTemplate(ctx context.Context, req *bo.CreateEmailTemplateBo) error {
	doEmailTemplate := req.ToDoEmailTemplate()
	if err := c.emailTemplateRepo.SaveEmailTemplate(ctx, doEmailTemplate); err != nil {
		c.helper.Errorw("msg", "create email template failed", "error", err, "name", doEmailTemplate.Name)
		return merr.ErrorInternal("create email template %s failed", doEmailTemplate.Name)
	}
	return nil
}

func (c *EmailConfig) UpdateEmailTemplate(ctx context.Context, req *bo.UpdateEmailTemplateBo) error {
	doEmailTemplate := req.ToDoEmailTemplate()
	if err := c.emailTemplateRepo.SaveEmailTemplate(ctx, doEmailTemplate); err != nil {
		c.helper.Errorw("msg", "update email template failed", "error", err, "name", doEmailTemplate.Name)
		return merr.ErrorInternal("update email template %s failed", doEmailTemplate.Name)
	}
	return nil
}

func (c *EmailConfig) UpdateEmailTemplateStatus(ctx context.Context, req *bo.UpdateEmailTemplateStatusBo) error {
	if err := c.emailTemplateRepo.UpdateEmailTemplateStatus(ctx, req.UID, req.Status); err != nil {
		c.helper.Errorw("msg", "update email template status failed", "error", err, "uid", req.UID)
		return merr.ErrorInternal("update email template status %s failed", req.UID)
	}
	return nil
}

func (c *EmailConfig) DeleteEmailTemplate(ctx context.Context, uid string) error {
	if err := c.emailTemplateRepo.DeleteEmailTemplate(ctx, uid); err != nil {
		c.helper.Errorw("msg", "delete email template failed", "error", err, "uid", uid)
		return merr.ErrorInternal("delete email template %s failed", uid)
	}
	return nil
}

func (c *EmailConfig) GetEmailTemplate(ctx context.Context, uid string) (*bo.EmailTemplateItemBo, error) {
	doEmailTemplate, err := c.emailTemplateRepo.GetEmailTemplate(ctx, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, merr.ErrorNotFound("email template %s not found", uid)
		}
		c.helper.Errorw("msg", "get email template failed", "error", err, "uid", uid)
		return nil, merr.ErrorInternal("get email template %s failed", uid)
	}
	return bo.NewEmailTemplateItemBo(doEmailTemplate), nil
}

func (c *EmailConfig) ListEmailTemplate(ctx context.Context, req *bo.ListEmailTemplateBo) (*bo.PageResponseBo[*bo.EmailTemplateItemBo], error) {
	pageResponseBo, err := c.emailTemplateRepo.ListEmailTemplate(ctx, req)
	if err != nil {
		c.helper.Errorw("msg", "list email template failed", "error", err, "req", req)
		return nil, merr.ErrorInternal("list email template failed")
	}
	items := make([]*bo.EmailTemplateItemBo, 0, len(pageResponseBo.GetItems()))
	for _, item := range pageResponseBo.GetItems() {
		items = append(items, bo.NewEmailTemplateItemBo(item))
	}
	return bo.NewPageResponseBo(pageResponseBo.PageRequestBo, items), nil
}
