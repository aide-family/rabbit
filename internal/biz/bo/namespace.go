package bo

import (
	"time"

	"github.com/bwmarrin/snowflake"

	"github.com/aide-family/magicbox/safety"
	"github.com/aide-family/rabbit/internal/biz/do"
	"github.com/aide-family/rabbit/internal/biz/vobj"
	apiv1 "github.com/aide-family/rabbit/pkg/api/v1"
	"github.com/aide-family/rabbit/pkg/enum"
)

type CreateNamespaceBo struct {
	Name     string
	Metadata map[string]string
}

func (b *CreateNamespaceBo) ToDoNamespace() *do.Namespace {
	return &do.Namespace{
		Name:     b.Name,
		Metadata: safety.NewMap(b.Metadata),
	}
}

func NewCreateNamespaceBo(req *apiv1.CreateNamespaceRequest) *CreateNamespaceBo {
	return &CreateNamespaceBo{
		Name:     req.Name,
		Metadata: req.Metadata,
	}
}

type UpdateNamespaceBo struct {
	UID      snowflake.ID
	Name     string
	Metadata map[string]string
}

func (b *UpdateNamespaceBo) ToDoNamespace() *do.Namespace {
	namespace := &do.Namespace{
		Name:     b.Name,
		Metadata: safety.NewMap(b.Metadata),
	}
	namespace.WithUID(b.UID)
	return namespace
}

func NewUpdateNamespaceBo(req *apiv1.UpdateNamespaceRequest) *UpdateNamespaceBo {
	return &UpdateNamespaceBo{
		UID:      snowflake.ParseInt64(req.Uid),
		Name:     req.Name,
		Metadata: req.Metadata,
	}
}

type UpdateNamespaceStatusBo struct {
	UID    snowflake.ID
	Status vobj.GlobalStatus
}

type ListNamespaceBo struct {
	*PageRequestBo
	Keyword string
	Status  vobj.GlobalStatus
}

type NamespaceItemBo struct {
	UID       snowflake.ID
	Name      string
	Metadata  map[string]string
	Status    vobj.GlobalStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewNamespaceItemBo(doNamespace *do.Namespace) *NamespaceItemBo {
	return &NamespaceItemBo{
		UID:       doNamespace.UID,
		Name:      doNamespace.Name,
		Metadata:  doNamespace.Metadata.Map(),
		Status:    doNamespace.Status,
		CreatedAt: doNamespace.CreatedAt,
		UpdatedAt: doNamespace.UpdatedAt,
	}
}

func (b *NamespaceItemBo) ToAPIV1NamespaceItem() *apiv1.NamespaceItem {
	return &apiv1.NamespaceItem{
		Uid:       b.UID.Int64(),
		Name:      b.Name,
		Metadata:  b.Metadata,
		Status:    enum.GlobalStatus(b.Status),
		CreatedAt: b.CreatedAt.Format(time.DateTime),
		UpdatedAt: b.UpdatedAt.Format(time.DateTime),
	}
}

func NewUpdateNamespaceStatusBo(req *apiv1.UpdateNamespaceStatusRequest) *UpdateNamespaceStatusBo {
	return &UpdateNamespaceStatusBo{
		UID:    snowflake.ParseInt64(req.Uid),
		Status: vobj.GlobalStatus(req.Status),
	}
}

func NewListNamespaceBo(req *apiv1.ListNamespaceRequest) *ListNamespaceBo {
	return &ListNamespaceBo{
		PageRequestBo: NewPageRequestBo(req.Page, req.PageSize),
		Keyword:       req.Keyword,
		Status:        vobj.GlobalStatus(req.Status),
	}
}

func ToAPIV1ListNamespaceReply(pageResponseBo *PageResponseBo[*NamespaceItemBo]) *apiv1.ListNamespaceReply {
	items := make([]*apiv1.NamespaceItem, 0, len(pageResponseBo.GetItems()))
	for _, item := range pageResponseBo.GetItems() {
		items = append(items, item.ToAPIV1NamespaceItem())
	}
	return &apiv1.ListNamespaceReply{
		Items:    items,
		Total:    pageResponseBo.GetTotal(),
		Page:     pageResponseBo.GetPage(),
		PageSize: pageResponseBo.GetPageSize(),
	}
}
