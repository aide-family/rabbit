package do

import (
	"github.com/aide-family/magicbox/strutil"
	"github.com/google/uuid"

	"github.com/aide-family/rabbit/internal/biz/vobj"
)

type EmailConfig struct {
	BaseModel

	UID      uuid.UUID             `gorm:"column:uid;type:varchar(36);not null;uniqueIndex"`
	Name     string                `gorm:"column:name;type:varchar(100);not null;uniqueIndex"`
	Host     string                `gorm:"column:host;type:varchar(255);not null"`
	Port     int32                 `gorm:"column:port;type:int(11);not null"`
	Username string                `gorm:"column:username;type:varchar(255);not null"`
	Password strutil.EncryptString `gorm:"column:password;type:varchar(512);not null"`
	Status   vobj.GlobalStatus     `gorm:"column:status;type:tinyint(2);not null;default:0"`
}

func (EmailConfig) TableName() string {
	return "email_configs"
}
