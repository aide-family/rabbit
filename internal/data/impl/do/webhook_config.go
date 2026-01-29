package do

import (
	"github.com/aide-family/magicbox/safety"
	"github.com/aide-family/magicbox/strutil"

	"github.com/aide-family/rabbit/pkg/enum"
)

type WebhookConfig struct {
	BaseModel

	App     enum.WebhookAPP             `gorm:"column:app;default:0"`
	Name    string                      `gorm:"column:name;uniqueIndex"`
	URL     string                      `gorm:"column:url;uniqueIndex"`
	Method  enum.HTTPMethod             `gorm:"column:method;default:0"`
	Headers *safety.Map[string, string] `gorm:"column:headers;type:json;"`
	Secret  strutil.EncryptString       `gorm:"column:secret;"`
	Status  enum.GlobalStatus           `gorm:"column:status;default:0"`
}

func (WebhookConfig) TableName() string {
	return "webhooks"
}
