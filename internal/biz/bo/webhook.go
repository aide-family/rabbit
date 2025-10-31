package bo

import (
	"time"

	"github.com/aide-family/magicbox/safety"
	"github.com/aide-family/magicbox/strutil"
	"github.com/bwmarrin/snowflake"

	"github.com/aide-family/rabbit/internal/biz/do"
	"github.com/aide-family/rabbit/internal/biz/vobj"
	apiv1 "github.com/aide-family/rabbit/pkg/api/v1"
	"github.com/aide-family/rabbit/pkg/enum"
)

type CreateWebhookBo struct {
	App     vobj.WebhookApp
	Name    string
	URL     string
	Method  vobj.HTTPMethod
	Headers map[string]string
	Secret  string
}

func (b *CreateWebhookBo) ToDoWebhookConfig() *do.WebhookConfig {
	return &do.WebhookConfig{
		App:     b.App,
		Name:    b.Name,
		URL:     b.URL,
		Method:  b.Method,
		Headers: safety.NewMap(b.Headers),
		Secret:  strutil.EncryptString(b.Secret),
	}
}

func NewCreateWebhookBo(req *apiv1.CreateWebhookRequest) *CreateWebhookBo {
	return &CreateWebhookBo{
		App:     vobj.WebhookApp(req.App),
		Name:    req.Name,
		URL:     req.Url,
		Method:  vobj.HTTPMethod(req.Method),
		Headers: req.Headers,
		Secret:  req.Secret,
	}
}

type UpdateWebhookBo struct {
	UID     snowflake.ID
	App     vobj.WebhookApp
	Name    string
	URL     string
	Method  vobj.HTTPMethod
	Headers map[string]string
	Secret  string
}

func (b *UpdateWebhookBo) ToDoWebhookConfig() *do.WebhookConfig {
	webhookConfig := &do.WebhookConfig{
		App:     b.App,
		Name:    b.Name,
		URL:     b.URL,
		Method:  b.Method,
		Headers: safety.NewMap(b.Headers),
		Secret:  strutil.EncryptString(b.Secret),
	}
	webhookConfig.WithUID(b.UID)
	return webhookConfig
}

func NewUpdateWebhookBo(req *apiv1.UpdateWebhookRequest) *UpdateWebhookBo {
	return &UpdateWebhookBo{
		UID:     snowflake.ParseInt64(req.Uid),
		App:     vobj.WebhookApp(req.App),
		Name:    req.Name,
		URL:     req.Url,
		Method:  vobj.HTTPMethod(req.Method),
		Headers: req.Headers,
		Secret:  req.Secret,
	}
}

type UpdateWebhookStatusBo struct {
	UID    snowflake.ID
	Status vobj.GlobalStatus
}

func NewUpdateWebhookStatusBo(req *apiv1.UpdateWebhookStatusRequest) *UpdateWebhookStatusBo {
	return &UpdateWebhookStatusBo{
		UID:    snowflake.ParseInt64(req.Uid),
		Status: vobj.GlobalStatus(req.Status),
	}
}

type WebhookItemBo struct {
	UID       snowflake.ID
	App       vobj.WebhookApp
	Name      string
	URL       string
	Method    vobj.HTTPMethod
	Headers   map[string]string
	Secret    string
	Status    vobj.GlobalStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewWebhookItemBo(doWebhook *do.WebhookConfig) *WebhookItemBo {
	return &WebhookItemBo{
		UID:       doWebhook.UID,
		App:       doWebhook.App,
		Name:      doWebhook.Name,
		URL:       doWebhook.URL,
		Method:    doWebhook.Method,
		Headers:   doWebhook.Headers.Map(),
		Secret:    string(doWebhook.Secret),
		Status:    doWebhook.Status,
		CreatedAt: doWebhook.CreatedAt,
		UpdatedAt: doWebhook.UpdatedAt,
	}
}

func (b *WebhookItemBo) ToAPIV1WebhookItem() *apiv1.WebhookItem {
	return &apiv1.WebhookItem{
		Uid:       b.UID.Int64(),
		App:       enum.WebhookAPP(b.App),
		Name:      b.Name,
		Url:       b.URL,
		Method:    enum.HTTPMethod(b.Method),
		Headers:   b.Headers,
		Secret:    b.Secret,
		Status:    enum.GlobalStatus(b.Status),
		CreatedAt: b.CreatedAt.Format(time.DateTime),
		UpdatedAt: b.UpdatedAt.Format(time.DateTime),
	}
}

type ListWebhookBo struct {
	*PageRequestBo
	App     vobj.WebhookApp
	Keyword string
}

func NewListWebhookBo(req *apiv1.ListWebhookRequest) *ListWebhookBo {
	return &ListWebhookBo{
		PageRequestBo: NewPageRequestBo(req.Page, req.PageSize),
		App:           vobj.WebhookApp(req.App),
		Keyword:       req.Keyword,
	}
}

func ToAPIV1ListWebhookReply(pageResponseBo *PageResponseBo[*WebhookItemBo]) *apiv1.ListWebhookReply {
	items := make([]*apiv1.WebhookItem, 0, len(pageResponseBo.GetItems()))
	for _, item := range pageResponseBo.GetItems() {
		items = append(items, item.ToAPIV1WebhookItem())
	}
	return &apiv1.ListWebhookReply{
		Items:    items,
		Total:    pageResponseBo.GetTotal(),
		Page:     pageResponseBo.GetPage(),
		PageSize: pageResponseBo.GetPageSize(),
	}
}
