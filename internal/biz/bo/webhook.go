package bo

import (
	"encoding/json"
	"time"

	"github.com/aide-family/magicbox/message"
	"github.com/aide-family/magicbox/safety"
	"github.com/aide-family/magicbox/serialize"
	"github.com/aide-family/magicbox/strutil"
	"github.com/bwmarrin/snowflake"

	"github.com/aide-family/rabbit/internal/biz/do"
	"github.com/aide-family/rabbit/internal/biz/vobj"
	apiv1 "github.com/aide-family/rabbit/pkg/api/v1"
	"github.com/aide-family/rabbit/pkg/enum"
	"github.com/aide-family/rabbit/pkg/merr"
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

type SendWebhookBo struct {
	UID  snowflake.ID `json:"uid"`
	Data string       `json:"data"`
}

// Message implements message.Message.
func (b *SendWebhookBo) Message(message.MessageChannel) ([]byte, error) {
	if !json.Valid([]byte(b.Data)) {
		return nil, merr.ErrorParams("invalid json data")
	}
	return []byte(b.Data), nil
}

func (b *SendWebhookBo) ToMessageLog(webhookConfig *do.WebhookConfig) (*do.MessageLog, error) {
	messageBytes, err := serialize.JSONMarshal(b)
	if err != nil {
		return nil, err
	}
	webhookConfigBytes, err := serialize.JSONMarshal(NewWebhookConfigItemBo(webhookConfig))
	if err != nil {
		return nil, err
	}
	return &do.MessageLog{
		SendAt:  time.Now(),
		Message: strutil.EncryptString(messageBytes),
		Config:  strutil.EncryptString(webhookConfigBytes),
		Type:    vobj.MessageTypeWebhook,
		Status:  vobj.MessageStatusPending,
	}, nil
}

func NewSendWebhookBo(req *apiv1.SendWebhookRequest) *SendWebhookBo {
	return &SendWebhookBo{
		UID:  snowflake.ParseInt64(req.Uid),
		Data: req.Data,
	}
}

type SendWebhookWithTemplateBo struct {
	UID         snowflake.ID
	TemplateUID snowflake.ID
	JSONData    []byte
}

func NewSendWebhookWithTemplateBo(req *apiv1.SendWebhookWithTemplateRequest) (*SendWebhookWithTemplateBo, error) {
	if !json.Valid([]byte(req.JsonData)) {
		return nil, merr.ErrorParams("invalid json data")
	}
	return &SendWebhookWithTemplateBo{
		UID:         snowflake.ParseInt64(req.Uid),
		TemplateUID: snowflake.ParseInt64(req.TemplateUID),
		JSONData:    []byte(req.JsonData),
	}, nil
}

func (b *SendWebhookWithTemplateBo) ToSendWebhookBo(templateDo *do.Template) (*SendWebhookBo, error) {
	if !templateDo.App.IsWebhookType() {
		return nil, merr.ErrorParams("invalid template app type, expected webhook type, got %s", templateDo.App)
	}
	if !templateDo.Status.IsEnabled() {
		return nil, merr.ErrorParams("template %s(%s) is disabled", templateDo.Name, templateDo.UID)
	}
	webhookTemplateData, err := templateDo.ToWebhookTemplateData()
	if err != nil {
		return nil, err
	}
	var jsonData map[string]any
	if err := serialize.JSONUnmarshal(b.JSONData, &jsonData); err != nil {
		return nil, merr.ErrorInternal("unmarshal json data failed").WithCause(err)
	}

	bodyData, err := strutil.ExecuteTextTemplate(webhookTemplateData.Body, jsonData)
	if err != nil {
		return nil, merr.ErrorParams("execute text template failed").WithCause(err)
	}

	return &SendWebhookBo{
		UID:  b.UID,
		Data: bodyData,
	}, nil
}

type WebhookConfigItemBo struct {
	UID       snowflake.ID      `json:"uid"`
	App       vobj.WebhookApp   `json:"app"`
	Name      string            `json:"name"`
	URL       string            `json:"url"`
	Method    vobj.HTTPMethod   `json:"method"`
	Headers   map[string]string `json:"headers"`
	Secret    string            `json:"secret"`
	Status    vobj.GlobalStatus `json:"status"`
	CreatedAt time.Time         `json:"-"`
	UpdatedAt time.Time         `json:"-"`
}

// GetSecret implements dingtalk.Config.
func (w *WebhookConfigItemBo) GetSecret() string {
	return w.Secret
}

// GetURL implements dingtalk.Config.
func (w *WebhookConfigItemBo) GetURL() string {
	return w.URL
}

func NewWebhookConfigItemBo(doWebhookConfig *do.WebhookConfig) *WebhookConfigItemBo {
	return &WebhookConfigItemBo{
		UID:       doWebhookConfig.UID,
		App:       doWebhookConfig.App,
		Name:      doWebhookConfig.Name,
		URL:       doWebhookConfig.URL,
		Method:    doWebhookConfig.Method,
		Headers:   doWebhookConfig.Headers.Map(),
		Secret:    string(doWebhookConfig.Secret),
		Status:    doWebhookConfig.Status,
		CreatedAt: doWebhookConfig.CreatedAt,
		UpdatedAt: doWebhookConfig.UpdatedAt,
	}
}
