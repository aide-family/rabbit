// Package do is the data object package for the Rabbit service.
package do

import (
	"context"
	"time"

	"github.com/aide-family/magicbox/strutil"
	"github.com/aide-family/rabbit/pkg/middler"
	"gorm.io/gorm"
)

func Models() []any {
	return []any{
		&Namespace{},
		&WebhookConfig{},
		&EmailConfig{},
		&EmailTemplate{},
		&WebhookTemplate{},
		&MessageLog{},
	}
}

type BaseModel struct {
	ID        uint32         `gorm:"column:id;type:int unsigned;primaryKey;autoIncrement"`
	CreatedAt time.Time      `gorm:"column:created_at;type:datetime;not null;"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:datetime;not null;"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetime;index"`
	Creator   string         `gorm:"column:creator;type:varchar(36);not null;index"`
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	if strutil.IsEmpty(b.Creator) {
		b.WithCreator(tx.Statement.Context)
	}
	return
}

func (b *BaseModel) WithCreator(ctx context.Context) *BaseModel {
	b.Creator = middler.GetBaseInfo(ctx).UserID
	return b
}

type NamespaceModel struct {
	BaseModel

	Namespace string `gorm:"column:namespace;type:varchar(100);not null;index"`
}

func (n *NamespaceModel) BeforeCreate(tx *gorm.DB) (err error) {
	if n.BaseModel.BeforeCreate(tx) != nil {
		return
	}
	if strutil.IsEmpty(n.Namespace) {
		n.WithNamespace(middler.GetNamespace(tx.Statement.Context))
	}
	return
}

func (n *NamespaceModel) WithNamespace(namespace string) *NamespaceModel {
	n.Namespace = namespace
	return n
}

func HasTable(tx *gorm.DB, tableName string) bool {
	return tx.Migrator().HasTable(tableName)
}

func getFirstMonday(date time.Time) time.Time {
	offset := int(time.Monday - date.Weekday())
	if offset > 0 {
		offset -= 7
	}
	return date.AddDate(0, 0, offset)
}
