// Package fileimpl is the implementation of the namespace repository for file config
package fileimpl

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"

	"github.com/aide-family/magicbox/safety"
	"github.com/aide-family/magicbox/strutil"
	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/do"
	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/internal/biz/vobj"
	"github.com/aide-family/rabbit/internal/conf"
	"github.com/aide-family/rabbit/internal/data"
	"github.com/aide-family/rabbit/pkg/merr"
)

func NewNamespaceRepository(d *data.Data) repository.Namespace {
	n := &namespaceRepositoryImpl{
		d:          d,
		namespaces: d.GetFileConfig().GetNamespaces(),
	}
	n.initNamespaces()
	d.RegisterReloadFunc(data.KeyNamespaces, func() {
		n.initNamespaces()
	})
	return n
}

type namespaceRepositoryImpl struct {
	d                  *data.Data
	namespaces         []*conf.Config_Namespace
	namespacesWithUID  *safety.SyncMap[snowflake.ID, *do.Namespace]
	namespacesWithName *safety.SyncMap[string, *do.Namespace]
}

func (n *namespaceRepositoryImpl) initNamespaces() {
	n.namespaces = n.d.GetFileConfig().GetNamespaces()
	n.namespacesWithUID = safety.NewSyncMap(make(map[snowflake.ID]*do.Namespace))
	n.namespacesWithName = safety.NewSyncMap(make(map[string]*do.Namespace))
	for _, namespace := range n.namespaces {
		uid := snowflake.ParseInt64(namespace.GetUid())
		name := namespace.GetName()
		item := n.toDoNamespace(namespace)
		n.namespacesWithName.Set(name, item)
		n.namespacesWithUID.Set(uid, item)
	}
}

func (n *namespaceRepositoryImpl) toDoNamespace(namespace *conf.Config_Namespace) *do.Namespace {
	createdAt, _ := time.Parse(time.DateTime, namespace.GetCreatedAt())
	updatedAt, _ := time.Parse(time.DateTime, namespace.GetUpdatedAt())
	metadata := safety.NewMap(namespace.GetMetadata())
	return &do.Namespace{
		BaseModel: do.BaseModel{
			ID:        namespace.GetId(),
			UID:       snowflake.ParseInt64(namespace.GetUid()),
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		},
		Name:     namespace.GetName(),
		Metadata: metadata,
		Status:   vobj.GlobalStatus(namespace.GetStatus()),
	}
}

// CreateNamespace implements repository.Namespace.
func (n *namespaceRepositoryImpl) CreateNamespace(ctx context.Context, req *do.Namespace) error {
	return merr.ErrorParamsNotSupportFileConfig()
}

// UpdateNamespace implements repository.Namespace.
func (n *namespaceRepositoryImpl) UpdateNamespace(ctx context.Context, req *do.Namespace) error {
	return merr.ErrorParamsNotSupportFileConfig()
}

// UpdateNamespaceStatus implements repository.Namespace.
func (n *namespaceRepositoryImpl) UpdateNamespaceStatus(ctx context.Context, req *bo.UpdateNamespaceStatusBo) error {
	return merr.ErrorParamsNotSupportFileConfig()
}

// DeleteNamespace implements repository.Namespace.
func (n *namespaceRepositoryImpl) DeleteNamespace(ctx context.Context, uid snowflake.ID) error {
	return merr.ErrorParamsNotSupportFileConfig()
}

// GetNamespace implements repository.Namespace.
func (n *namespaceRepositoryImpl) GetNamespace(ctx context.Context, uid snowflake.ID) (*do.Namespace, error) {
	namespace, ok := n.namespacesWithUID.Get(uid)
	if !ok {
		return nil, merr.ErrorNotFound("namespace %s not found", uid)
	}
	return namespace, nil
}

// GetNamespaceByName implements repository.Namespace.
func (n *namespaceRepositoryImpl) GetNamespaceByName(ctx context.Context, name string) (*do.Namespace, error) {
	namespace, ok := n.namespacesWithName.Get(name)
	if !ok {
		return nil, merr.ErrorNotFound("namespace %s not found", name)
	}
	return namespace, nil
}

// ListNamespace implements repository.Namespace.
func (n *namespaceRepositoryImpl) ListNamespace(ctx context.Context, req *bo.ListNamespaceBo) (*bo.PageResponseBo[*do.Namespace], error) {
	namespaces := make([]*do.Namespace, 0, n.namespacesWithUID.Len())
	for _, namespace := range n.namespacesWithUID.Values() {
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
	sort.Slice(namespaces, func(i, j int) bool {
		return namespaces[i].CreatedAt.After(namespaces[j].CreatedAt)
	})
	return bo.NewPageResponseBo(req.PageRequestBo, namespaces), nil
}

// SelectNamespace implements repository.Namespace.
func (n *namespaceRepositoryImpl) SelectNamespace(ctx context.Context, req *bo.SelectNamespaceBo) (*bo.SelectNamespaceResult, error) {
	namespaces := make([]*do.Namespace, 0, n.namespacesWithUID.Len())
	for _, namespace := range n.namespacesWithUID.Values() {
		if strutil.IsNotEmpty(req.Keyword) && !strings.Contains(namespace.Name, req.Keyword) {
			continue
		}
		if req.Status.Exist() && !req.Status.IsUnknown() && namespace.Status != req.Status {
			continue
		}
		if req.LastUID > 0 && namespace.UID >= req.LastUID {
			continue
		}
		namespaces = append(namespaces, namespace)
	}
	total := int64(len(namespaces))
	sort.Slice(namespaces, func(i, j int) bool {
		return namespaces[i].UID > namespaces[j].UID
	})
	if int32(len(namespaces)) > req.Limit {
		namespaces = namespaces[:req.Limit]
	}
	var lastUID snowflake.ID
	if len(namespaces) > 0 {
		lastUID = namespaces[len(namespaces)-1].UID
	}
	return &bo.SelectNamespaceResult{
		Items:   namespaces,
		Total:   total,
		LastUID: lastUID,
	}, nil
}
