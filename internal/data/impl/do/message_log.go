package do

import (
	"strings"
	"time"

	"github.com/aide-family/magicbox/strutil"
	"gorm.io/gorm"

	"github.com/aide-family/rabbit/pkg/enum"
)

const (
	TableNameMessageLog = "message_logs"
)

type MessageLog struct {
	BaseModel

	SendAt     time.Time             `gorm:"column:send_at;"`
	Message    strutil.EncryptString `gorm:"column:message;"`
	Config     strutil.EncryptString `gorm:"column:config;"`
	Type       enum.MessageType      `gorm:"column:type;default:0"`
	Status     enum.MessageStatus    `gorm:"column:status;default:0"`
	RetryTotal int32                 `gorm:"column:retry_total;default:0"`
	LastError  string                `gorm:"column:last_error;"`
}

func (m *MessageLog) TableName() string {
	return TableNameMessageLog
}

func GenMessageLogTableName(namespace string, sendAt time.Time) string {
	weekStart := getFirstMonday(sendAt)
	return strings.Join([]string{TableNameMessageLog, namespace, weekStart.Format("20060102")}, "__")
}

func GenMessageLogTableNames(tx *gorm.DB, namespace string, startAt time.Time, endAt time.Time) []string {
	if startAt.After(endAt) {
		return nil
	}
	tableNames := make([]string, 0)
	firstMonday := getFirstMonday(startAt)
	for current := firstMonday; current.Before(endAt); current = current.AddDate(0, 0, 7) {
		if tableName := GenMessageLogTableName(namespace, current); HasTable(tx, tableName) {
			tableNames = append(tableNames, tableName)
		}
	}
	return tableNames
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
