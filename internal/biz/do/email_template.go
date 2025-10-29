package do

import (
	"github.com/aide-family/magicbox/safety"
	"github.com/aide-family/rabbit/internal/biz/vobj"
)

type EmailTemplate struct {
	NamespaceModel

	Name        string                      `gorm:"column:name;type:varchar(100);not null;uniqueIndex"`
	Subject     string                      `gorm:"column:subject;type:varchar(255);not null"`
	Body        string                      `gorm:"column:body;type:text;not null"`
	ContentType string                      `gorm:"column:content_type;type:varchar(255);not null"`
	Headers     *safety.Map[string, string] `gorm:"column:headers;type:json;"`
	Status      vobj.GlobalStatus           `gorm:"column:status;type:tinyint(2);not null;default:0"`
}

func (EmailTemplate) TableName() string {
	return "email_templates"
}
