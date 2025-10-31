// Package sender implements message senders for different message types.
package sender

import (
	"context"

	"github.com/aide-family/rabbit/internal/biz/do"
	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/internal/biz/vobj"
	klog "github.com/go-kratos/kratos/v2/log"
)

// NewEmailSender 创建邮件发送器
func NewEmailSender(helper *klog.Helper) repository.MessageSender {
	return &emailSender{
		helper: klog.NewHelper(klog.With(helper.Logger(), "impl.sender", "email")),
	}
}

type emailSender struct {
	helper *klog.Helper
}

// Type 返回发送器支持的消息类型
func (e *emailSender) Type() vobj.MessageType {
	return vobj.MessageTypeEmail
}

// Send 发送邮件
func (e *emailSender) Send(ctx context.Context, messageLog *do.MessageLog) error {
	return nil
}
