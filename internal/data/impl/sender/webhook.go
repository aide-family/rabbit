// Package sender implements message senders for different message types.
package sender

import (
	"context"
	"strings"

	"github.com/aide-family/magicbox/message"
	"github.com/aide-family/magicbox/message/hook"
	"github.com/aide-family/magicbox/message/hook/dingtalk"
	"github.com/aide-family/magicbox/message/hook/feishu"
	"github.com/aide-family/magicbox/message/hook/wechat"
	"github.com/aide-family/magicbox/safety"
	"github.com/aide-family/magicbox/serialize"
	"github.com/aide-family/magicbox/strutil"
	klog "github.com/go-kratos/kratos/v2/log"

	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/internal/biz/vobj"
	"github.com/aide-family/rabbit/pkg/merr"
)

// NewWebhookSender 创建Webhook发送器
func NewWebhookSender(helper *klog.Helper) repository.MessageSender {
	w := &webhookSender{
		helper:     klog.NewHelper(klog.With(helper.Logger(), "impl.sender", "webhook")),
		senders:    safety.NewSyncMap(make(map[int64]message.Sender)),
		sendHashes: safety.NewSyncMap(make(map[int64]string)),
		drivers:    safety.NewSyncMap(make(map[vobj.WebhookApp]func(hook.Config) message.Driver)),
	}
	w.drivers.Set(vobj.WebhookAppDingTalk, dingtalk.SenderDriver)
	w.drivers.Set(vobj.WebhookAppWechat, wechat.SenderDriver)
	w.drivers.Set(vobj.WebhookAppFeishu, feishu.SenderDriver)
	return w
}

type webhookSender struct {
	helper *klog.Helper

	senders    *safety.SyncMap[int64, message.Sender]
	sendHashes *safety.SyncMap[int64, string]

	drivers *safety.SyncMap[vobj.WebhookApp, func(hook.Config) message.Driver]
}

// Type 返回发送器支持的消息类型
func (w *webhookSender) Type() vobj.MessageType {
	return vobj.MessageTypeWebhook
}

// Send 发送Webhook请求
func (w *webhookSender) Send(ctx context.Context, messageLog *bo.MessageLogItemBo) error {
	webhookMessage, err := w.buildWebhookMessage([]byte(string(messageLog.Message)))
	if err != nil {
		return err
	}
	webhookSender, err := w.getSender([]byte(string(messageLog.Config)))
	if err != nil {
		return err
	}

	if err := webhookSender.Send(ctx, webhookMessage); err != nil {
		return merr.ErrorInternal("send webhook message failed").WithCause(err)
	}
	return nil
}

func (w *webhookSender) buildWebhookMessage(messageBytes []byte) (message.Message, error) {
	var webhookMessage bo.SendWebhookBo
	if err := serialize.JSONUnmarshal(messageBytes, &webhookMessage); err != nil {
		return nil, merr.ErrorInternal("unmarshal webhook message failed").WithCause(err)
	}
	return &webhookMessage, nil
}

func (w *webhookSender) getSender(configBytes []byte) (message.Sender, error) {
	var webhookConfig bo.WebhookItemBo
	if err := serialize.JSONUnmarshal(configBytes, &webhookConfig); err != nil {
		return nil, merr.ErrorInternal("unmarshal webhook config failed").WithCause(err)
	}
	sendHash := strutil.SHA256(string(configBytes))
	hash, ok := w.sendHashes.Get(webhookConfig.UID.Int64())
	if ok && strings.EqualFold(sendHash, hash) {
		sender, ok := w.senders.Get(webhookConfig.UID.Int64())
		if !ok {
			return nil, merr.ErrorParams("webhook sender not found")
		}
		return sender, nil
	}

	driverFunc, ok := w.drivers.Get(webhookConfig.App)
	if !ok {
		return nil, merr.ErrorParams("invalid webhook app type, expected webhook type, got %s", webhookConfig.App)
	}

	sender, err := message.NewSender(driverFunc(&webhookConfig))
	if err != nil {
		return nil, merr.ErrorInternal("create webhook sender failed").WithCause(err)
	}
	w.senders.Set(webhookConfig.UID.Int64(), sender)
	w.sendHashes.Set(webhookConfig.UID.Int64(), sendHash)
	return sender, nil
}
