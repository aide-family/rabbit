package repository

import (
	"context"

	"github.com/bwmarrin/snowflake"

	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/do"
	"github.com/aide-family/rabbit/internal/biz/vobj"
)

type Template interface {
	CreateTemplate(ctx context.Context, req *do.Template) error
	UpdateTemplate(ctx context.Context, req *do.Template) error
	UpdateTemplateStatus(ctx context.Context, uid snowflake.ID, status vobj.GlobalStatus) error
	DeleteTemplate(ctx context.Context, uid snowflake.ID) error
	GetTemplate(ctx context.Context, uid snowflake.ID) (*do.Template, error)
	GetTemplateByName(ctx context.Context, name string) (*do.Template, error)
	ListTemplate(ctx context.Context, req *bo.ListTemplateBo) (*bo.PageResponseBo[*do.Template], error)
	SelectTemplate(ctx context.Context, req *bo.SelectTemplateBo) (*bo.SelectTemplateResult, error)
}
