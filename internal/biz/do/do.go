// Package do is the data object package for the Rabbit service.
package do

import (
	"time"

	"gorm.io/gorm"
)

func Models() []any {
	return []any{
		&Namespace{},
	}
}

type BaseModel struct {
	ID        uint32         `gorm:"column:id;type:int unsigned;primaryKey;autoIncrement"`
	CreatedAt time.Time      `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetime;index"`
}
