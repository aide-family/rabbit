package bo

import (
	"time"

	"github.com/bwmarrin/snowflake"

	"github.com/aide-family/magicbox/strutil"
	"github.com/aide-family/rabbit/internal/biz/do"
	"github.com/aide-family/rabbit/internal/biz/vobj"
	apiv1 "github.com/aide-family/rabbit/pkg/api/v1"
	"github.com/aide-family/rabbit/pkg/enum"
)

type CreateMessageLogBo struct {
	SendAt  time.Time
	Message string
	Type    vobj.MessageType
	Status  vobj.MessageStatus
}

func (b *CreateMessageLogBo) ToDoMessageLog() *do.MessageLog {
	return &do.MessageLog{
		SendAt:  b.SendAt,
		Message: strutil.EncryptString(b.Message),
		Type:    b.Type,
		Status:  b.Status,
	}
}

func NewCreateMessageLogBo(sendAt time.Time, message string, messageType vobj.MessageType, status vobj.MessageStatus) *CreateMessageLogBo {
	return &CreateMessageLogBo{
		SendAt:  sendAt,
		Message: message,
		Type:    messageType,
		Status:  status,
	}
}

type MessageLogItemBo struct {
	UID        snowflake.ID
	SendAt     time.Time
	Message    strutil.EncryptString
	Config     strutil.EncryptString
	Type       vobj.MessageType
	Status     vobj.MessageStatus
	RetryTotal int32
	LastError  string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewMessageLogItemBo(doMessageLog *do.MessageLog) *MessageLogItemBo {
	return &MessageLogItemBo{
		UID:        doMessageLog.UID,
		SendAt:     doMessageLog.SendAt,
		Message:    doMessageLog.Message,
		Config:     doMessageLog.Config,
		Type:       doMessageLog.Type,
		Status:     doMessageLog.Status,
		RetryTotal: doMessageLog.RetryTotal,
		LastError:  doMessageLog.LastError,
		CreatedAt:  doMessageLog.CreatedAt,
		UpdatedAt:  doMessageLog.UpdatedAt,
	}
}

func (b *MessageLogItemBo) ToAPIV1MessageLogItem() *apiv1.MessageLogItem {
	return &apiv1.MessageLogItem{
		Uid:        b.UID.Int64(),
		Type:       enum.MessageType(b.Type),
		Status:     enum.MessageStatus(b.Status),
		SendAt:     b.SendAt.Format(time.DateTime),
		Message:    string(b.Message),
		Config:     string(b.Config),
		RetryTotal: b.RetryTotal,
		LastError:  b.LastError,
		CreatedAt:  b.CreatedAt.Format(time.DateTime),
		UpdatedAt:  b.UpdatedAt.Format(time.DateTime),
	}
}

type ListMessageLogBo struct {
	*PageRequestBo
	StartAt time.Time
	EndAt   time.Time
	Status  vobj.MessageStatus
	Type    vobj.MessageType
}

func NewListMessageLogBo(req *apiv1.ListMessageLogRequest) *ListMessageLogBo {
	return &ListMessageLogBo{
		PageRequestBo: NewPageRequestBo(req.Page, req.PageSize),
		StartAt:       time.Unix(req.StartAtUnix, 0),
		EndAt:         time.Unix(req.EndAtUnix, 0),
		Status:        vobj.MessageStatus(req.Status),
		Type:          vobj.MessageType(req.Type),
	}
}

func ToAPIV1ListMessageLogReply(pageResponseBo *PageResponseBo[*MessageLogItemBo]) *apiv1.ListMessageLogReply {
	items := make([]*apiv1.MessageLogItem, 0, len(pageResponseBo.GetItems()))
	for _, item := range pageResponseBo.GetItems() {
		items = append(items, item.ToAPIV1MessageLogItem())
	}
	return &apiv1.ListMessageLogReply{
		Items:    items,
		Total:    pageResponseBo.GetTotal(),
		Page:     pageResponseBo.GetPage(),
		PageSize: pageResponseBo.GetPageSize(),
	}
}
