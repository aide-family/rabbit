package biz

import (
	"context"

	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/pkg/merr"
	klog "github.com/go-kratos/kratos/v2/log"
)

func NewEmail(
	messageLogRepo repository.MessageLog,
	messageBus repository.MessageBus,
	helper *klog.Helper,
) *Email {
	return &Email{
		messageLogRepo: messageLogRepo,
		messageBus:     messageBus,
		helper:         klog.NewHelper(klog.With(helper.Logger(), "biz", "email")),
	}
}

type Email struct {
	messageLogRepo repository.MessageLog
	messageBus     repository.MessageBus
	helper         *klog.Helper
}

func (e *Email) AppendEmailMessage(ctx context.Context, req *bo.SendEmailBo) error {
	messageLog, err := req.ToMessageLog()
	if err != nil {
		e.helper.Errorw("msg", "create message log failed", "error", err)
		return merr.ErrorInternal("generate message log failed")
	}
	if err := e.messageLogRepo.CreateMessageLog(ctx, messageLog); err != nil {
		e.helper.Errorw("msg", "create message log failed", "error", err)
		return merr.ErrorInternal("create message log failed")
	}

	if err := e.messageBus.AppendMessage(ctx, messageLog); err != nil {
		e.helper.Errorw("msg", "append email message failed", "error", err)
		return merr.ErrorInternal("append email message failed")
	}

	return nil
}
