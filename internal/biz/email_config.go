package biz

import (
	"context"

	"github.com/aide-family/rabbit/internal/biz/bo"
)

func NewEmailConfig() *EmailConfig {
	return &EmailConfig{}
}

type EmailConfig struct{}

func (c *EmailConfig) CreateEmailConfig(ctx context.Context, req *bo.CreateEmailConfigBo) error {
	return nil
}

func (c *EmailConfig) UpdateEmailConfig(ctx context.Context, req *bo.UpdateEmailConfigBo) error {
	return nil
}

func (c *EmailConfig) UpdateEmailConfigStatus(ctx context.Context, req *bo.UpdateEmailConfigStatusBo) error {
	return nil
}

func (c *EmailConfig) DeleteEmailConfig(ctx context.Context, uid string) error {
	return nil
}

func (c *EmailConfig) GetEmailConfig(ctx context.Context, uid string) (*bo.EmailConfigItemBo, error) {
	return nil, nil
}

func (c *EmailConfig) ListEmailConfig(ctx context.Context, req *bo.ListEmailConfigBo) (*bo.PageResponseBo[*bo.EmailConfigItemBo], error) {
	return nil, nil
}

func (c *EmailConfig) CreateEmailTemplate(ctx context.Context, req *bo.CreateEmailTemplateBo) error {
	return nil
}

func (c *EmailConfig) UpdateEmailTemplate(ctx context.Context, req *bo.UpdateEmailTemplateBo) error {
	return nil
}

func (c *EmailConfig) UpdateEmailTemplateStatus(ctx context.Context, req *bo.UpdateEmailTemplateStatusBo) error {
	return nil
}

func (c *EmailConfig) DeleteEmailTemplate(ctx context.Context, uid string) error {
	return nil
}

func (c *EmailConfig) GetEmailTemplate(ctx context.Context, uid string) (*bo.EmailTemplateItemBo, error) {
	return nil, nil
}

func (c *EmailConfig) ListEmailTemplate(ctx context.Context, req *bo.ListEmailTemplateBo) (*bo.PageResponseBo[*bo.EmailTemplateItemBo], error) {
	return nil, nil
}
