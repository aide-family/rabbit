// Package model is the model package for the namespace service.
package model

import (
	"errors"
	"time"

	"github.com/aide-family/magicbox/hello"
	"github.com/aide-family/magicbox/safety"

	"github.com/bwmarrin/snowflake"
	"gorm.io/gorm"
)

func Models() []any {
	return []any{
		&Namespace{},
	}
}

type Namespace struct {
	ID        uint32         `gorm:"column:id;primaryKey;autoIncrement"`
	UID       snowflake.ID   `gorm:"column:uid;not null;uniqueIndex"`
	CreatedAt time.Time      `gorm:"column:created_at;type:datetime;not null;"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:datetime;not null;"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetime;index"`
	Creator   snowflake.ID   `gorm:"column:creator;not null;index"`

	Name     string                      `gorm:"column:name;type:varchar(100);not null;uniqueIndex"`
	Metadata *safety.Map[string, string] `gorm:"column:metadata;type:json;"`
	Status   uint8                       `gorm:"column:status;type:tinyint;not null;default:0"`
}

func (Namespace) TableName() string {
	return "namespaces"
}

func (n *Namespace) WithCreator(creator snowflake.ID) *Namespace {
	n.Creator = creator
	return n
}

func (n *Namespace) BeforeCreate(tx *gorm.DB) (err error) {
	if n.Creator == 0 {
		return errors.New("creator is required")
	}

	node, err := snowflake.NewNode(hello.NodeID())
	if err != nil {
		return err
	}
	n.UID = node.Generate()
	if n.Status <= 0 {
		return errors.New("status is required")
	}
	return
}
