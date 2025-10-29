package repository

import (
	"context"

	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/do"
)

type Namespace interface {
	SaveNamespace(ctx context.Context, req *do.Namespace) error
	UpdateNamespaceStatus(ctx context.Context, req *bo.UpdateNamespaceStatusBo) error
	DeleteNamespace(ctx context.Context, name string) error
	GetNamespace(ctx context.Context, name string) (*do.Namespace, error)
	ListNamespace(ctx context.Context, req *bo.ListNamespaceBo) (*bo.PageResponseBo[*do.Namespace], error)
}
