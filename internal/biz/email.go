package biz

import (
	"context"
	"errors"

	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/pkg/merr"
	klog "github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

func NewEmail(
	emailConfigRepo repository.EmailConfig,
	messageLogRepo repository.MessageLog,
	messageBus repository.MessageBus,
	helper *klog.Helper,
) *Email {
	return &Email{
		emailConfigRepo: emailConfigRepo,
		messageLogRepo:  messageLogRepo,
		messageBus:      messageBus,
		helper:          klog.NewHelper(klog.With(helper.Logger(), "biz", "email")),
	}
}

type Email struct {
	emailConfigRepo repository.EmailConfig
	messageLogRepo  repository.MessageLog
	messageBus      repository.MessageBus
	helper          *klog.Helper
}

func (e *Email) AppendEmailMessage(ctx context.Context, req *bo.SendEmailBo) error {
	// 获取邮箱配置
	emailConfig, err := e.emailConfigRepo.GetEmailConfig(ctx, req.UID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return merr.ErrorParams("email config not found")
		}
		e.helper.Errorw("msg", "get email config failed", "error", err)
		return merr.ErrorInternal("get email config failed")
	}
	messageLog, err := req.ToMessageLog(emailConfig)
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
