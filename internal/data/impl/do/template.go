package do

import (
	"encoding/json"

	"github.com/aide-family/rabbit/pkg/enum"
	"github.com/bwmarrin/snowflake"
)

// Template 统一的模板结构
type Template struct {
	BaseModel

	NamespaceUID snowflake.ID      `gorm:"column:namespace_uid;uniqueIndex:template__namespace_uid__name"`
	Name         string            `gorm:"column:name;uniqueIndex:template__namespace_uid__name"`
	MessageType  enum.MessageType  `gorm:"column:message_type;default:0"`
	JSONData     json.RawMessage   `gorm:"column:json_data;"`
	Status       enum.GlobalStatus `gorm:"column:status;default:0"`
}

func (Template) TableName() string {
	return "templates"
}
