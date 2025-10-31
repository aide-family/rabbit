package do

import (
	"strings"
	"time"

	"github.com/aide-family/rabbit/internal/biz/vobj"
	"gorm.io/gorm"
)

const (
	tableNameMessageLog = "message_logs"
)

type MessageLog struct {
	NamespaceModel

	SendAt  time.Time          `gorm:"column:send_at;type:datetime;not null"`
	Message string             `gorm:"column:message;type:text;not null"`
	Type    vobj.MessageType   `gorm:"column:type;type:tinyint(2);not null;default:0"`
	Status  vobj.MessageStatus `gorm:"column:status;type:tinyint(2);not null;default:0"`
}

func (m *MessageLog) TableName() string {
	return GenMessageLogTableName(m.Namespace, m.SendAt)
}

func GenMessageLogTableName(namespace string, sendAt time.Time) string {
	weekStart := getFirstMonday(sendAt)
	return strings.Join([]string{tableNameMessageLog, namespace, weekStart.Format("20060102")}, "__")
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
