package bo

import (
	"time"

	"github.com/aide-family/rabbit/internal/biz/do"
	"github.com/aide-family/rabbit/internal/biz/vobj"
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

// NewCreateWebhookTemplateBo 已被移除，请使用 template 统一接口
// 使用 bo.CreateTemplateBoFromWebhook 将旧的 BO 转换为新的统一 Template BO

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

// NewUpdateWebhookTemplateBo 已被移除，请使用 template 统一接口

type UpdateWebhookTemplateStatusBo struct {
	UID    string
	Status vobj.GlobalStatus
}

// NewUpdateWebhookTemplateStatusBo 已被移除，请使用 template 统一接口

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

// ToAPIV1WebhookTemplateItem 已被移除，请使用 template 统一接口

type ListWebhookTemplateBo struct {
	*PageRequestBo
	Keyword string
	Status  vobj.GlobalStatus
	App     vobj.WebhookApp
}

// NewListWebhookTemplateBo 已被移除，请使用 template 统一接口

// ToAPIV1ListWebhookTemplateReply 已被移除，请使用 template 统一接口
