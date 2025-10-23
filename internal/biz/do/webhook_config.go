package do

import (
	"github.com/aide-family/magicbox/strutil"
	"github.com/google/uuid"

	"github.com/aide-family/rabbit/internal/biz/vobj"
)

type WebhookConfig struct {
	BaseModel

	UID     uuid.UUID             `gorm:"column:uid;type:varchar(36);not null;uniqueIndex"`
	App     vobj.WebhookApp       `gorm:"column:app;type:tinyint(2);not null;default:0"`
	Name    string                `gorm:"column:name;type:varchar(100);not null;uniqueIndex"`
	URL     string                `gorm:"column:url;type:varchar(255);not null;uniqueIndex"`
	Method  vobj.HTTPMethod       `gorm:"column:method;type:tinyint(2);not null;default:0"`
	Headers map[string]string     `gorm:"column:headers;type:json;"`
	Secret  strutil.EncryptString `gorm:"column:secret;type:varchar(512);not null"`
	Status  vobj.GlobalStatus     `gorm:"column:status;type:tinyint(2);not null;default:0"`
}

func (WebhookConfig) TableName() string {
	return "webhooks"
}
