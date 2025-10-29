package do

import (
	"github.com/aide-family/magicbox/strutil"
	"github.com/aide-family/rabbit/internal/biz/vobj"
	"github.com/aide-family/rabbit/pkg/merr"
	"gorm.io/gorm"
)

type WebhookTemplate struct {
	NamespaceModel

	Name   string            `gorm:"column:name;type:varchar(100);not null;uniqueIndex"`
	App    vobj.WebhookApp   `gorm:"column:app;type:tinyint(2);not null;default:0"`
	Body   string            `gorm:"column:body;type:text;not null"`
	Status vobj.GlobalStatus `gorm:"column:status;type:tinyint(2);not null;default:0"`
}

func (WebhookTemplate) TableName() string {
	return "webhook_templates"
}

func (w *WebhookTemplate) BeforeCreate(tx *gorm.DB) (err error) {
	if err = w.NamespaceModel.BeforeCreate(tx); err != nil {
		return
	}
	if !w.Status.Exist() || w.Status.IsUnknown() {
		w.Status = vobj.GlobalStatusEnabled
	}
	if strutil.IsEmpty(w.Name) {
		return merr.ErrorParams("webhook template name is required")
	}
	if !w.App.Exist() || w.App.IsUnknown() {
		return merr.ErrorParams("invalid webhook template app")
	}
	if strutil.IsEmpty(w.Body) {
		return merr.ErrorParams("webhook template body is required")
	}
	return
}
