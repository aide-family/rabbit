// Package sender implements message senders for different message types.
package sender

import (
	"context"
	"strings"

	"github.com/aide-family/magicbox/message"
	"github.com/aide-family/magicbox/message/email"
	"github.com/aide-family/magicbox/safety"
	"github.com/aide-family/magicbox/serialize"
	"github.com/aide-family/magicbox/strutil"
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
		helper:     klog.NewHelper(klog.With(helper.Logger(), "impl.sender", "email")),
		senders:    safety.NewSyncMap(make(map[int64]message.Sender)),
		sendHashes: safety.NewSyncMap(make(map[int64]string)),
	}
}

type emailSender struct {
	helper     *klog.Helper
	senders    *safety.SyncMap[int64, message.Sender]
	sendHashes *safety.SyncMap[int64, string]
}

// Type 返回发送器支持的消息类型
func (e *emailSender) Type() vobj.MessageType {
	return vobj.MessageTypeEmail
}

// Send 发送邮件
func (e *emailSender) Send(ctx context.Context, messageLog *do.MessageLog) error {
	msg, err := e.buildEmailMessage([]byte(string(messageLog.Message)))
	if err != nil {
		return merr.ErrorInternal("convert to email message failed").WithCause(err)
	}
	sender, err := e.getSender([]byte(string(messageLog.Config)))
	if err != nil {
		return err
	}
	if err := sender.Send(ctx, msg); err != nil {
		e.helper.Errorw("msg", "send email failed", "error", err, "uid", messageLog.UID)
		return merr.ErrorInternal("send email failed").WithCause(err)
	}
	return nil
}

func (e *emailSender) buildEmailMessage(messageBytes []byte) (*email.Message, error) {
	var emailMessage bo.SendEmailBo
	if err := serialize.JSONUnmarshal(messageBytes, &emailMessage); err != nil {
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

func (e *emailSender) getSender(emailConfigBytes []byte) (message.Sender, error) {
	var emailConfig bo.EmailConfigItemBo
	if err := serialize.JSONUnmarshal(emailConfigBytes, &emailConfig); err != nil {
		e.helper.Errorw("msg", "unmarshal email config failed", "error", err)
		return nil, merr.ErrorInternal("unmarshal email config failed")
	}
	sendHash := strutil.SHA256(string(emailConfigBytes))
	hash, ok := e.sendHashes.Get(emailConfig.UID.Int64())
	if !ok || !strings.EqualFold(sendHash, hash) {
		sender, err := message.NewSender(email.SenderDriver(&emailConfig))
		if err != nil {
			e.helper.Errorw("msg", "create email sender failed", "error", err)
			return nil, merr.ErrorInternal("create email sender failed")
		}
		e.senders.Set(emailConfig.UID.Int64(), sender)
		e.sendHashes.Set(emailConfig.UID.Int64(), sendHash)
	}
	sender, ok := e.senders.Get(emailConfig.UID.Int64())
	if !ok {
		e.helper.Errorw("msg", "email sender not found", "uid", emailConfig.UID)
		return nil, merr.ErrorParams("email sender not found")
	}
	return sender, nil
}
