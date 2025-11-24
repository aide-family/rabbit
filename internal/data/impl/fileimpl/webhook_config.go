// Package fileimpl is the implementation of the webhook config repository for file config
package fileimpl

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"

	"github.com/aide-family/magicbox/safety"
	"github.com/aide-family/magicbox/strutil"
	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/do"
	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/internal/biz/vobj"
	"github.com/aide-family/rabbit/internal/conf"
	"github.com/aide-family/rabbit/internal/data"
	"github.com/aide-family/rabbit/pkg/merr"
	"github.com/aide-family/rabbit/pkg/middler"
)

func NewWebhookConfigRepository(d *data.Data) repository.WebhookConfig {
	w := &webhookConfigRepositoryImpl{
		d:              d,
		webhookConfigs: d.GetFileConfig().GetWebhooks(),
	}
	w.initWebhookConfigs()
	d.RegisterReloadFunc(data.KeyWebhooks, func() {
		w.initWebhookConfigs()
	})
	return w
}

type webhookConfigRepositoryImpl struct {
	d                      *data.Data
	webhookConfigs         []*conf.Config_Webhook
	webhookConfigsWithUID  *safety.SyncMap[string, *safety.SyncMap[snowflake.ID, *do.WebhookConfig]]
	webhookConfigsWithName *safety.SyncMap[string, *safety.SyncMap[string, *do.WebhookConfig]]
}

func (w *webhookConfigRepositoryImpl) initWebhookConfigs() {
	w.webhookConfigs = w.d.GetFileConfig().GetWebhooks()
	w.webhookConfigsWithUID = safety.NewSyncMap(make(map[string]*safety.SyncMap[snowflake.ID, *do.WebhookConfig]))
	w.webhookConfigsWithName = safety.NewSyncMap(make(map[string]*safety.SyncMap[string, *do.WebhookConfig]))
	for _, webhookConfig := range w.webhookConfigs {
		namespace := webhookConfig.GetNamespace()
		uid := snowflake.ParseInt64(webhookConfig.GetUid())
		name := webhookConfig.GetName()
		if _, ok := w.webhookConfigsWithUID.Get(namespace); !ok {
			w.webhookConfigsWithUID.Set(namespace, safety.NewSyncMap(map[snowflake.ID]*do.WebhookConfig{}))
			w.webhookConfigsWithName.Set(namespace, safety.NewSyncMap(map[string]*do.WebhookConfig{}))
		}
		item := w.toDoWebhookConfig(webhookConfig)
		if namespaceWebhookConfigsByName, ok := w.webhookConfigsWithName.Get(namespace); ok {
			namespaceWebhookConfigsByName.Set(name, item)
		}
		if namespaceWebhookConfigsByUID, ok := w.webhookConfigsWithUID.Get(namespace); ok {
			namespaceWebhookConfigsByUID.Set(uid, item)
		}
	}
}

