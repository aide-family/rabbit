package bo

import (
	"time"

	"github.com/aide-family/rabbit/internal/biz/vobj"
)

type MessageLogItemBo struct {
	UID     string
	SendAt  time.Time
	Message string
	Type    vobj.MessageType
	Status  vobj.MessageStatus
}

type ListMessageLogBo struct {
	*PageRequestBo
	StartAt time.Time
	EndAt   time.Time
	Keyword string
	Status  vobj.MessageStatus
	Type    vobj.MessageType
}
