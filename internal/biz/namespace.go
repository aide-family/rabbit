package biz

import (
	"context"
	"errors"

	klog "github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"

	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/pkg/merr"
)

func NewNamespace(namespaceRepo repository.Namespace, helper *klog.Helper) *Namespace {
	return &Namespace{
		namespaceRepo: namespaceRepo,
		helper:        klog.NewHelper(klog.With(helper.Logger(), "biz", "namespace")),
	}
}

type Namespace struct {
	helper        *klog.Helper
	namespaceRepo repository.Namespace
}

func (n *Namespace) SaveNamespace(ctx context.Context, req *bo.SaveNamespaceBo) error {
	doNamespace := req.ToDoNamespace()
	if err := n.namespaceRepo.SaveNamespace(ctx, doNamespace); err != nil {
		n.helper.Errorw("msg", "save namespace failed", "error", err, "name", doNamespace.Name)
		return merr.ErrorInternal("save namespace %s failed", doNamespace.Name)
	}
	return nil
}

func (n *Namespace) UpdateNamespaceStatus(ctx context.Context, req *bo.UpdateNamespaceStatusBo) error {
	if err := n.namespaceRepo.UpdateNamespaceStatus(ctx, req); err != nil {
		n.helper.Errorw("msg", "update namespace status failed", "error", err, "name", req.Name)
		return merr.ErrorInternal("update namespace status %s failed", req.Name)
	}
	return nil
}

func (n *Namespace) DeleteNamespace(ctx context.Context, name string) error {
	if err := n.namespaceRepo.DeleteNamespace(ctx, name); err != nil {
		n.helper.Errorw("msg", "delete namespace failed", "error", err, "name", name)
		return merr.ErrorInternal("delete namespace %s failed", name)
	}
	return nil
}

func (n *Namespace) GetNamespace(ctx context.Context, name string) (*bo.NamespaceItemBo, error) {
	doNamespace, err := n.namespaceRepo.GetNamespace(ctx, name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, merr.ErrorNotFound("namespace %s not found", name)
		}

		n.helper.Errorw("msg", "get namespace failed", "error", err, "name", name)
		return nil, merr.ErrorInternal("get namespace %s failed", name)
	}
	return bo.NewNamespaceItemBo(doNamespace), nil
}

func (n *Namespace) ListNamespace(ctx context.Context, req *bo.ListNamespaceBo) (*bo.PageResponseBo[*bo.NamespaceItemBo], error) {
	pageResponseBo, err := n.namespaceRepo.ListNamespace(ctx, req)
	if err != nil {
		n.helper.Errorw("msg", "list namespace failed", "error", err, "req", req)
		return nil, merr.ErrorInternal("list namespace failed")
	}
	items := make([]*bo.NamespaceItemBo, 0, len(pageResponseBo.GetItems()))
	for _, item := range pageResponseBo.GetItems() {
		items = append(items, bo.NewNamespaceItemBo(item))
	}
	return bo.NewPageResponseBo(pageResponseBo.PageRequestBo, items), nil
}
