package biz

import (
	"context"

	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/pkg/merr"
	"github.com/bwmarrin/snowflake"
	klog "github.com/go-kratos/kratos/v2/log"
)

func NewMessage(
	messageLogRepo repository.MessageLog,
	messageRepo repository.Message,
	helper *klog.Helper,
) *Message {
	return &Message{
		messageLogRepo: messageLogRepo,
		messageRepo:    messageRepo,
		helper:         klog.NewHelper(klog.With(helper.Logger(), "biz", "message")),
	}
}

type Message struct {
	helper         *klog.Helper
	messageLogRepo repository.MessageLog
	messageRepo    repository.Message
}

func (m *Message) SendMessage(ctx context.Context, uid snowflake.ID) error {
	messageLog, err := m.messageLogRepo.GetMessageLog(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return err
		}
		m.helper.Errorw("msg", "get message log failed", "error", err, "uid", uid)
		return merr.ErrorInternal("get message log failed").WithCause(err)
	}
	if messageLog.Status.IsSent() || messageLog.Status.IsSending() || messageLog.Status.IsCancelled() {
		m.helper.Warnw("msg", "message already sent or sending or cancelled", "uid", uid, "status", messageLog.Status)
		return nil
	}
	return m.messageRepo.SendMessage(ctx, uid)
}
