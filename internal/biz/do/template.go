package do

import (
	"encoding/json"

	"github.com/aide-family/magicbox/strutil"
	"github.com/aide-family/rabbit/internal/biz/vobj"
	"github.com/aide-family/rabbit/pkg/merr"
	"gorm.io/gorm"
)

// Template 统一的模板结构
type Template struct {
	NamespaceModel

	Name     string            `gorm:"column:name;type:varchar(100);not null;uniqueIndex"`
	App      vobj.TemplateApp  `gorm:"column:app;type:tinyint(2);not null;default:0"`
	JSONData json.RawMessage   `gorm:"column:json_data;type:json;not null"`
	Status   vobj.GlobalStatus `gorm:"column:status;type:tinyint(2);not null;default:0"`
}

func (Template) TableName() string {
	return "templates"
}

func (t *Template) BeforeCreate(tx *gorm.DB) (err error) {
	if err = t.NamespaceModel.BeforeCreate(tx); err != nil {
		return
	}
	if !t.Status.Exist() || t.Status.IsUnknown() {
		t.Status = vobj.GlobalStatusEnabled
	}
	if strutil.IsEmpty(t.Name) {
		return merr.ErrorParams("template name is required")
	}
	if !t.App.Exist() || t.App.IsUnknown() {
		return merr.ErrorParams("invalid template app")
	}
	if len(t.JSONData) == 0 {
		return merr.ErrorParams("template json_data is required")
	} else if !json.Valid(t.JSONData) {
		return merr.ErrorParams("invalid template json_data")
	}
	return
}