func (w *webhookConfigRepositoryImpl) toDoWebhookConfig(webhookConfig *conf.Config_Webhook) *do.WebhookConfig {
	createdAt, _ := time.Parse(time.DateTime, webhookConfig.GetCreatedAt())
	updatedAt, _ := time.Parse(time.DateTime, webhookConfig.GetUpdatedAt())
	headers := safety.NewMap(webhookConfig.GetHeaders())
	return &do.WebhookConfig{
		NamespaceModel: do.NamespaceModel{
			Namespace: webhookConfig.GetNamespace(),
			BaseModel: do.BaseModel{
				ID:        webhookConfig.GetId(),
				UID:       snowflake.ParseInt64(webhookConfig.GetUid()),
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
		},
		App:     vobj.WebhookApp(webhookConfig.GetApp()),
		Name:    webhookConfig.GetName(),
		URL:     webhookConfig.GetUrl(),
		Method:  vobj.HTTPMethod(webhookConfig.GetMethod()),
		Headers: headers,
		Secret:  strutil.EncryptString(webhookConfig.GetSecret()),
		Status:  vobj.GlobalStatus(webhookConfig.GetStatus()),
	}
}

// CreateWebhookConfig implements repository.WebhookConfig.
func (w *webhookConfigRepositoryImpl) CreateWebhookConfig(ctx context.Context, req *do.WebhookConfig) error {
	return merr.ErrorParamsNotSupportFileConfig()
}

// UpdateWebhookConfig implements repository.WebhookConfig.
func (w *webhookConfigRepositoryImpl) UpdateWebhookConfig(ctx context.Context, req *do.WebhookConfig) error {
	return merr.ErrorParamsNotSupportFileConfig()
}

// UpdateWebhookStatus implements repository.WebhookConfig.
func (w *webhookConfigRepositoryImpl) UpdateWebhookStatus(ctx context.Context, uid snowflake.ID, status vobj.GlobalStatus) error {
	return merr.ErrorParamsNotSupportFileConfig()
}

// DeleteWebhookConfig implements repository.WebhookConfig.
func (w *webhookConfigRepositoryImpl) DeleteWebhookConfig(ctx context.Context, uid snowflake.ID) error {
	return merr.ErrorParamsNotSupportFileConfig()
}

// GetWebhookConfig implements repository.WebhookConfig.
func (w *webhookConfigRepositoryImpl) GetWebhookConfig(ctx context.Context, uid snowflake.ID) (*do.WebhookConfig, error) {
	namespace := middler.GetNamespace(ctx)
	webhookConfigWithUID, ok := w.webhookConfigsWithUID.Get(namespace)
	if !ok {
		return nil, merr.ErrorNotFound("webhook config not found")
	}
	webhookConfig, ok := webhookConfigWithUID.Get(uid)
	if !ok {
		return nil, merr.ErrorNotFound("webhook config not found")
	}
	return webhookConfig, nil
}

// GetWebhookConfigByName implements repository.WebhookConfig.
func (w *webhookConfigRepositoryImpl) GetWebhookConfigByName(ctx context.Context, name string) (*do.WebhookConfig, error) {
	namespace := middler.GetNamespace(ctx)
	webhookConfigWithName, ok := w.webhookConfigsWithName.Get(namespace)
	if !ok {
		return nil, merr.ErrorNotFound("webhook config not found")
	}
	webhookConfig, ok := webhookConfigWithName.Get(name)
	if !ok {
		return nil, merr.ErrorNotFound("webhook config not found")
	}
	return webhookConfig, nil
}

// ListWebhookConfig implements repository.WebhookConfig.
func (w *webhookConfigRepositoryImpl) ListWebhookConfig(ctx context.Context, req *bo.ListWebhookBo) (*bo.PageResponseBo[*do.WebhookConfig], error) {
	namespace := middler.GetNamespace(ctx)
	webhookConfigWithUID, ok := w.webhookConfigsWithUID.Get(namespace)
	if !ok {
		pageRequestBo := bo.NewPageRequestBo(req.Page, req.PageSize)
		pageRequestBo.WithTotal(0)
		req.PageRequestBo = pageRequestBo
		return bo.NewPageResponseBo(req.PageRequestBo, []*do.WebhookConfig{}), nil
	}
	webhookConfigs := make([]*do.WebhookConfig, 0, webhookConfigWithUID.Len())
	for _, webhookConfig := range webhookConfigWithUID.Values() {
		if req.App.Exist() && !req.App.IsUnknown() && webhookConfig.App != req.App {
			continue
		}
		if strutil.IsNotEmpty(req.Keyword) && !strings.Contains(webhookConfig.Name, req.Keyword) {
			continue
		}
		webhookConfigs = append(webhookConfigs, webhookConfig)
	}
	total := int64(len(webhookConfigs))
	pageRequestBo := bo.NewPageRequestBo(req.Page, req.PageSize)
	pageRequestBo.WithTotal(total)
	req.PageRequestBo = pageRequestBo
	sort.Slice(webhookConfigs, func(i, j int) bool {
		return webhookConfigs[i].CreatedAt.After(webhookConfigs[j].CreatedAt)
	})
	return bo.NewPageResponseBo(req.PageRequestBo, webhookConfigs), nil
}

// SelectWebhookConfig implements repository.WebhookConfig.
func (w *webhookConfigRepositoryImpl) SelectWebhookConfig(ctx context.Context, req *bo.SelectWebhookBo) (*bo.SelectWebhookResult, error) {
	namespace := middler.GetNamespace(ctx)
	webhookConfigWithUID, ok := w.webhookConfigsWithUID.Get(namespace)
	if !ok {
		return &bo.SelectWebhookResult{
			Items:   []*do.WebhookConfig{},
			Total:   0,
			LastUID: 0,
		}, nil
	}
	webhookConfigs := make([]*do.WebhookConfig, 0, webhookConfigWithUID.Len())
	for _, webhookConfig := range webhookConfigWithUID.Values() {
		if strutil.IsNotEmpty(req.Keyword) && !strings.Contains(webhookConfig.Name, req.Keyword) {
			continue
		}
		if req.Status.Exist() && !req.Status.IsUnknown() && webhookConfig.Status != req.Status {
			continue
		}
		if req.App.Exist() && !req.App.IsUnknown() && webhookConfig.App != req.App {
			continue
		}
		if req.LastUID > 0 && webhookConfig.UID >= req.LastUID {
			continue
		}
		webhookConfigs = append(webhookConfigs, webhookConfig)
	}
	total := int64(len(webhookConfigs))
	sort.Slice(webhookConfigs, func(i, j int) bool {
		return webhookConfigs[i].UID > webhookConfigs[j].UID
	})
	if int32(len(webhookConfigs)) > req.Limit {
		webhookConfigs = webhookConfigs[:req.Limit]
	}
	var lastUID snowflake.ID
	if len(webhookConfigs) > 0 {
		lastUID = webhookConfigs[len(webhookConfigs)-1].UID
	}
	return &bo.SelectWebhookResult{
		Items:   webhookConfigs,
		Total:   total,
		LastUID: lastUID,
	}, nil
}
