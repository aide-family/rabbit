package biz

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/aide-family/magicbox/safety"
	"github.com/aide-family/magicbox/strutil"
	"github.com/bwmarrin/snowflake"
	klog "github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"

	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/internal/biz/vobj"
	"github.com/aide-family/rabbit/internal/conf"
	"github.com/aide-family/rabbit/pkg/merr"
	"github.com/aide-family/rabbit/pkg/middler"
)

func NewWebhookConfig(
	bc *conf.Bootstrap,
	webhookConfigRepo repository.WebhookConfig,
	helper *klog.Helper,
) *WebhookConfig {
	w := &WebhookConfig{
		useDatabase:       bc.GetUseDatabase() == "true",
		webhookConfigRepo: webhookConfigRepo,
		helper:            klog.NewHelper(klog.With(helper.Logger(), "biz", "webhookConfig")),
	}
	w.webhookConfigs = w.initWebhookConfigs()
	conf.RegisterReloadFunc(conf.KeyWebhooks, func() {
		w.helper.Infow("msg", "webhooks changed, reloading webhooks")
		w.webhookConfigs = w.initWebhookConfigs()
	})
	return w
}

type WebhookConfig struct {
	helper            *klog.Helper
	useDatabase       bool
	webhookConfigRepo repository.WebhookConfig
	webhookConfigs    *safety.SyncMap[string, *safety.SyncMap[snowflake.ID, *bo.WebhookItemBo]]
}

