package biz

import (
	"context"
	"errors"

	klog "github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"

	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/pkg/merr"
)

func NewWebhook(
	webhookConfigBiz *WebhookConfig,
	messageLogBiz *MessageLog,
	eventBusBiz *EventBus,
	templateBiz *Template,
	helper *klog.Helper,
) *Webhook {
	return &Webhook{
		webhookConfigBiz: webhookConfigBiz,
		messageLogBiz:    messageLogBiz,
		eventBusBiz:      eventBusBiz,
		templateBiz:      templateBiz,
		helper:           klog.NewHelper(klog.With(helper.Logger(), "biz", "webhook")),
	}
}

type Webhook struct {
	webhookConfigBiz *WebhookConfig
	messageLogBiz    *MessageLog
	eventBusBiz      *EventBus
	templateBiz      *Template
	helper           *klog.Helper
}

func (w *Webhook) AppendWebhookMessage(ctx context.Context, req *bo.SendWebhookBo) error {
	// 获取webhook配置
	webhookConfig, err := w.webhookConfigBiz.GetWebhook(ctx, req.UID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return merr.ErrorParams("webhook config not found")
		}
		w.helper.Errorw("msg", "get webhook config failed", "error", err)
		return merr.ErrorInternal("get webhook config failed").WithCause(err)
	}
	messageLog, err := req.ToMessageLog(webhookConfig)
	if err != nil {
		w.helper.Errorw("msg", "create message log failed", "error", err)
		return merr.ErrorInternal("generate message log failed").WithCause(err)
	}
	if err := w.messageLogBiz.createMessageLog(ctx, messageLog); err != nil {
		w.helper.Errorw("msg", "create message log failed", "error", err)
		return merr.ErrorInternal("create message log failed").WithCause(err)
	}

	if err := w.eventBusBiz.appendMessage(ctx, messageLog.UID); err != nil {
		w.helper.Errorw("msg", "append webhook message failed", "error", err, "uid", messageLog.UID)
		return merr.ErrorInternal("append webhook message failed").WithCause(err)
	}

	return nil
}

func (w *Webhook) AppendWebhookMessageWithTemplate(ctx context.Context, req *bo.SendWebhookWithTemplateBo) error {
	// 获取模板
	templateDo, err := w.templateBiz.GetTemplate(ctx, req.TemplateUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return merr.ErrorParams("template not found")
		}
		w.helper.Errorw("msg", "get template failed", "error", err)
		return merr.ErrorInternal("get template failed")
	}
	sendWebhookBo, err := req.ToSendWebhookBo(templateDo)
	if err != nil {
		w.helper.Errorw("msg", "convert template to webhook template data failed", "error", err)
		return merr.ErrorInternal("convert template to webhook template data failed")
	}
	return w.AppendWebhookMessage(ctx, sendWebhookBo)
}
