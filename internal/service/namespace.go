package service

import (
	"context"

	"github.com/aide-family/rabbit/internal/biz"
	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/vobj"
	apiv1 "github.com/aide-family/rabbit/pkg/api/v1"
)

func NewNamespaceService(namespaceBiz *biz.Namespace) *NamespaceService {
	return &NamespaceService{
		namespaceBiz: namespaceBiz,
	}
}

type NamespaceService struct {
	apiv1.UnimplementedNamespaceServer

	namespaceBiz *biz.Namespace
}

func (s *NamespaceService) SaveNamespace(ctx context.Context, req *apiv1.SaveNamespaceRequest) (*apiv1.SaveNamespaceReply, error) {
	saveNamespaceBo := bo.NewSaveNamespaceBo(req.Name, req.Metadata)
	if err := s.namespaceBiz.SaveNamespace(ctx, saveNamespaceBo); err != nil {
		return nil, err
	}
	return &apiv1.SaveNamespaceReply{}, nil
}

func (s *NamespaceService) UpdateNamespaceStatus(ctx context.Context, req *apiv1.UpdateNamespaceStatusRequest) (*apiv1.UpdateNamespaceStatusReply, error) {
	updateNamespaceStatusBo := bo.NewUpdateNamespaceStatusBo(req.Name, vobj.GlobalStatus(req.Status))
	if err := s.namespaceBiz.UpdateNamespaceStatus(ctx, updateNamespaceStatusBo); err != nil {
		return nil, err
	}
	return &apiv1.UpdateNamespaceStatusReply{}, nil
}

func (s *NamespaceService) DeleteNamespace(ctx context.Context, req *apiv1.DeleteNamespaceRequest) (*apiv1.DeleteNamespaceReply, error) {
	if err := s.namespaceBiz.DeleteNamespace(ctx, req.Name); err != nil {
		return nil, err
	}
	return &apiv1.DeleteNamespaceReply{}, nil
}

func (s *NamespaceService) GetNamespace(ctx context.Context, req *apiv1.GetNamespaceRequest) (*apiv1.NamespaceItem, error) {
	namespaceItemBo, err := s.namespaceBiz.GetNamespace(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	return namespaceItemBo.ToAPIV1NamespaceItem(), nil
}

func (s *NamespaceService) ListNamespace(ctx context.Context, req *apiv1.ListNamespaceRequest) (*apiv1.ListNamespaceReply, error) {
	listNamespaceBo := bo.NewListNamespaceBo(req)
	listNamespacePageResponseBo, err := s.namespaceBiz.ListNamespace(ctx, listNamespaceBo)
	if err != nil {
		return nil, err
	}
	return bo.ToAPIV1ListNamespaceReply(listNamespacePageResponseBo), nil
}