func (w *WebhookConfig) initWebhookConfigs() *safety.SyncMap[string, *safety.SyncMap[snowflake.ID, *bo.WebhookItemBo]] {
	webhookConfigs := safety.NewSyncMap(make(map[string]*safety.SyncMap[snowflake.ID, *bo.WebhookItemBo]))
	for _, webhookConfig := range conf.GetFileConfig().GetWebhooks() {
		namespace := webhookConfig.GetNamespace()
		uid := snowflake.ParseInt64(webhookConfig.GetUid())
		if _, ok := webhookConfigs.Get(namespace); !ok {
			webhookConfigs.Set(namespace, safety.NewSyncMap(map[snowflake.ID]*bo.WebhookItemBo{}))
		}
		namespaceWebhookConfigs, ok := webhookConfigs.Get(namespace)
		if !ok {
			continue
		}
		createdAt, _ := time.Parse(time.DateTime, webhookConfig.GetCreatedAt())
		updatedAt, _ := time.Parse(time.DateTime, webhookConfig.GetUpdatedAt())
		namespaceWebhookConfigs.Set(uid, &bo.WebhookItemBo{
			UID:       uid,
			App:       vobj.WebhookApp(webhookConfig.GetApp()),
			Name:      webhookConfig.GetName(),
			URL:       webhookConfig.GetUrl(),
			Method:    vobj.HTTPMethod(webhookConfig.GetMethod()),
			Headers:   webhookConfig.GetHeaders(),
			Secret:    webhookConfig.GetSecret(),
			Status:    vobj.GlobalStatus(webhookConfig.GetStatus()),
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
	}
	return webhookConfigs
}

func (w *WebhookConfig) CreateWebhook(ctx context.Context, req *bo.CreateWebhookBo) error {
	if !w.useDatabase {
		return merr.ErrorParamsNotSupportFileConfig()
	}
	doWebhookConfig := req.ToDoWebhookConfig()
	if _, err := w.webhookConfigRepo.GetWebhookConfigByName(ctx, doWebhookConfig.Name); err == nil {
		return merr.ErrorParams("webhook config %s already exists", doWebhookConfig.Name)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		w.helper.Errorw("msg", "check webhook config exists failed", "error", err, "name", doWebhookConfig.Name)
		return merr.ErrorInternal("create webhook config %s failed", doWebhookConfig.Name)
	}
	if err := w.webhookConfigRepo.CreateWebhookConfig(ctx, doWebhookConfig); err != nil {
		w.helper.Errorw("msg", "create webhook config failed", "error", err, "name", doWebhookConfig.Name)
		return merr.ErrorInternal("create webhook config %s failed", doWebhookConfig.Name)
	}
	return nil
}

func (w *WebhookConfig) UpdateWebhook(ctx context.Context, req *bo.UpdateWebhookBo) error {
	if !w.useDatabase {
		return merr.ErrorParamsNotSupportFileConfig()
	}
	doWebhookConfig := req.ToDoWebhookConfig()
	if err := w.webhookConfigRepo.UpdateWebhookConfig(ctx, doWebhookConfig); err != nil {
		w.helper.Errorw("msg", "update webhook config failed", "error", err, "uid", doWebhookConfig.UID)
		return merr.ErrorInternal("update webhook config %s failed", doWebhookConfig.UID)
	}
	return nil
}

func (w *WebhookConfig) UpdateWebhookStatus(ctx context.Context, req *bo.UpdateWebhookStatusBo) error {
	if !w.useDatabase {
		return merr.ErrorParamsNotSupportFileConfig()
	}
	if err := w.webhookConfigRepo.UpdateWebhookStatus(ctx, req.UID, req.Status); err != nil {
		w.helper.Errorw("msg", "update webhook status failed", "error", err, "uid", req.UID)
		return merr.ErrorInternal("update webhook status %s failed", req.UID)
	}
	return nil
}

func (w *WebhookConfig) DeleteWebhook(ctx context.Context, uid snowflake.ID) error {
	if !w.useDatabase {
		return merr.ErrorParamsNotSupportFileConfig()
	}
	if err := w.webhookConfigRepo.DeleteWebhookConfig(ctx, uid); err != nil {
		w.helper.Errorw("msg", "delete webhook config failed", "error", err, "uid", uid)
		return merr.ErrorInternal("delete webhook config %s failed", uid)
	}
	return nil
}

func (w *WebhookConfig) getWebhookByFileConfigWithUID(ctx context.Context, uid snowflake.ID) (*bo.WebhookItemBo, error) {
	namespaceWebhookConfigs, ok := w.webhookConfigs.Get(middler.GetNamespace(ctx))
	if !ok {
		return nil, merr.ErrorNotFound("webhook config %s not found", uid)
	}
	webhookConfig, ok := namespaceWebhookConfigs.Get(uid)
	if !ok {
		return nil, merr.ErrorNotFound("webhook config %s not found", uid)
	}
	return webhookConfig, nil
}

func (w *WebhookConfig) GetWebhook(ctx context.Context, uid snowflake.ID) (*bo.WebhookItemBo, error) {
	if !w.useDatabase {
		return w.getWebhookByFileConfigWithUID(ctx, uid)
	}
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

func (w *WebhookConfig) getWebhookByFileConfigWithNamespace(ctx context.Context, req *bo.ListWebhookBo) ([]*bo.WebhookItemBo, error) {
	namespace := middler.GetNamespace(ctx)
	namespaceWebhookConfigs, ok := w.webhookConfigs.Get(namespace)
	if !ok {
		return nil, merr.ErrorNotFound("webhook config %s not found", namespace)
	}
	webhookConfigs := make([]*bo.WebhookItemBo, 0, namespaceWebhookConfigs.Len())
	for _, webhookConfig := range namespaceWebhookConfigs.Values() {
		if strutil.IsNotEmpty(req.Keyword) && !strings.Contains(webhookConfig.Name, req.Keyword) {
			continue
		}
		if req.App.Exist() && !req.App.IsUnknown() && webhookConfig.App != req.App {
			continue
		}
		webhookConfigs = append(webhookConfigs, webhookConfig)
	}
	total := int64(len(webhookConfigs))
	pageRequestBo := bo.NewPageRequestBo(req.Page, req.PageSize)
	pageRequestBo.WithTotal(total)
	req.PageRequestBo = pageRequestBo
	sort.Slice(webhookConfigs, func(i, j int) bool {
		return webhookConfigs[i].UID < webhookConfigs[j].UID
	})
	return webhookConfigs, nil
}

func (w *WebhookConfig) ListWebhook(ctx context.Context, req *bo.ListWebhookBo) (*bo.PageResponseBo[*bo.WebhookItemBo], error) {
	if !w.useDatabase {
		webhookConfigs, err := w.getWebhookByFileConfigWithNamespace(ctx, req)
		if err != nil {
			return nil, err
		}
		return bo.NewPageResponseBo(req.PageRequestBo, webhookConfigs), nil
	}
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

func (w *WebhookConfig) getSelectWebhookByFileConfig(ctx context.Context, req *bo.SelectWebhookBo) ([]*bo.WebhookItemSelectBo, error) {
	namespace := middler.GetNamespace(ctx)
	namespaceWebhookConfigs, ok := w.webhookConfigs.Get(namespace)
	if !ok {
		return nil, merr.ErrorNotFound("webhook config %s not found", namespace)
	}
	webhookConfigs := make([]*bo.WebhookItemSelectBo, 0, namespaceWebhookConfigs.Len())
	for _, webhookConfig := range namespaceWebhookConfigs.Values() {
		if strutil.IsNotEmpty(req.Keyword) && !strings.Contains(webhookConfig.Name, req.Keyword) {
			continue
		}
		if req.Status.Exist() && !req.Status.IsUnknown() && webhookConfig.Status != req.Status {
			continue
		}
		if req.App.Exist() && !req.App.IsUnknown() && webhookConfig.App != req.App {
			continue
		}
		webhookConfigs = append(webhookConfigs, &bo.WebhookItemSelectBo{
			UID:      webhookConfig.UID,
			Name:     webhookConfig.Name,
			Status:   webhookConfig.Status,
			Disabled: webhookConfig.Status != vobj.GlobalStatusEnabled,
			Tooltip:  "",
		})
	}
	total := int64(len(webhookConfigs))
	req.Limit = int32(total)
	req.LastUID = 0
	sort.Slice(webhookConfigs, func(i, j int) bool {
		return webhookConfigs[i].UID < webhookConfigs[j].UID
	})
	return webhookConfigs, nil
}

func (w *WebhookConfig) SelectWebhook(ctx context.Context, req *bo.SelectWebhookBo) (*bo.SelectWebhookBoResult, error) {
	if !w.useDatabase {
		webhookConfigs, err := w.getSelectWebhookByFileConfig(ctx, req)
		if err != nil {
			return nil, err
		}
		return &bo.SelectWebhookBoResult{
			Items:   webhookConfigs,
			Total:   int64(len(webhookConfigs)),
			LastUID: 0,
		}, nil
	}
	result, err := w.webhookConfigRepo.SelectWebhookConfig(ctx, req)
	if err != nil {
		w.helper.Errorw("msg", "select webhook config failed", "error", err, "req", req)
		return nil, merr.ErrorInternal("select webhook config failed")
	}
	items := make([]*bo.WebhookItemSelectBo, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, bo.NewWebhookItemSelectBo(item))
	}
	return &bo.SelectWebhookBoResult{
		Items:   items,
		Total:   result.Total,
		LastUID: result.LastUID,
	}, nil
}
