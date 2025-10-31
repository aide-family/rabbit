package do

import (
	"time"

	"github.com/bwmarrin/snowflake"
)

type MessageRetryLog struct {
	NamespaceModel

	MessageLogID snowflake.ID `gorm:"column:message_log_id;type:bigint(20);not null"`
	RetryAt      time.Time    `gorm:"column:retry_at;type:datetime;not null"`
	Error        string       `gorm:"column:error;type:text;not null"`
}

func (MessageRetryLog) TableName() string {
	return "message_retry_logs"
}
