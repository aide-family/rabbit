package do

import (
	"strings"
	"time"

	"github.com/aide-family/magicbox/strutil"
	"gorm.io/gorm"

	"github.com/aide-family/rabbit/internal/biz/vobj"
)

const (
	TableNameMessageLog = "message_logs"
)

type MessageLog struct {
	NamespaceModel

	SendAt     time.Time             `gorm:"column:send_at;type:datetime;not null"`
	Message    strutil.EncryptString `gorm:"column:message;type:text;not null"`
	Config     strutil.EncryptString `gorm:"column:config;type:text;not null"`
	Type       vobj.MessageType      `gorm:"column:type;type:tinyint(2);not null;default:0"`
	Status     vobj.MessageStatus    `gorm:"column:status;type:tinyint(2);not null;default:0"`
	RetryTotal int32                 `gorm:"column:retry_total;type:int(11);not null;default:0"`
	LastError  string                `gorm:"column:last_error;type:text;not null"`
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
