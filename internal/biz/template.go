package biz

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
	klog "github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"

	"github.com/aide-family/magicbox/safety"
	"github.com/aide-family/magicbox/strutil"
	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/internal/biz/vobj"
	"github.com/aide-family/rabbit/internal/conf"
	"github.com/aide-family/rabbit/pkg/merr"
	"github.com/aide-family/rabbit/pkg/middler"
)

func NewTemplate(
	bc *conf.Bootstrap,
	templateRepo repository.Template,
	helper *klog.Helper,
) *Template {
	templates := safety.NewSyncMap(make(map[string]*safety.SyncMap[snowflake.ID, *bo.TemplateItemBo]))
	for _, template := range conf.GetFileConfig().GetTemplates() {
		namespace := template.GetNamespace()
		uid := snowflake.ParseInt64(template.GetUid())
		if _, ok := templates.Get(namespace); !ok {
			templates.Set(namespace, safety.NewSyncMap(map[snowflake.ID]*bo.TemplateItemBo{}))
		}
		namespaceTemplates, ok := templates.Get(namespace)
		if !ok {
			continue
		}
		createdAt, _ := time.Parse(time.DateTime, template.GetCreatedAt())
		updatedAt, _ := time.Parse(time.DateTime, template.GetUpdatedAt())
		namespaceTemplates.Set(uid, &bo.TemplateItemBo{
			UID:       uid,
			Name:      template.GetName(),
			App:       vobj.TemplateApp(template.GetApp()),
			JSONData:  template.GetJsonData(),
			Status:    vobj.GlobalStatus(template.GetStatus()),
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
	}
	return &Template{
		useDatabase:  bc.GetUseDatabase() == "true",
		templateRepo: templateRepo,
		templates:    templates,
		helper:       klog.NewHelper(klog.With(helper.Logger(), "biz", "template")),
	}
}

type Template struct {
	helper       *klog.Helper
	useDatabase  bool
	templateRepo repository.Template
	templates    *safety.SyncMap[string, *safety.SyncMap[snowflake.ID, *bo.TemplateItemBo]]
}

func (t *Template) CreateTemplate(ctx context.Context, req *bo.CreateTemplateBo) error {
	if !t.useDatabase {
		return merr.ErrorParamsNotSupportFileConfig()
	}
	doTemplate := req.ToDoTemplate()
	if _, err := t.templateRepo.GetTemplateByName(ctx, doTemplate.Name); err == nil {
		return merr.ErrorParams("template %s already exists", doTemplate.Name)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.helper.Errorw("msg", "check template exists failed", "error", err, "name", doTemplate.Name)
		return merr.ErrorInternal("create template %s failed", doTemplate.Name)
	}
	if err := t.templateRepo.CreateTemplate(ctx, doTemplate); err != nil {
		t.helper.Errorw("msg", "create template failed", "error", err, "name", doTemplate.Name)
		return merr.ErrorInternal("create template %s failed", doTemplate.Name)
	}
	return nil
}

func (t *Template) UpdateTemplate(ctx context.Context, req *bo.UpdateTemplateBo) error {
	if !t.useDatabase {
		return merr.ErrorParamsNotSupportFileConfig()
	}
	doTemplate := req.ToDoTemplate()
	existTemplate, err := t.templateRepo.GetTemplateByName(ctx, doTemplate.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		t.helper.Errorw("msg", "check template exists failed", "error", err, "name", doTemplate.Name)
		return merr.ErrorInternal("update template %s failed", doTemplate.Name)
	} else if existTemplate != nil && existTemplate.UID != doTemplate.UID {
		return merr.ErrorParams("template %s already exists", doTemplate.Name)
	}
	if err := t.templateRepo.UpdateTemplate(ctx, doTemplate); err != nil {
		t.helper.Errorw("msg", "update template failed", "error", err, "uid", doTemplate.UID)
		return merr.ErrorInternal("update template %s failed", doTemplate.UID)
	}
	return nil
}

func (t *Template) UpdateTemplateStatus(ctx context.Context, req *bo.UpdateTemplateStatusBo) error {
	if !t.useDatabase {
		return merr.ErrorParamsNotSupportFileConfig()
	}
	if err := t.templateRepo.UpdateTemplateStatus(ctx, req.UID, req.Status); err != nil {
		t.helper.Errorw("msg", "update template status failed", "error", err, "uid", req.UID)
		return merr.ErrorInternal("update template status %s failed", req.UID)
	}
	return nil
}

func (t *Template) DeleteTemplate(ctx context.Context, uid snowflake.ID) error {
	if !t.useDatabase {
		return merr.ErrorParamsNotSupportFileConfig()
	}
	if err := t.templateRepo.DeleteTemplate(ctx, uid); err != nil {
		t.helper.Errorw("msg", "delete template failed", "error", err, "uid", uid)
		return merr.ErrorInternal("delete template %s failed", uid)
	}
	return nil
}

func (t *Template) getTemplateByFileConfigWithUID(ctx context.Context, uid snowflake.ID) (*bo.TemplateItemBo, error) {
	namespaceTemplates, ok := t.templates.Get(middler.GetNamespace(ctx))
	if !ok {
		return nil, merr.ErrorNotFound("template %s not found", uid)
	}
	template, ok := namespaceTemplates.Get(uid)
	if !ok {
		return nil, merr.ErrorNotFound("template %s not found", uid)
	}
	return template, nil
}

func (t *Template) GetTemplate(ctx context.Context, uid snowflake.ID) (*bo.TemplateItemBo, error) {
	if !t.useDatabase {
		return t.getTemplateByFileConfigWithUID(ctx, uid)
	}
	doTemplate, err := t.templateRepo.GetTemplate(ctx, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, merr.ErrorNotFound("template %s not found", uid)
		}
		t.helper.Errorw("msg", "get template failed", "error", err, "uid", uid)
		return nil, merr.ErrorInternal("get template %s failed", uid)
	}
	return bo.NewTemplateItemBo(doTemplate), nil
}

func (t *Template) getTemplateByFileConfigWithNamespace(ctx context.Context, req *bo.ListTemplateBo) ([]*bo.TemplateItemBo, error) {
	namespace := middler.GetNamespace(ctx)
	namespaceTemplates, ok := t.templates.Get(namespace)
	if !ok {
		return nil, merr.ErrorNotFound("template %s not found", namespace)
	}
	templates := make([]*bo.TemplateItemBo, 0, namespaceTemplates.Len())
	for _, template := range namespaceTemplates.Values() {
		if strutil.IsNotEmpty(req.Keyword) && !strings.Contains(template.Name, req.Keyword) {
			continue
		}
		if req.Status.Exist() && !req.Status.IsUnknown() && template.Status != req.Status {
			continue
		}
		if req.App.Exist() && !req.App.IsUnknown() && template.App != req.App {
			continue
		}
		templates = append(templates, template)
	}
	total := int64(len(templates))
	pageRequestBo := bo.NewPageRequestBo(req.Page, req.PageSize)
	pageRequestBo.WithTotal(total)
	req.PageRequestBo = pageRequestBo
	return templates, nil
}

func (t *Template) ListTemplate(ctx context.Context, req *bo.ListTemplateBo) (*bo.PageResponseBo[*bo.TemplateItemBo], error) {
	if !t.useDatabase {
		templates, err := t.getTemplateByFileConfigWithNamespace(ctx, req)
		if err != nil {
			return nil, err
		}
		return bo.NewPageResponseBo(req.PageRequestBo, templates), nil
	}
	pageResponseBo, err := t.templateRepo.ListTemplate(ctx, req)
	if err != nil {
		t.helper.Errorw("msg", "list template failed", "error", err, "req", req)
		return nil, merr.ErrorInternal("list template failed")
	}
	items := make([]*bo.TemplateItemBo, 0, len(pageResponseBo.GetItems()))
	for _, item := range pageResponseBo.GetItems() {
		items = append(items, bo.NewTemplateItemBo(item))
	}
	return bo.NewPageResponseBo(pageResponseBo.PageRequestBo, items), nil
}

func (t *Template) getSelectTemplateByFileConfig(ctx context.Context, req *bo.SelectTemplateBo) ([]*bo.TemplateItemSelectBo, error) {
	namespace := middler.GetNamespace(ctx)
	namespaceTemplates, ok := t.templates.Get(namespace)
	if !ok {
		return nil, merr.ErrorNotFound("template %s not found", namespace)
	}
	templates := make([]*bo.TemplateItemSelectBo, 0, namespaceTemplates.Len())
	for _, template := range namespaceTemplates.Values() {
		if strutil.IsNotEmpty(req.Keyword) && !strings.Contains(template.Name, req.Keyword) {
			continue
		}
		if req.Status.Exist() && !req.Status.IsUnknown() && template.Status != req.Status {
			continue
		}
		if req.App.Exist() && !req.App.IsUnknown() && template.App != req.App {
			continue
		}
		templates = append(templates, &bo.TemplateItemSelectBo{
			UID:      template.UID,
			Name:     template.Name,
			Status:   template.Status,
			Disabled: template.Status != vobj.GlobalStatusEnabled,
			Tooltip:  "",
		})
	}
	total := int64(len(templates))
	req.Limit = int32(total)
	req.LastUID = 0
	return templates, nil
}

func (t *Template) SelectTemplate(ctx context.Context, req *bo.SelectTemplateBo) (*bo.SelectTemplateBoResult, error) {
	if !t.useDatabase {
		templates, err := t.getSelectTemplateByFileConfig(ctx, req)
		if err != nil {
			return nil, err
		}
		return &bo.SelectTemplateBoResult{
			Items:   templates,
			Total:   int64(len(templates)),
			LastUID: 0,
		}, nil
	}
	result, err := t.templateRepo.SelectTemplate(ctx, req)
	if err != nil {
		t.helper.Errorw("msg", "select template failed", "error", err, "req", req)
		return nil, merr.ErrorInternal("select template failed")
	}
	items := make([]*bo.TemplateItemSelectBo, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, bo.NewTemplateItemSelectBo(item))
	}
	return &bo.SelectTemplateBoResult{
		Items:   items,
		Total:   result.Total,
		LastUID: result.LastUID,
	}, nil
}
