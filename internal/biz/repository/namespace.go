package repository

import (
	"context"

	"github.com/bwmarrin/snowflake"

	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/do"
)

type Namespace interface {
	CreateNamespace(ctx context.Context, req *do.Namespace) error
	UpdateNamespace(ctx context.Context, req *do.Namespace) error
	UpdateNamespaceStatus(ctx context.Context, req *bo.UpdateNamespaceStatusBo) error
	DeleteNamespace(ctx context.Context, uid snowflake.ID) error
	GetNamespace(ctx context.Context, uid snowflake.ID) (*do.Namespace, error)
	GetNamespaceByName(ctx context.Context, name string) (*do.Namespace, error)
	ListNamespace(ctx context.Context, req *bo.ListNamespaceBo) (*bo.PageResponseBo[*do.Namespace], error)
	SelectNamespace(ctx context.Context, req *bo.SelectNamespaceBo) (*bo.SelectNamespaceResult, error)
}
