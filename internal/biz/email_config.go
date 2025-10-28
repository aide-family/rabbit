package biz

import (
	"context"

	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/repository"
)

func NewEmailConfig(emailConfigRepository repository.EmailConfig, emailTemplateRepository repository.EmailTemplate) *EmailConfig {
	return &EmailConfig{
		emailConfigRepository:   emailConfigRepository,
		emailTemplateRepository: emailTemplateRepository,
	}
}

type EmailConfig struct {
	emailConfigRepository   repository.EmailConfig
	emailTemplateRepository repository.EmailTemplate
}

func (c *EmailConfig) CreateEmailConfig(ctx context.Context, req *bo.CreateEmailConfigBo) error {
	doEmailConfig := req.ToDoEmailConfig()
	return c.emailConfigRepository.SaveEmailConfig(ctx, doEmailConfig)
}

func (c *EmailConfig) UpdateEmailConfig(ctx context.Context, req *bo.UpdateEmailConfigBo) error {
	doEmailConfig := req.ToDoEmailConfig()
	return c.emailConfigRepository.SaveEmailConfig(ctx, doEmailConfig)
}

func (c *EmailConfig) UpdateEmailConfigStatus(ctx context.Context, req *bo.UpdateEmailConfigStatusBo) error {
	return c.emailConfigRepository.UpdateEmailConfigStatus(ctx, req.UID, req.Status)
}

func (c *EmailConfig) DeleteEmailConfig(ctx context.Context, uid string) error {
	return c.emailConfigRepository.DeleteEmailConfig(ctx, uid)
}

func (c *EmailConfig) GetEmailConfig(ctx context.Context, uid string) (*bo.EmailConfigItemBo, error) {
	doEmailConfig, err := c.emailConfigRepository.GetEmailConfig(ctx, uid)
	if err != nil {
		return nil, err
	}
	return bo.NewEmailConfigItemBo(doEmailConfig), nil
}

func (c *EmailConfig) ListEmailConfig(ctx context.Context, req *bo.ListEmailConfigBo) (*bo.PageResponseBo[*bo.EmailConfigItemBo], error) {
	pageResponseBo, err := c.emailConfigRepository.ListEmailConfig(ctx, req)
	if err != nil {
		return nil, err
	}
	items := make([]*bo.EmailConfigItemBo, 0, len(pageResponseBo.GetItems()))
	for _, item := range pageResponseBo.GetItems() {
		items = append(items, bo.NewEmailConfigItemBo(item))
	}
	return bo.NewPageResponseBo(pageResponseBo.PageRequestBo, items), nil
}

func (c *EmailConfig) CreateEmailTemplate(ctx context.Context, req *bo.CreateEmailTemplateBo) error {
	doEmailTemplate := req.ToDoEmailTemplate()
	return c.emailTemplateRepository.SaveEmailTemplate(ctx, doEmailTemplate)
}

func (c *EmailConfig) UpdateEmailTemplate(ctx context.Context, req *bo.UpdateEmailTemplateBo) error {
	doEmailTemplate := req.ToDoEmailTemplate()
	return c.emailTemplateRepository.SaveEmailTemplate(ctx, doEmailTemplate)
}

func (c *EmailConfig) UpdateEmailTemplateStatus(ctx context.Context, req *bo.UpdateEmailTemplateStatusBo) error {
	return c.emailTemplateRepository.UpdateEmailTemplateStatus(ctx, req.UID, req.Status)
}

func (c *EmailConfig) DeleteEmailTemplate(ctx context.Context, uid string) error {
	return c.emailTemplateRepository.DeleteEmailTemplate(ctx, uid)
}

func (c *EmailConfig) GetEmailTemplate(ctx context.Context, uid string) (*bo.EmailTemplateItemBo, error) {
	doEmailTemplate, err := c.emailTemplateRepository.GetEmailTemplate(ctx, uid)
	if err != nil {
		return nil, err
	}
	return bo.NewEmailTemplateItemBo(doEmailTemplate), nil
}

func (c *EmailConfig) ListEmailTemplate(ctx context.Context, req *bo.ListEmailTemplateBo) (*bo.PageResponseBo[*bo.EmailTemplateItemBo], error) {
	pageResponseBo, err := c.emailTemplateRepository.ListEmailTemplate(ctx, req)
	if err != nil {
		return nil, err
	}
	items := make([]*bo.EmailTemplateItemBo, 0, len(pageResponseBo.GetItems()))
	for _, item := range pageResponseBo.GetItems() {
		items = append(items, bo.NewEmailTemplateItemBo(item))
	}
	return bo.NewPageResponseBo(pageResponseBo.PageRequestBo, items), nil
}
