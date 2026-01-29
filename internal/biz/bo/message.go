package bo

import (
	"time"

	"github.com/bwmarrin/snowflake"

	"github.com/aide-family/magicbox/strutil"
	apiv1 "github.com/aide-family/rabbit/pkg/api/v1"
	"github.com/aide-family/rabbit/pkg/enum"
)

type CreateMessageLogBo struct {
	SendAt      time.Time
	Message     string
	MessageType enum.MessageType
	Status      enum.MessageStatus
}

func NewCreateMessageLogBo(sendAt time.Time, message string, messageType enum.MessageType, status enum.MessageStatus) *CreateMessageLogBo {
	return &CreateMessageLogBo{
		SendAt:      sendAt,
		Message:     message,
		MessageType: messageType,
		Status:      status,
	}
}

type MessageLogItemBo struct {
	UID         snowflake.ID
	SendAt      time.Time
	Message     strutil.EncryptString
	Config      strutil.EncryptString
	MessageType enum.MessageType
	Status      enum.MessageStatus
	RetryTotal  int32
	LastError   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (b *MessageLogItemBo) ToAPIV1MessageLogItem() *apiv1.MessageLogItem {
	return &apiv1.MessageLogItem{
		Uid:         b.UID.Int64(),
		MessageType: b.MessageType,
		Status:      enum.MessageStatus(b.Status),
		SendAt:      b.SendAt.Format(time.DateTime),
		Message:     string(b.Message),
		Config:      string(b.Config),
		RetryTotal:  b.RetryTotal,
		LastError:   b.LastError,
		CreatedAt:   b.CreatedAt.Format(time.DateTime),
		UpdatedAt:   b.UpdatedAt.Format(time.DateTime),
	}
}

type ListMessageLogBo struct {
	*PageRequestBo
	StartAt     time.Time
	EndAt       time.Time
	Status      enum.MessageStatus
	MessageType enum.MessageType
}

func NewListMessageLogBo(req *apiv1.ListMessageLogRequest) *ListMessageLogBo {
	return &ListMessageLogBo{
		PageRequestBo: NewPageRequestBo(req.Page, req.PageSize),
		StartAt:       time.Unix(req.StartAtUnix, 0),
		EndAt:         time.Unix(req.EndAtUnix, 0),
		Status:        req.Status,
		MessageType:   req.MessageType,
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
