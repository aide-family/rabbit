package biz

import (
	"context"

	klog "github.com/go-kratos/kratos/v2/log"

	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/pkg/merr"
)

func NewEmail(
	emailConfigBiz *EmailConfig,
	templateBiz *Template,
	messageLogBiz *MessageLog,
	messageBusBiz *MessageBus,
	helper *klog.Helper,
) *Email {
	return &Email{
		emailConfigBiz: emailConfigBiz,
		messageLogBiz:  messageLogBiz,
		messageBusBiz:  messageBusBiz,
		templateBiz:    templateBiz,
		helper:         klog.NewHelper(klog.With(helper.Logger(), "biz", "email")),
	}
}

type Email struct {
	emailConfigBiz *EmailConfig
	templateBiz    *Template
	messageLogBiz  *MessageLog
	messageBusBiz  *MessageBus
	helper         *klog.Helper
}

func (e *Email) AppendEmailMessage(ctx context.Context, req *bo.SendEmailBo) error {
	// 获取邮箱配置
	emailConfig, err := e.emailConfigBiz.GetEmailConfig(ctx, req.UID)
	if err != nil {
		return err
	}
	messageLog, err := req.ToMessageLog(emailConfig)
	if err != nil {
		e.helper.Errorw("msg", "create message log failed", "error", err)
		return merr.ErrorInternal("generate message log failed").WithCause(err)
	}
	if err := e.messageLogBiz.createMessageLog(ctx, messageLog); err != nil {
		e.helper.Errorw("msg", "create message log failed", "error", err)
		return merr.ErrorInternal("create message log failed").WithCause(err)
	}

	if err := e.messageBusBiz.appendMessage(ctx, messageLog.UID); err != nil {
		e.helper.Errorw("msg", "append email message failed", "error", err, "uid", messageLog.UID)
		return merr.ErrorInternal("append email message failed").WithCause(err)
	}

	return nil
}

func (e *Email) AppendEmailMessageWithTemplate(ctx context.Context, req *bo.SendEmailWithTemplateBo) error {
	// 获取模板
	templateBo, err := e.templateBiz.GetTemplate(ctx, req.TemplateUID)
	if err != nil {
		return err
	}
	sendEmailBo, err := req.ToSendEmailBo(templateBo)
	if err != nil {
		e.helper.Errorw("msg", "convert template to email template data failed", "error", err)
		return merr.ErrorInternal("convert template to email template data failed")
	}
	return e.AppendEmailMessage(ctx, sendEmailBo)
}
