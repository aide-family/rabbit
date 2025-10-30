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
	GetEmailConfigByName(ctx context.Context, name string) (*do.EmailConfig, error)
	ListEmailConfig(ctx context.Context, req *bo.ListEmailConfigBo) (*bo.PageResponseBo[*do.EmailConfig], error)
}
