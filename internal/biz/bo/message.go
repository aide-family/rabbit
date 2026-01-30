package bo

import (
	"encoding/json"
	"time"

	"github.com/bwmarrin/snowflake"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/aide-family/magicbox/strutil"
	apiv1 "github.com/aide-family/rabbit/pkg/api/v1"
	"github.com/aide-family/rabbit/pkg/config"
	"github.com/aide-family/rabbit/pkg/enum"
	"github.com/aide-family/rabbit/pkg/merr"
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

// ToMessageConfig 将 bo.MessageLogItemBo 的 Config（BO 的 JSON）转换为 *config.MessageConfig。
// 当前 Config 存储的是各业务 BO 的 JSON（如 WebhookItemBo、EmailConfigItemBo），按 messageType 反序列化后
// 组装为 Driver 所需的 MessageConfig（MessageType + Options *anypb.Any）。
func (b *MessageLogItemBo) ToMessageConfig() (*config.MessageConfig, error) {
	configBytes := []byte(string(b.Config))
	msgType := b.MessageType
	switch {
	case msgType == enum.MessageType_EMAIL:
		var boConfig EmailConfigItemBo
		if err := json.Unmarshal(configBytes, &boConfig); err != nil {
			return nil, merr.ErrorInternal("unmarshal email config failed: %v", err)
		}
		options, err := anypb.New(&config.MessageEmailConfig{
			Host:     boConfig.Host,
			Port:     boConfig.Port,
			Username: boConfig.Username,
			Password: boConfig.Password,
		})
		if err != nil {
			return nil, err
		}
		return &config.MessageConfig{MessageType: msgType, Options: options}, nil
	case msgType >= enum.MessageType_WEBHOOK_OTHER && msgType < 3000:
		var boConfig WebhookItemBo
		if err := json.Unmarshal(configBytes, &boConfig); err != nil {
			return nil, merr.ErrorInternal("unmarshal webhook config failed: %v", err)
		}
		options, err := anypb.New(&config.MessageWebhookConfig{
			App:     boConfig.App,
			Url:     boConfig.URL,
			Secret:  boConfig.Secret,
			Method:  boConfig.Method,
			Headers: boConfig.Headers,
		})
		if err != nil {
			return nil, err
		}
		return &config.MessageConfig{MessageType: msgType, Options: options}, nil
	default:
		return nil, merr.ErrorInternal("unsupported message type for config conversion: %s", msgType)
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
