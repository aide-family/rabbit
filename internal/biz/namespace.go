package biz

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/aide-family/magicbox/safety"
	"github.com/aide-family/magicbox/strutil"
	"github.com/bwmarrin/snowflake"
	klog "github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"

	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/internal/biz/vobj"
	"github.com/aide-family/rabbit/internal/conf"
	"github.com/aide-family/rabbit/pkg/merr"
)

func NewNamespace(
	bc *conf.Bootstrap,
	namespaceRepo repository.Namespace,
	helper *klog.Helper,
) *Namespace {
	namespaces := safety.NewSyncMap(make(map[snowflake.ID]*bo.NamespaceItemBo))
	for _, namespace := range conf.GetFileConfig().GetNamespaces() {
		uid := snowflake.ParseInt64(namespace.GetUid())
		createdAt, _ := time.Parse(time.DateTime, namespace.GetCreatedAt())
		updatedAt, _ := time.Parse(time.DateTime, namespace.GetUpdatedAt())
		// Note: Config_Namespace.GetName() returns int64 in proto, convert to string
		name := namespace.GetName()
		namespaces.Set(uid, &bo.NamespaceItemBo{
			UID:       uid,
			Name:      name,
			Metadata:  namespace.GetMetadata(),
			Status:    vobj.GlobalStatus(namespace.GetStatus()),
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
	}
	return &Namespace{
		useDatabase:   bc.GetUseDatabase() == "true",
		namespaceRepo: namespaceRepo,
		namespaces:    namespaces,
		helper:        klog.NewHelper(klog.With(helper.Logger(), "biz", "namespace")),
	}
}

type Namespace struct {
	helper        *klog.Helper
	useDatabase   bool
	namespaceRepo repository.Namespace
	namespaces    *safety.SyncMap[snowflake.ID, *bo.NamespaceItemBo]
}

func (n *Namespace) CreateNamespace(ctx context.Context, req *bo.CreateNamespaceBo) error {
	if !n.useDatabase {
		return merr.ErrorParamsNotSupportFileConfig()
	}
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
	if !n.useDatabase {
		return merr.ErrorParamsNotSupportFileConfig()
	}
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
	if !n.useDatabase {
		return merr.ErrorParamsNotSupportFileConfig()
	}
	if err := n.namespaceRepo.UpdateNamespaceStatus(ctx, req); err != nil {
		n.helper.Errorw("msg", "update namespace status failed", "error", err, "uid", req.UID)
		return merr.ErrorInternal("update namespace status %s failed", req.UID)
	}
	return nil
}

func (n *Namespace) DeleteNamespace(ctx context.Context, uid snowflake.ID) error {
	if !n.useDatabase {
		return merr.ErrorParamsNotSupportFileConfig()
	}
	if err := n.namespaceRepo.DeleteNamespace(ctx, uid); err != nil {
		n.helper.Errorw("msg", "delete namespace failed", "error", err, "uid", uid)
		return merr.ErrorInternal("delete namespace %s failed", uid)
	}
	return nil
}

func (n *Namespace) getNamespaceByFileConfigWithUID(ctx context.Context, uid snowflake.ID) (*bo.NamespaceItemBo, error) {
	namespace, ok := n.namespaces.Get(uid)
	if !ok {
		return nil, merr.ErrorNotFound("namespace %s not found", uid)
	}
	return namespace, nil
}

func (n *Namespace) GetNamespace(ctx context.Context, uid snowflake.ID) (*bo.NamespaceItemBo, error) {
	if !n.useDatabase {
		return n.getNamespaceByFileConfigWithUID(ctx, uid)
	}
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

func (n *Namespace) getNamespaceByFileConfigWithName(ctx context.Context, name string) (*bo.NamespaceItemBo, error) {
	for _, namespace := range n.namespaces.Values() {
		if namespace.Name == name {
			return namespace, nil
		}
	}
	return nil, merr.ErrorNotFound("namespace %s not found", name)
}

func (n *Namespace) GetNamespaceByName(ctx context.Context, name string) (*bo.NamespaceItemBo, error) {
	if !n.useDatabase {
		return n.getNamespaceByFileConfigWithName(ctx, name)
	}
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

func (n *Namespace) getNamespaceByFileConfig(ctx context.Context, req *bo.ListNamespaceBo) ([]*bo.NamespaceItemBo, error) {
	namespaces := make([]*bo.NamespaceItemBo, 0, n.namespaces.Len())
	for _, namespace := range n.namespaces.Values() {
		if strutil.IsNotEmpty(req.Keyword) && !strings.Contains(namespace.Name, req.Keyword) {
			continue
		}
		if req.Status.Exist() && !req.Status.IsUnknown() && namespace.Status != req.Status {
			continue
		}
		namespaces = append(namespaces, namespace)
	}
	total := int64(len(namespaces))
	pageRequestBo := bo.NewPageRequestBo(req.Page, req.PageSize)
	pageRequestBo.WithTotal(total)
	req.PageRequestBo = pageRequestBo
	return namespaces, nil
}

func (n *Namespace) ListNamespace(ctx context.Context, req *bo.ListNamespaceBo) (*bo.PageResponseBo[*bo.NamespaceItemBo], error) {
	if !n.useDatabase {
		namespaces, err := n.getNamespaceByFileConfig(ctx, req)
		if err != nil {
			return nil, err
		}
		return bo.NewPageResponseBo(req.PageRequestBo, namespaces), nil
	}
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

func (n *Namespace) getSelectNamespaceByFileConfig(ctx context.Context, req *bo.SelectNamespaceBo) ([]*bo.NamespaceItemSelectBo, error) {
	namespaces := make([]*bo.NamespaceItemSelectBo, 0, n.namespaces.Len())
	for _, namespace := range n.namespaces.Values() {
		if strutil.IsNotEmpty(req.Keyword) && !strings.Contains(namespace.Name, req.Keyword) {
			continue
		}
		if req.Status.Exist() && !req.Status.IsUnknown() && namespace.Status != req.Status {
			continue
		}
		namespaces = append(namespaces, &bo.NamespaceItemSelectBo{
			UID:      namespace.UID,
			Name:     namespace.Name,
			Status:   namespace.Status,
			Disabled: namespace.Status != vobj.GlobalStatusEnabled,
			Tooltip:  "",
		})
	}
	total := int64(len(namespaces))
	req.Limit = int32(total)
	req.LastUID = 0
	return namespaces, nil
}

func (n *Namespace) SelectNamespace(ctx context.Context, req *bo.SelectNamespaceBo) (*bo.SelectNamespaceBoResult, error) {
	if !n.useDatabase {
		namespaces, err := n.getSelectNamespaceByFileConfig(ctx, req)
		if err != nil {
			return nil, err
		}
		return &bo.SelectNamespaceBoResult{
			Items:   namespaces,
			Total:   int64(len(namespaces)),
			LastUID: 0,
		}, nil
	}
	result, err := n.namespaceRepo.SelectNamespace(ctx, req)
	if err != nil {
		n.helper.Errorw("msg", "select namespace failed", "error", err, "req", req)
		return nil, merr.ErrorInternal("select namespace failed")
	}
	items := make([]*bo.NamespaceItemSelectBo, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, bo.NewNamespaceItemSelectBo(item))
	}
	return &bo.SelectNamespaceBoResult{
		Items:   items,
		Total:   result.Total,
		LastUID: result.LastUID,
	}, nil
}
