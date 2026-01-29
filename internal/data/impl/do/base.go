// Package do is the data package for the internal data.
package do

import (
	"errors"
	"time"

	"github.com/aide-family/magicbox/hello"
	"github.com/bwmarrin/snowflake"
	"gorm.io/gorm"
)

func Models() []any {
	return []any{
		&EmailConfig{},
		&MessageLog{},
		&MessageRetryLog{},
		&RecipientGroup{},
		&RecipientMember{},
		&Template{},
		&WebhookConfig{},
	}
}

type BaseModel struct {
	ID           uint32         `gorm:"column:id;primaryKey;autoIncrement"`
	UID          snowflake.ID   `gorm:"column:uid;uniqueIndex"`
	CreatedAt    time.Time      `gorm:"column:created_at;"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index"`
	Creator      snowflake.ID   `gorm:"column:creator;index"`
	NamespaceUID snowflake.ID   `gorm:"column:namespace_uid;index"`
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	if b.Creator == 0 {
		return errors.New("creator is required")
	}
	if b.NamespaceUID == 0 {
		return errors.New("namespaceUID is required")
	}
	node, err := snowflake.NewNode(hello.NodeID())
	if err != nil {
		return err
	}
	b.UID = node.Generate()
	return nil
}

func (b *BaseModel) WithCreator(creator snowflake.ID) *BaseModel {
	b.Creator = creator
	return b
}

func (b *BaseModel) WithNamespaceUID(namespaceUID snowflake.ID) *BaseModel {
	b.NamespaceUID = namespaceUID
	return b
}
