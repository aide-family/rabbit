package bo

import (
	"time"

	"github.com/aide-family/rabbit/internal/biz/do"
	"github.com/aide-family/rabbit/internal/biz/vobj"
	apiv1 "github.com/aide-family/rabbit/pkg/api/v1"
	"github.com/aide-family/rabbit/pkg/enum"
)

type CreateWebhookTemplateBo struct {
	App  vobj.WebhookApp
	Name string
	Body string
}

func (b *CreateWebhookTemplateBo) ToDoWebhookTemplate() *do.WebhookTemplate {
	return &do.WebhookTemplate{
		App:  b.App,
		Name: b.Name,
		Body: b.Body,
	}
}

func NewCreateWebhookTemplateBo(req *apiv1.CreateWebhookTemplateRequest) *CreateWebhookTemplateBo {
	return &CreateWebhookTemplateBo{
		App:  vobj.WebhookApp(req.App),
		Name: req.Name,
		Body: req.Body,
	}
}

type UpdateWebhookTemplateBo struct {
	UID  string
	Name string
	Body string
	App  vobj.WebhookApp
}

func (b *UpdateWebhookTemplateBo) ToDoWebhookTemplate() *do.WebhookTemplate {
	return &do.WebhookTemplate{
		NamespaceModel: do.NamespaceModel{
			UID: b.UID,
		},
		Name: b.Name,
		Body: b.Body,
		App:  b.App,
	}
}

func NewUpdateWebhookTemplateBo(req *apiv1.UpdateWebhookTemplateRequest) *UpdateWebhookTemplateBo {
	return &UpdateWebhookTemplateBo{
		UID:  req.Uid,
		Name: req.Name,
		Body: req.Body,
		App:  vobj.WebhookApp(req.App),
	}
}

type UpdateWebhookTemplateStatusBo struct {
	UID    string
	Status vobj.GlobalStatus
}

func NewUpdateWebhookTemplateStatusBo(req *apiv1.UpdateWebhookTemplateStatusRequest) *UpdateWebhookTemplateStatusBo {
	return &UpdateWebhookTemplateStatusBo{
		UID:    req.Uid,
		Status: vobj.GlobalStatus(req.Status),
	}
}

type WebhookTemplateItemBo struct {
	UID       string
	App       vobj.WebhookApp
	Name      string
	Body      string
	Status    vobj.GlobalStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewWebhookTemplateItemBo(doTemplate *do.WebhookTemplate) *WebhookTemplateItemBo {
	return &WebhookTemplateItemBo{
		UID:       doTemplate.UID,
		App:       doTemplate.App,
		Name:      doTemplate.Name,
		Body:      doTemplate.Body,
		Status:    doTemplate.Status,
		CreatedAt: doTemplate.CreatedAt,
		UpdatedAt: doTemplate.UpdatedAt,
	}
}

func (b *WebhookTemplateItemBo) ToAPIV1WebhookTemplateItem() *apiv1.WebhookTemplateItem {
	return &apiv1.WebhookTemplateItem{
		Uid:       b.UID,
		App:       enum.WebhookAPP(b.App),
		Name:      b.Name,
		Body:      b.Body,
		Status:    enum.GlobalStatus(b.Status),
		CreatedAt: b.CreatedAt.Format(time.DateTime),
		UpdatedAt: b.UpdatedAt.Format(time.DateTime),
	}
}

type ListWebhookTemplateBo struct {
	*PageRequestBo
	Keyword string
	Status  vobj.GlobalStatus
	App     vobj.WebhookApp
}

func NewListWebhookTemplateBo(req *apiv1.ListWebhookTemplateRequest) *ListWebhookTemplateBo {
	return &ListWebhookTemplateBo{
		PageRequestBo: NewPageRequestBo(req.Page, req.PageSize),
		Keyword:       req.Keyword,
		Status:        vobj.GlobalStatus(req.Status),
		App:           vobj.WebhookApp(req.App),
	}
}

func ToAPIV1ListWebhookTemplateReply(pageResponseBo *PageResponseBo[*WebhookTemplateItemBo]) *apiv1.ListWebhookTemplateReply {
	items := make([]*apiv1.WebhookTemplateItem, 0, len(pageResponseBo.GetItems()))
	for _, item := range pageResponseBo.GetItems() {
		items = append(items, item.ToAPIV1WebhookTemplateItem())
	}
	return &apiv1.ListWebhookTemplateReply{
		Items:    items,
		Total:    int32(pageResponseBo.GetTotal()),
		Page:     pageResponseBo.GetPage(),
		PageSize: pageResponseBo.GetPageSize(),
	}
}
