package do

import (
	"encoding/json"
	"net/http"

	"github.com/aide-family/magicbox/strutil"
	"github.com/aide-family/rabbit/internal/biz/vobj"
	"github.com/aide-family/rabbit/pkg/merr"
	"gorm.io/gorm"
)

// Template 统一的模板结构
type Template struct {
	NamespaceModel

	Name     string            `gorm:"column:name;type:varchar(100);not null;uniqueIndex"`
	App      vobj.TemplateApp  `gorm:"column:app;type:tinyint(2);not null;default:0"`
	JSONData json.RawMessage   `gorm:"column:json_data;type:json;not null"`
	Status   vobj.GlobalStatus `gorm:"column:status;type:tinyint(2);not null;default:0"`
}

func (Template) TableName() string {
	return "templates"
}

func (t *Template) BeforeCreate(tx *gorm.DB) (err error) {
	if err = t.NamespaceModel.BeforeCreate(tx); err != nil {
		return
	}
	if !t.Status.Exist() || t.Status.IsUnknown() {
		t.Status = vobj.GlobalStatusEnabled
	}
	if strutil.IsEmpty(t.Name) {
		return merr.ErrorParams("template name is required")
	}
	if !t.App.Exist() || t.App.IsUnknown() {
		return merr.ErrorParams("invalid template app")
	}
	if len(t.JSONData) == 0 {
		return merr.ErrorParams("template json_data is required")
	}
	return
}

// EmailTemplateData Email 模板的数据结构
type EmailTemplateData struct {
	Subject     string      `json:"subject"`
	Body        string      `json:"body"`
	ContentType string      `json:"content_type"`
	Headers     http.Header `json:"headers,omitempty"`
}

// SMSTemplateData SMS 模板的数据结构
type SMSTemplateData struct {
	Content string            `json:"content"`
	Params  map[string]string `json:"params,omitempty"`
}

// WebhookTemplateData Webhook 模板的数据结构
type WebhookTemplateData string

// ToEmailTemplateData 将 JSONData 转换为 EmailTemplateData
func (t *Template) ToEmailTemplateData() (*EmailTemplateData, error) {
	var data EmailTemplateData
	if err := json.Unmarshal(t.JSONData, &data); err != nil {
		return nil, err
	}
	if data.ContentType == "" {
		data.ContentType = "text/html"
	}
	return &data, nil
}

// ToSMSTemplateData 将 JSONData 转换为 SMSTemplateData
func (t *Template) ToSMSTemplateData() (*SMSTemplateData, error) {
	var data SMSTemplateData
	if err := json.Unmarshal(t.JSONData, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// ToWebhookTemplateData 将 JSONData 转换为 WebhookTemplateData
func (t *Template) ToWebhookTemplateData() (WebhookTemplateData, error) {
	return WebhookTemplateData(t.JSONData), nil
}

// SetEmailTemplateData 设置 Email 模板数据
func (t *Template) SetEmailTemplateData(data *EmailTemplateData) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	t.JSONData = jsonData
	return nil
}

// SetSMSTemplateData 设置 SMS 模板数据
func (t *Template) SetSMSTemplateData(data *SMSTemplateData) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	t.JSONData = jsonData
	return nil
}

// SetWebhookTemplateData 设置 Webhook 模板数据
func (t *Template) SetWebhookTemplateData(data *WebhookTemplateData) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	t.JSONData = jsonData
	return nil
}
