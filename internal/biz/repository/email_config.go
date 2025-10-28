package repository

import (
	"context"

	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/do"
	"github.com/aide-family/rabbit/internal/biz/vobj"
)

type EmailConfig interface {
	SaveEmailConfig(ctx context.Context, req *do.EmailConfig) error
	UpdateEmailConfigStatus(ctx context.Context, uid string, status vobj.GlobalStatus) error
	DeleteEmailConfig(ctx context.Context, uid string) error
	GetEmailConfig(ctx context.Context, uid string) (*do.EmailConfig, error)
	ListEmailConfig(ctx context.Context, req *bo.ListEmailConfigBo) (*bo.PageResponseBo[*do.EmailConfig], error)
}

type EmailTemplate interface {
	SaveEmailTemplate(ctx context.Context, req *do.EmailTemplate) error
	UpdateEmailTemplateStatus(ctx context.Context, uid string, status vobj.GlobalStatus) error
	DeleteEmailTemplate(ctx context.Context, uid string) error
	GetEmailTemplate(ctx context.Context, uid string) (*do.EmailTemplate, error)
	ListEmailTemplate(ctx context.Context, req *bo.ListEmailTemplateBo) (*bo.PageResponseBo[*do.EmailTemplate], error)
}
