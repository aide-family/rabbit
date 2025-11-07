package biz

import (
	"context"
	"errors"

	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/pkg/merr"
	"github.com/bwmarrin/snowflake"
	klog "github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

func NewMessage(
	messageLogRepo repository.MessageLog,
	messageBus repository.MessageBus,
	helper *klog.Helper,
) *Message {
	return &Message{
		messageLogRepo: messageLogRepo,
		messageBus:     messageBus,
		helper:         klog.NewHelper(klog.With(helper.Logger(), "biz", "message")),
	}
}

type Message struct {
	helper         *klog.Helper
	messageLogRepo repository.MessageLog
	messageBus     repository.MessageBus
}

func (m *Message) SendMessage(ctx context.Context, uid snowflake.ID) error {
	messageLog, err := m.messageLogRepo.GetMessageLog(ctx, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return merr.ErrorParams("message log not found")
		}
		m.helper.Errorw("msg", "get message log failed", "error", err, "uid", uid)
		return merr.ErrorInternal("get message log failed")
	}
	if messageLog.Status.IsSent() || messageLog.Status.IsSending() || messageLog.Status.IsCancelled() {
		m.helper.Warnw("msg", "message already sent or sending or cancelled", "uid", uid, "status", messageLog.Status)
		return nil
	}
	return m.messageBus.SendMessage(ctx, uid)
}
