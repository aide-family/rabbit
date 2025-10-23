package do

import (
	"github.com/google/uuid"

	"github.com/aide-family/rabbit/internal/biz/vobj"
)

type WebhookTemplate struct {
	BaseModel

	UID    uuid.UUID         `gorm:"column:uid;type:varchar(36);not null;uniqueIndex"`
	Name   string            `gorm:"column:name;type:varchar(100);not null;uniqueIndex"`
	Body   string            `gorm:"column:body;type:text;not null"`
	Status vobj.GlobalStatus `gorm:"column:status;type:tinyint(2);not null;default:0"`
}

func (WebhookTemplate) TableName() string {
	return "webhook_templates"
}
