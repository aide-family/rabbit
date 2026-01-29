package do

import (
	"encoding/json"

	"github.com/aide-family/rabbit/pkg/enum"
)

// Template 统一的模板结构
type Template struct {
	BaseModel

	Name        string            `gorm:"column:name;uniqueIndex"`
	MessageType enum.MessageType  `gorm:"column:message_type;default:0"`
	JSONData    json.RawMessage   `gorm:"column:json_data;"`
	Status      enum.GlobalStatus `gorm:"column:status;default:0"`
}

func (Template) TableName() string {
	return "templates"
}
