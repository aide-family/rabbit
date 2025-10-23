package service

import (
	"context"

	apiv1 "github.com/aide-family/rabbit/pkg/api/v1"
)

type NamespaceService struct {
	apiv1.UnimplementedNamespaceServer
}

func NewNamespaceService() *NamespaceService {
	return &NamespaceService{}
}

func (s *NamespaceService) CreateNamespace(ctx context.Context, req *apiv1.CreateNamespaceRequest) (*apiv1.CreateNamespaceReply, error) {
	return &apiv1.CreateNamespaceReply{}, nil
}

func (s *NamespaceService) UpdateNamespace(ctx context.Context, req *apiv1.UpdateNamespaceRequest) (*apiv1.UpdateNamespaceReply, error) {
	return &apiv1.UpdateNamespaceReply{}, nil
}

func (s *NamespaceService) UpdateNamespaceStatus(ctx context.Context, req *apiv1.UpdateNamespaceStatusRequest) (*apiv1.UpdateNamespaceStatusReply, error) {
	return &apiv1.UpdateNamespaceStatusReply{}, nil
}

func (s *NamespaceService) DeleteNamespace(ctx context.Context, req *apiv1.DeleteNamespaceRequest) (*apiv1.DeleteNamespaceReply, error) {
	return &apiv1.DeleteNamespaceReply{}, nil
}

func (s *NamespaceService) GetNamespace(ctx context.Context, req *apiv1.GetNamespaceRequest) (*apiv1.GetNamespaceReply, error) {
	return &apiv1.GetNamespaceReply{}, nil
}

func (s *NamespaceService) ListNamespace(ctx context.Context, req *apiv1.ListNamespaceRequest) (*apiv1.ListNamespaceReply, error) {
	return &apiv1.ListNamespaceReply{}, nil
}
