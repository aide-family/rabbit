// Package namespacev1 is the namespace service implementation.
package namespacev1

import (
	context "context"
)

type Repository interface {
	CreateNamespace(ctx context.Context, req *CreateNamespaceRequest) (*NamespaceModel, error)
	GetNamespace(ctx context.Context, req *GetNamespaceRequest) (*NamespaceModel, error)
	UpdateNamespace(ctx context.Context, req *UpdateNamespaceRequest) (*ResultInfo, error)
	DeleteNamespace(ctx context.Context, req *DeleteNamespaceRequest) (*ResultInfo, error)
	ListNamespace(ctx context.Context, req *ListNamespaceRequest) (*ListNamespaceResponse, error)
	SelectNamespace(ctx context.Context, req *SelectNamespaceRequest) (*SelectNamespaceResponse, error)
	UpdateNamespaceStatus(ctx context.Context, req *UpdateNamespaceStatusRequest) (*ResultInfo, error)
	GetNamespaceByName(ctx context.Context, req *GetNamespaceByNameRequest) (*NamespaceModel, error)
}
