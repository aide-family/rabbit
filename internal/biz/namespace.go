package biz

import (
	"context"
	"errors"

	"github.com/bwmarrin/snowflake"
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

func (n *Namespace) CreateNamespace(ctx context.Context, req *bo.CreateNamespaceBo) error {
	if _, err := n.namespaceRepo.GetNamespaceByName(ctx, req.Name); err == nil {
		return merr.ErrorParams("namespace %s already exists", req.Name)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		n.helper.Errorw("msg", "check namespace exists failed", "error", err, "name", req.Name)
		return merr.ErrorInternal("create namespace %s failed", req.Name)
	}
	if err := n.namespaceRepo.CreateNamespace(ctx, req.ToDoNamespace()); err != nil {
		n.helper.Errorw("msg", "create namespace failed", "error", err, "name", req.Name)
		return merr.ErrorInternal("create namespace %s failed", req.Name)
	}
	return nil
}

func (n *Namespace) UpdateNamespace(ctx context.Context, req *bo.UpdateNamespaceBo) error {
	doNamespace := req.ToDoNamespace()
	existNamespace, err := n.namespaceRepo.GetNamespaceByName(ctx, doNamespace.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		n.helper.Errorw("msg", "check namespace exists failed", "error", err, "name", doNamespace.Name)
		return merr.ErrorInternal("update namespace %s failed", doNamespace.Name)
	} else if existNamespace != nil && existNamespace.UID != doNamespace.UID {
		return merr.ErrorParams("namespace %s already exists", doNamespace.Name)
	}
	if err := n.namespaceRepo.UpdateNamespace(ctx, doNamespace); err != nil {
		n.helper.Errorw("msg", "update namespace failed", "error", err, "uid", req.UID)
		return merr.ErrorInternal("update namespace %s failed", req.UID)
	}
	return nil
}

func (n *Namespace) UpdateNamespaceStatus(ctx context.Context, req *bo.UpdateNamespaceStatusBo) error {
	if err := n.namespaceRepo.UpdateNamespaceStatus(ctx, req); err != nil {
		n.helper.Errorw("msg", "update namespace status failed", "error", err, "uid", req.UID)
		return merr.ErrorInternal("update namespace status %s failed", req.UID)
	}
	return nil
}

func (n *Namespace) DeleteNamespace(ctx context.Context, uid snowflake.ID) error {
	if err := n.namespaceRepo.DeleteNamespace(ctx, uid); err != nil {
		n.helper.Errorw("msg", "delete namespace failed", "error", err, "uid", uid)
		return merr.ErrorInternal("delete namespace %s failed", uid)
	}
	return nil
}

func (n *Namespace) GetNamespace(ctx context.Context, uid snowflake.ID) (*bo.NamespaceItemBo, error) {
	doNamespace, err := n.namespaceRepo.GetNamespace(ctx, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, merr.ErrorNotFound("namespace %s not found", uid)
		}

		n.helper.Errorw("msg", "get namespace failed", "error", err, "uid", uid)
		return nil, merr.ErrorInternal("get namespace %s failed", uid)
	}
	return bo.NewNamespaceItemBo(doNamespace), nil
}

func (n *Namespace) GetNamespaceByName(ctx context.Context, name string) (*bo.NamespaceItemBo, error) {
	doNamespace, err := n.namespaceRepo.GetNamespaceByName(ctx, name)
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
