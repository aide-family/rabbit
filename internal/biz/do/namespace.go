package do

import "github.com/aide-family/rabbit/internal/biz/vobj"

type Namespace struct {
	BaseModel

	Name     string               `gorm:"column:name;type:varchar(100);not null;unique"`
	Metadata map[string]string    `gorm:"column:metadata;type:json;"`
	Status   vobj.NamespaceStatus `gorm:"column:status;type:tinyint(2);not null;default:0"`
}

func (Namespace) TableName() string {
	return "namespace"
}
