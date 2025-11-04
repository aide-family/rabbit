package bo

import (
	"time"

	"github.com/bwmarrin/snowflake"

	"github.com/aide-family/rabbit/internal/biz/do"
	"github.com/aide-family/rabbit/internal/biz/vobj"
	apiv1 "github.com/aide-family/rabbit/pkg/api/v1"
	"github.com/aide-family/rabbit/pkg/enum"
)

// CreateTemplateBo 创建模板的 BO
type CreateTemplateBo struct {
	Name     string
	App      vobj.TemplateApp
	JSONData []byte
}

// ToDoTemplate 转换为 DO
func (c *CreateTemplateBo) ToDoTemplate() *do.Template {
	return &do.Template{
		Name:     c.Name,
		App:      c.App,
		JSONData: c.JSONData,
	}
}

// NewCreateTemplateBo 从 API 请求创建 BO
func NewCreateTemplateBo(req *apiv1.CreateTemplateRequest) *CreateTemplateBo {
	return &CreateTemplateBo{
		Name:     req.Name,
		App:      vobj.TemplateApp(req.App),
		JSONData: req.JsonData,
	}
}

// UpdateTemplateBo 更新模板的 BO
type UpdateTemplateBo struct {
	UID      snowflake.ID
	Name     string
	App      vobj.TemplateApp
	JSONData []byte
}

// ToDoTemplate 转换为 DO
func (u *UpdateTemplateBo) ToDoTemplate() *do.Template {
	template := &do.Template{
		Name:     u.Name,
		App:      u.App,
		JSONData: u.JSONData,
	}
	template.WithUID(u.UID)
	return template
}

// NewUpdateTemplateBo 从 API 请求创建 BO
func NewUpdateTemplateBo(req *apiv1.UpdateTemplateRequest) *UpdateTemplateBo {
	return &UpdateTemplateBo{
		UID:      snowflake.ParseInt64(req.Uid),
		Name:     req.Name,
		App:      vobj.TemplateApp(req.App),
		JSONData: req.JsonData,
	}
}

// UpdateTemplateStatusBo 更新模板状态的 BO
type UpdateTemplateStatusBo struct {
	UID    snowflake.ID
	Status vobj.GlobalStatus
}

// NewUpdateTemplateStatusBo 从 API 请求创建 BO
func NewUpdateTemplateStatusBo(req *apiv1.UpdateTemplateStatusRequest) *UpdateTemplateStatusBo {
	return &UpdateTemplateStatusBo{
		UID:    snowflake.ParseInt64(req.Uid),
		Status: vobj.GlobalStatus(req.Status),
	}
}

// TemplateItemBo 模板项的 BO
type TemplateItemBo struct {
	UID       snowflake.ID
	Name      string
	App       vobj.TemplateApp
	JSONData  []byte
	Status    vobj.GlobalStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewTemplateItemBo 从 DO 创建 BO
func NewTemplateItemBo(doTemplate *do.Template) *TemplateItemBo {
	return &TemplateItemBo{
		UID:       doTemplate.UID,
		Name:      doTemplate.Name,
		App:       doTemplate.App,
		JSONData:  doTemplate.JSONData,
		Status:    doTemplate.Status,
		CreatedAt: doTemplate.CreatedAt,
		UpdatedAt: doTemplate.UpdatedAt,
	}
}

// ToAPIV1TemplateItem 转换为 API 响应
func (b *TemplateItemBo) ToAPIV1TemplateItem() *apiv1.TemplateItem {
	return &apiv1.TemplateItem{
		Uid:       b.UID.Int64(),
		Name:      b.Name,
		App:       enum.TemplateAPP(b.App),
		JsonData:  b.JSONData,
		Status:    enum.GlobalStatus(b.Status),
		CreatedAt: b.CreatedAt.Format(time.DateTime),
		UpdatedAt: b.UpdatedAt.Format(time.DateTime),
	}
}

// ListTemplateBo 列表查询的 BO
type ListTemplateBo struct {
	*PageRequestBo
	Keyword string
	Status  vobj.GlobalStatus
	App     vobj.TemplateApp
}

// NewListTemplateBo 从 API 请求创建 BO
func NewListTemplateBo(req *apiv1.ListTemplateRequest) *ListTemplateBo {
	return &ListTemplateBo{
		PageRequestBo: NewPageRequestBo(req.Page, req.PageSize),
		Keyword:       req.Keyword,
		Status:        vobj.GlobalStatus(req.Status),
		App:           vobj.TemplateApp(req.App),
	}
}

// ToAPIV1ListTemplateReply 转换为 API 响应
func ToAPIV1ListTemplateReply(pageResponseBo *PageResponseBo[*TemplateItemBo]) *apiv1.ListTemplateReply {
	items := make([]*apiv1.TemplateItem, 0, len(pageResponseBo.GetItems()))
	for _, item := range pageResponseBo.GetItems() {
		items = append(items, item.ToAPIV1TemplateItem())
	}
	return &apiv1.ListTemplateReply{
		Items:    items,
		Total:    pageResponseBo.GetTotal(),
		Page:     pageResponseBo.GetPage(),
		PageSize: pageResponseBo.GetPageSize(),
	}
}
