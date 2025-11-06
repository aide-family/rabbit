package service

import (
	"context"

	"github.com/bwmarrin/snowflake"

	"github.com/aide-family/magicbox/strutil"
	"github.com/aide-family/magicbox/strutil/cnst"
	"github.com/aide-family/rabbit/internal/biz"
	"github.com/aide-family/rabbit/internal/biz/bo"
	apiv1 "github.com/aide-family/rabbit/pkg/api/v1"
	"github.com/aide-family/rabbit/pkg/merr"
	"github.com/aide-family/rabbit/pkg/middler"
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

func (s *NamespaceService) CreateNamespace(ctx context.Context, req *apiv1.CreateNamespaceRequest) (*apiv1.CreateNamespaceReply, error) {
	createNamespaceBo := bo.NewCreateNamespaceBo(req)
	if err := s.namespaceBiz.CreateNamespace(ctx, createNamespaceBo); err != nil {
		return nil, err
	}
	return &apiv1.CreateNamespaceReply{}, nil
}

func (s *NamespaceService) UpdateNamespace(ctx context.Context, req *apiv1.UpdateNamespaceRequest) (*apiv1.UpdateNamespaceReply, error) {
	updateNamespaceBo := bo.NewUpdateNamespaceBo(req)
	if err := s.namespaceBiz.UpdateNamespace(ctx, updateNamespaceBo); err != nil {
		return nil, err
	}
	return &apiv1.UpdateNamespaceReply{}, nil
}

func (s *NamespaceService) UpdateNamespaceStatus(ctx context.Context, req *apiv1.UpdateNamespaceStatusRequest) (*apiv1.UpdateNamespaceStatusReply, error) {
	updateNamespaceStatusBo := bo.NewUpdateNamespaceStatusBo(req)
	if err := s.namespaceBiz.UpdateNamespaceStatus(ctx, updateNamespaceStatusBo); err != nil {
		return nil, err
	}
	return &apiv1.UpdateNamespaceStatusReply{}, nil
}

func (s *NamespaceService) DeleteNamespace(ctx context.Context, req *apiv1.DeleteNamespaceRequest) (*apiv1.DeleteNamespaceReply, error) {
	if err := s.namespaceBiz.DeleteNamespace(ctx, snowflake.ParseInt64(req.Uid)); err != nil {
		return nil, err
	}
	return &apiv1.DeleteNamespaceReply{}, nil
}

func (s *NamespaceService) GetNamespace(ctx context.Context, req *apiv1.GetNamespaceRequest) (*apiv1.NamespaceItem, error) {
	namespaceItemBo, err := s.namespaceBiz.GetNamespace(ctx, snowflake.ParseInt64(req.Uid))
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

func (s *NamespaceService) HasNamespace(ctx context.Context) error {
	ns := middler.GetNamespace(ctx)
	if strutil.IsEmpty(ns) {
		return merr.ErrorForbidden("namespace is required, please set the namespace in the request header or metadata, Example: %s: default", cnst.HTTPHeaderXNamespace)
	}
	namespaceItemBo, err := s.namespaceBiz.GetNamespaceByName(ctx, ns)
	if err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorForbidden("namespace %s not found", ns)
		}
		return err
	}
	if !namespaceItemBo.Status.IsEnabled() {
		return merr.ErrorForbidden("namespace %s is not enabled", ns)
	}
	return nil
}
