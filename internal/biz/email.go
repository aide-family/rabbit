package biz

import (
	"context"
	"errors"

	klog "github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"

	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/pkg/merr"
)

func NewEmail(
	emailConfigRepo repository.EmailConfig,
	messageLogRepo repository.MessageLog,
	messageBus repository.MessageBus,
	templateRepo repository.Template,
	helper *klog.Helper,
) *Email {
	return &Email{
		emailConfigRepo: emailConfigRepo,
		messageLogRepo:  messageLogRepo,
		messageBus:      messageBus,
		templateRepo:    templateRepo,
		helper:          klog.NewHelper(klog.With(helper.Logger(), "biz", "email")),
	}
}

type Email struct {
	emailConfigRepo repository.EmailConfig
	messageLogRepo  repository.MessageLog
	messageBus      repository.MessageBus
	templateRepo    repository.Template
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

func (e *Email) AppendEmailMessageWithTemplate(ctx context.Context, req *bo.SendEmailWithTemplateBo) error {
	// 获取模板
	templateDo, err := e.templateRepo.GetTemplate(ctx, req.TemplateUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return merr.ErrorParams("template not found")
		}
		e.helper.Errorw("msg", "get template failed", "error", err)
		return merr.ErrorInternal("get template failed")
	}
	sendEmailBo, err := req.ToSendEmailBo(templateDo)
	if err != nil {
		e.helper.Errorw("msg", "convert template to email template data failed", "error", err)
		return merr.ErrorInternal("convert template to email template data failed")
	}
	return e.AppendEmailMessage(ctx, sendEmailBo)
}
