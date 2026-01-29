package do

import (
	"github.com/aide-family/magicbox/safety"
	"github.com/aide-family/magicbox/strutil"
	"github.com/aide-family/rabbit/pkg/enum"
	"github.com/bwmarrin/snowflake"
)

type RecipientGroup struct {
	BaseModel

	Name         string                      `gorm:"column:name;uniqueIndex"`
	Metadata     *safety.Map[string, string] `gorm:"column:metadata;type:json;"`
	Status       enum.GlobalStatus           `gorm:"column:status;default:0"`
	Templates    []*Template                 `gorm:"many2many:recipient_group__templates;"`
	EmailConfigs []*EmailConfig              `gorm:"many2many:recipient_group__email_configs;"`
	Webhooks     []*WebhookConfig            `gorm:"many2many:recipient_group__webhooks;"`
	Members      []*RecipientMember          `gorm:"many2many:recipient_group__members;"`
}

func (RecipientGroup) TableName() string {
	return "recipient_groups"
}

type RecipientMember struct {
	BaseModel

	UserID snowflake.ID          `gorm:"column:user_id;uniqueIndex"`
	Email  strutil.EncryptString `gorm:"column:email;uniqueIndex"`
	Phone  strutil.EncryptString `gorm:"column:phone;uniqueIndex"`
	Status enum.GlobalStatus     `gorm:"column:status;default:0"`
}

func (RecipientMember) TableName() string {
	return "recipient_members"
}
