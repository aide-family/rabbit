package do

import (
	"github.com/aide-family/magicbox/strutil"
	"github.com/aide-family/rabbit/pkg/enum"
)

type EmailConfig struct {
	BaseModel

	Name     string                `gorm:"column:name;uniqueIndex"`
	Host     string                `gorm:"column:host;"`
	Port     int32                 `gorm:"column:port;"`
	Username string                `gorm:"column:username;"`
	Password strutil.EncryptString `gorm:"column:password;"`
	Status   enum.GlobalStatus     `gorm:"column:status;default:0"`
}

func (EmailConfig) TableName() string {
	return "email_configs"
}
