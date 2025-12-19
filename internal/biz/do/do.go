// Package do is the data object package for the Rabbit service.
package do

import (
	"context"
	"time"

	"github.com/aide-family/magicbox/hello"
	"github.com/aide-family/magicbox/strutil"
	"github.com/bwmarrin/snowflake"
	"gorm.io/gorm"

	"github.com/aide-family/rabbit/pkg/middler"
)

func Models() []any {
	return []any{
		&Namespace{},
		&WebhookConfig{},
		&EmailConfig{},
		&Template{},
		&MessageLog{},
		&MessageRetryLog{},
	}
}

type BaseModel struct {
	ID        uint32         `gorm:"column:id;type:int unsigned;primaryKey;autoIncrement"`
	UID       snowflake.ID   `gorm:"column:uid;type:bigint(20) unsigned;not null;uniqueIndex"`
	CreatedAt time.Time      `gorm:"column:created_at;type:datetime;not null;"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:datetime;not null;"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetime;index"`
	Creator   snowflake.ID   `gorm:"column:creator;type:bigint(20) unsigned;not null;index"`
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	if b.Creator == 0 {
		b.WithCreator(tx.Statement.Context)
	}

	node, err := snowflake.NewNode(hello.NodeID())
	if err != nil {
		return err
	}
	b.WithUID(node.Generate())

	return
}

func (b *BaseModel) WithCreator(ctx context.Context) *BaseModel {
	b.Creator = middler.GetBaseInfo(ctx).UserID
	return b
}

func (b *BaseModel) WithUID(uid snowflake.ID) *BaseModel {
	b.UID = uid
	return b
}

type NamespaceModel struct {
	BaseModel

	Namespace string `gorm:"column:namespace;type:varchar(100);not null;index"`
}

func (n *NamespaceModel) BeforeCreate(tx *gorm.DB) (err error) {
	if err = n.BaseModel.BeforeCreate(tx); err != nil {
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
