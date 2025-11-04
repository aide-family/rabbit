// Package sender implements message senders for different message types.
package sender

import (
	"context"

	"github.com/aide-family/magicbox/message"
	"github.com/aide-family/magicbox/message/email"
	"github.com/aide-family/magicbox/safety"
	"github.com/aide-family/magicbox/serialize"
	klog "github.com/go-kratos/kratos/v2/log"

	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/do"
	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/internal/biz/vobj"
	"github.com/aide-family/rabbit/pkg/merr"
)

// NewEmailSender 创建邮件发送器
func NewEmailSender(helper *klog.Helper) repository.MessageSender {
	return &emailSender{
		helper:  klog.NewHelper(klog.With(helper.Logger(), "impl.sender", "email")),
		senders: safety.NewSyncMap(make(map[int64]message.Sender)),
	}
}

type emailSender struct {
	helper  *klog.Helper
	senders *safety.SyncMap[int64, message.Sender]
}

// Type 返回发送器支持的消息类型
func (e *emailSender) Type() vobj.MessageType {
	return vobj.MessageTypeEmail
}

// Send 发送邮件
func (e *emailSender) Send(ctx context.Context, messageLog *do.MessageLog) error {
	msg, err := e.buildEmailMessage(messageLog)
	if err != nil {
		return merr.ErrorInternal("convert to email message failed").WithCause(err)
	}
	sender, ok := e.senders.Get(messageLog.UID.Int64())
	if !ok {
		e.helper.Errorw("msg", "email sender not found", "uid", messageLog.UID)
		return merr.ErrorInternal("email sender not found")
	}
	if err := sender.Send(ctx, msg); err != nil {
		e.helper.Errorw("msg", "send email failed", "error", err, "uid", messageLog.UID)
		return merr.ErrorInternal("send email failed").WithCause(err)
	}
	return nil
}

func (e *emailSender) buildEmailMessage(messageLog *do.MessageLog) (*email.Message, error) {
	var emailConfig bo.EmailConfigItemBo
	if err := serialize.JSONUnmarshal([]byte(messageLog.Config), &emailConfig); err != nil {
		e.helper.Errorw("msg", "unmarshal email config failed", "error", err)
		return nil, merr.ErrorInternal("unmarshal email config failed")
	}
	sender, err := message.NewSender(email.SenderDriver(&emailConfig))
	if err != nil {
		e.helper.Errorw("msg", "create email sender failed", "error", err)
		return nil, merr.ErrorInternal("create email sender failed")
	}
	e.senders.Set(emailConfig.UID.Int64(), sender)
	var emailMessage bo.SendEmailBo
	if err := serialize.JSONUnmarshal([]byte(messageLog.Message), &emailMessage); err != nil {
		e.helper.Errorw("msg", "unmarshal email message failed", "error", err)
		return nil, merr.ErrorInternal("unmarshal email message failed")
	}

	return &email.Message{
		To:          emailMessage.To,
		Cc:          emailMessage.Cc,
		Subject:     emailMessage.Subject,
		Body:        emailMessage.Body,
		ContentType: emailMessage.ContentType,
		Attachments: []*email.Attachment{},
		Headers:     emailMessage.Headers,
	}, nil
}
