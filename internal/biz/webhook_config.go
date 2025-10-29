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

func NewWebhookConfig(
	webhookConfigRepo repository.WebhookConfig,
	webhookTemplateRepo repository.WebhookTemplate,
	namespace string,
	helper *klog.Helper,
) *WebhookConfig {
	return &WebhookConfig{
		webhookConfigRepo:   webhookConfigRepo,
		webhookTemplateRepo: webhookTemplateRepo,
		helper:              klog.NewHelper(klog.With(helper.Logger(), "biz", "webhookConfig")),
	}
}

type WebhookConfig struct {
	helper              *klog.Helper
	webhookConfigRepo   repository.WebhookConfig
	webhookTemplateRepo repository.WebhookTemplate
}

func (w *WebhookConfig) CreateWebhook(ctx context.Context, req *bo.CreateWebhookBo) error {
	doWebhookConfig := req.ToDoWebhookConfig()
	if err := w.webhookConfigRepo.SaveWebhookConfig(ctx, doWebhookConfig); err != nil {
		w.helper.Errorw("msg", "create webhook config failed", "error", err, "name", doWebhookConfig.Name)
		return merr.ErrorInternal("create webhook config %s failed", doWebhookConfig.Name)
	}
	return nil
}

func (w *WebhookConfig) UpdateWebhook(ctx context.Context, req *bo.UpdateWebhookBo) error {
	doWebhookConfig := req.ToDoWebhookConfig()
	if err := w.webhookConfigRepo.SaveWebhookConfig(ctx, doWebhookConfig); err != nil {
		w.helper.Errorw("msg", "update webhook config failed", "error", err, "uid", doWebhookConfig.UID)
		return merr.ErrorInternal("update webhook config %s failed", doWebhookConfig.UID)
	}
	return nil
}

func (w *WebhookConfig) UpdateWebhookStatus(ctx context.Context, req *bo.UpdateWebhookStatusBo) error {
	if err := w.webhookConfigRepo.UpdateWebhookStatus(ctx, req.UID, req.Status); err != nil {
		w.helper.Errorw("msg", "update webhook status failed", "error", err, "uid", req.UID)
		return merr.ErrorInternal("update webhook status %s failed", req.UID)
	}
	return nil
}

func (w *WebhookConfig) DeleteWebhook(ctx context.Context, uid string) error {
	if err := w.webhookConfigRepo.DeleteWebhookConfig(ctx, uid); err != nil {
		w.helper.Errorw("msg", "delete webhook config failed", "error", err, "uid", uid)
		return merr.ErrorInternal("delete webhook config %s failed", uid)
	}
	return nil
}

func (w *WebhookConfig) GetWebhook(ctx context.Context, uid string) (*bo.WebhookItemBo, error) {
	doWebhookConfig, err := w.webhookConfigRepo.GetWebhookConfig(ctx, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, merr.ErrorNotFound("webhook config %s not found", uid)
		}
		w.helper.Errorw("msg", "get webhook config failed", "error", err, "uid", uid)
		return nil, merr.ErrorInternal("get webhook config %s failed", uid)
	}
	return bo.NewWebhookItemBo(doWebhookConfig), nil
}

func (w *WebhookConfig) ListWebhook(ctx context.Context, req *bo.ListWebhookBo) (*bo.PageResponseBo[*bo.WebhookItemBo], error) {
	pageResponseBo, err := w.webhookConfigRepo.ListWebhookConfig(ctx, req)
	if err != nil {
		w.helper.Errorw("msg", "list webhook config failed", "error", err, "req", req)
		return nil, merr.ErrorInternal("list webhook config failed")
	}
	items := make([]*bo.WebhookItemBo, 0, len(pageResponseBo.GetItems()))
	for _, item := range pageResponseBo.GetItems() {
		items = append(items, bo.NewWebhookItemBo(item))
	}
	return bo.NewPageResponseBo(pageResponseBo.PageRequestBo, items), nil
}

func (w *WebhookConfig) CreateWebhookTemplate(ctx context.Context, req *bo.CreateWebhookTemplateBo) error {
	doWebhookTemplate := req.ToDoWebhookTemplate()
	if err := w.webhookTemplateRepo.SaveWebhookTemplate(ctx, doWebhookTemplate); err != nil {
		w.helper.Errorw("msg", "create webhook template failed", "error", err, "name", doWebhookTemplate.Name)
		return merr.ErrorInternal("create webhook template %s failed", doWebhookTemplate.Name)
	}
	return nil
}

func (w *WebhookConfig) UpdateWebhookTemplate(ctx context.Context, req *bo.UpdateWebhookTemplateBo) error {
	doWebhookTemplate := req.ToDoWebhookTemplate()
	if err := w.webhookTemplateRepo.SaveWebhookTemplate(ctx, doWebhookTemplate); err != nil {
		w.helper.Errorw("msg", "update webhook template failed", "error", err, "uid", doWebhookTemplate.UID)
		return merr.ErrorInternal("update webhook template %s failed", doWebhookTemplate.UID)
	}
	return nil
}

func (w *WebhookConfig) UpdateWebhookTemplateStatus(ctx context.Context, req *bo.UpdateWebhookTemplateStatusBo) error {
	if err := w.webhookTemplateRepo.UpdateWebhookTemplateStatus(ctx, req.UID, req.Status); err != nil {
		w.helper.Errorw("msg", "update webhook template status failed", "error", err, "uid", req.UID)
		return merr.ErrorInternal("update webhook template status %s failed", req.UID)
	}
	return nil
}

func (w *WebhookConfig) DeleteWebhookTemplate(ctx context.Context, uid string) error {
	if err := w.webhookTemplateRepo.DeleteWebhookTemplate(ctx, uid); err != nil {
		w.helper.Errorw("msg", "delete webhook template failed", "error", err, "uid", uid)
		return merr.ErrorInternal("delete webhook template %s failed", uid)
	}
	return nil
}

func (w *WebhookConfig) GetWebhookTemplate(ctx context.Context, uid string) (*bo.WebhookTemplateItemBo, error) {
	doWebhookTemplate, err := w.webhookTemplateRepo.GetWebhookTemplate(ctx, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, merr.ErrorNotFound("webhook template %s not found", uid)
		}
		w.helper.Errorw("msg", "get webhook template failed", "error", err, "uid", uid)
		return nil, merr.ErrorInternal("get webhook template %s failed", uid)
	}
	return bo.NewWebhookTemplateItemBo(doWebhookTemplate), nil
}

func (w *WebhookConfig) ListWebhookTemplate(ctx context.Context, req *bo.ListWebhookTemplateBo) (*bo.PageResponseBo[*bo.WebhookTemplateItemBo], error) {
	pageResponseBo, err := w.webhookTemplateRepo.ListWebhookTemplate(ctx, req)
	if err != nil {
		w.helper.Errorw("msg", "list webhook template failed", "error", err, "req", req)
		return nil, merr.ErrorInternal("list webhook template failed")
	}
	items := make([]*bo.WebhookTemplateItemBo, 0, len(pageResponseBo.GetItems()))
	for _, item := range pageResponseBo.GetItems() {
		items = append(items, bo.NewWebhookTemplateItemBo(item))
	}
	return bo.NewPageResponseBo(pageResponseBo.PageRequestBo, items), nil
}
