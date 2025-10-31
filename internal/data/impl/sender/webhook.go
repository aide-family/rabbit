// Package sender implements message senders for different message types.
package sender

import (
	"context"
	"net/http"

	"github.com/aide-family/rabbit/internal/biz/do"
	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/internal/biz/vobj"
	klog "github.com/go-kratos/kratos/v2/log"
)

// NewWebhookSender 创建Webhook发送器
func NewWebhookSender(helper *klog.Helper) repository.MessageSender {
	return &webhookSender{
		helper: klog.NewHelper(klog.With(helper.Logger(), "impl.sender", "webhook")),
	}
}

type webhookSender struct {
	helper     *klog.Helper
	httpClient *http.Client
}

// Type 返回发送器支持的消息类型
func (w *webhookSender) Type() vobj.MessageType {
	return vobj.MessageTypeWebhook
}

// Send 发送Webhook请求
func (w *webhookSender) Send(ctx context.Context, messageLog *do.MessageLog) error {
	return nil
}
