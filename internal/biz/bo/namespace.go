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

// SelectNamespaceBo 选择Namespace的 BO
type SelectNamespaceBo struct {
	Keyword string
	Limit   int32
	LastUID snowflake.ID
	Status  vobj.GlobalStatus
}

// NewSelectNamespaceBo 从 API 请求创建 BO
func NewSelectNamespaceBo(req *apiv1.SelectNamespaceRequest) *SelectNamespaceBo {
	var lastUID snowflake.ID
	if req.LastUID > 0 {
		lastUID = snowflake.ParseInt64(req.LastUID)
	}
	return &SelectNamespaceBo{
		Keyword: req.Keyword,
		Limit:   req.Limit,
		LastUID: lastUID,
		Status:  vobj.GlobalStatus(req.Status),
	}
}

// NamespaceItemSelectBo Namespace选择项的 BO
type NamespaceItemSelectBo struct {
	UID      snowflake.ID
	Name     string
	Status   vobj.GlobalStatus
	Disabled bool
	Tooltip  string
}

// NewNamespaceItemSelectBo 从 DO 创建 BO
func NewNamespaceItemSelectBo(doNamespace *do.Namespace) *NamespaceItemSelectBo {
	return &NamespaceItemSelectBo{
		UID:      doNamespace.UID,
		Name:     doNamespace.Name,
		Status:   doNamespace.Status,
		Disabled: doNamespace.Status != vobj.GlobalStatusEnabled,
		Tooltip:  "",
	}
}

// ToAPIV1NamespaceItemSelect 转换为 API 响应
func (b *NamespaceItemSelectBo) ToAPIV1NamespaceItemSelect() *apiv1.NamespaceItemSelect {
	return &apiv1.NamespaceItemSelect{
		Value:    b.UID.Int64(),
		Label:    b.Name,
		Disabled: b.Disabled,
		Tooltip:  b.Tooltip,
	}
}

// SelectNamespaceResult Repository层返回结果
type SelectNamespaceResult struct {
	Items   []*do.Namespace
	Total   int64
	LastUID snowflake.ID
}

// SelectNamespaceBoResult Biz层返回结果
type SelectNamespaceBoResult struct {
	Items   []*NamespaceItemSelectBo
	Total   int64
	LastUID snowflake.ID
}

// SelectNamespaceReplyParams 转换为API响应的参数
type SelectNamespaceReplyParams struct {
	Items   []*NamespaceItemSelectBo
	Total   int64
	LastUID snowflake.ID
	Limit   int32
}

// ToAPIV1SelectNamespaceReply 转换为 API 响应
func ToAPIV1SelectNamespaceReply(params *SelectNamespaceReplyParams) *apiv1.SelectNamespaceReply {
	selectItems := make([]*apiv1.NamespaceItemSelect, 0, len(params.Items))
	for _, item := range params.Items {
		selectItems = append(selectItems, item.ToAPIV1NamespaceItemSelect())
	}
	var lastUIDInt64 int64
	if params.LastUID > 0 {
		lastUIDInt64 = params.LastUID.Int64()
	}
	// hasMore: 如果返回的记录数等于limit，说明可能还有更多记录
	// 如果返回的记录数小于limit，说明已经查询完了
	hasMore := int32(len(params.Items)) == params.Limit
	return &apiv1.SelectNamespaceReply{
		Items:   selectItems,
		Total:   params.Total,
		LastUID: lastUIDInt64,
		HasMore: hasMore,
	}
}
