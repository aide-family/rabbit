package bo

import (
	"time"

	"github.com/aide-family/rabbit/internal/biz/do"
	"github.com/aide-family/rabbit/internal/biz/vobj"
	apiv1 "github.com/aide-family/rabbit/pkg/api/v1"
	"github.com/aide-family/rabbit/pkg/enum"
)

type SaveNamespaceBo struct {
	Name     string
	Metadata map[string]string
}

func (b *SaveNamespaceBo) ToDoNamespace() *do.Namespace {
	return &do.Namespace{
		Name:     b.Name,
		Metadata: b.Metadata,
	}
}

type UpdateNamespaceStatusBo struct {
	Name   string
	Status vobj.GlobalStatus
}

type ListNamespaceBo struct {
	*PageRequestBo
	Keyword string
	Status  vobj.GlobalStatus
}

type NamespaceItemBo struct {
	Name      string
	Metadata  map[string]string
	Status    vobj.GlobalStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewNamespaceItemBo(doNamespace *do.Namespace) *NamespaceItemBo {
	return &NamespaceItemBo{
		Name:      doNamespace.Name,
		Metadata:  doNamespace.Metadata,
		Status:    doNamespace.Status,
		CreatedAt: doNamespace.CreatedAt,
		UpdatedAt: doNamespace.UpdatedAt,
	}
}

func (b *NamespaceItemBo) ToAPIV1NamespaceItem() *apiv1.NamespaceItem {
	return &apiv1.NamespaceItem{
		Name:      b.Name,
		Metadata:  b.Metadata,
		Status:    enum.GlobalStatus(b.Status),
		CreatedAt: b.CreatedAt.Format(time.DateTime),
		UpdatedAt: b.UpdatedAt.Format(time.DateTime),
	}
}

func NewSaveNamespaceBo(name string, metadata map[string]string) *SaveNamespaceBo {
	return &SaveNamespaceBo{
		Name:     name,
		Metadata: metadata,
	}
}

func NewUpdateNamespaceStatusBo(name string, status vobj.GlobalStatus) *UpdateNamespaceStatusBo {
	return &UpdateNamespaceStatusBo{
		Name:   name,
		Status: status,
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
