package convert

import (
	"context"

	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/data/impl/do"
	"github.com/aide-family/rabbit/pkg/contextx"
)

func ToMessageLogItemBo(messageLogDo *do.MessageLog) *bo.MessageLogItemBo {
	return &bo.MessageLogItemBo{
		UID:         messageLogDo.UID,
		SendAt:      messageLogDo.SendAt,
		Message:     messageLogDo.Message,
		Config:      messageLogDo.Config,
		MessageType: messageLogDo.Type,
		Status:      messageLogDo.Status,
		RetryTotal:  messageLogDo.RetryTotal,
		LastError:   messageLogDo.LastError,
		CreatedAt:   messageLogDo.CreatedAt,
		UpdatedAt:   messageLogDo.UpdatedAt,
	}
}

func ToMessageLogDO(ctx context.Context, messageLogBo *bo.MessageLogItemBo) *do.MessageLog {
	model := &do.MessageLog{
		SendAt:     messageLogBo.SendAt,
		Message:    messageLogBo.Message,
		Config:     messageLogBo.Config,
		Type:       messageLogBo.MessageType,
		Status:     messageLogBo.Status,
		RetryTotal: messageLogBo.RetryTotal,
		LastError:  messageLogBo.LastError,
	}
	model.WithNamespaceUID(contextx.GetNamespaceUID(ctx)).WithCreator(contextx.GetUserUID(ctx))
	return model
}
