// Package fileimpl is the implementation of the template repository for file config
package fileimpl

import (
	"context"
	"encoding/json"
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

func NewTemplateRepository(d *data.Data) repository.Template {
	t := &templateRepositoryImpl{
		d:         d,
		templates: d.GetFileConfig().GetTemplates(),
	}
	t.initTemplates()
	d.RegisterReloadFunc(data.KeyTemplates, func() {
		t.initTemplates()
	})
	return t
}

type templateRepositoryImpl struct {
	d                 *data.Data
	templates         []*conf.Config_Template
	templatesWithUID  *safety.SyncMap[string, *safety.SyncMap[snowflake.ID, *do.Template]]
	templatesWithName *safety.SyncMap[string, *safety.SyncMap[string, *do.Template]]
}

func (t *templateRepositoryImpl) initTemplates() {
	t.templates = t.d.GetFileConfig().GetTemplates()
	t.templatesWithUID = safety.NewSyncMap(make(map[string]*safety.SyncMap[snowflake.ID, *do.Template]))
	t.templatesWithName = safety.NewSyncMap(make(map[string]*safety.SyncMap[string, *do.Template]))
	for _, template := range t.templates {
		namespace := template.GetNamespace()
		uid := snowflake.ParseInt64(template.GetUid())
		name := template.GetName()
		if _, ok := t.templatesWithUID.Get(namespace); !ok {
			t.templatesWithUID.Set(namespace, safety.NewSyncMap(map[snowflake.ID]*do.Template{}))
			t.templatesWithName.Set(namespace, safety.NewSyncMap(map[string]*do.Template{}))
		}
		item := t.toDoTemplate(template)
		if namespaceTemplatesByName, ok := t.templatesWithName.Get(namespace); ok {
			namespaceTemplatesByName.Set(name, item)
		}
		if namespaceTemplatesByUID, ok := t.templatesWithUID.Get(namespace); ok {
			namespaceTemplatesByUID.Set(uid, item)
		}
	}
}

func (t *templateRepositoryImpl) toDoTemplate(template *conf.Config_Template) *do.Template {
	createdAt, _ := time.Parse(time.DateTime, template.GetCreatedAt())
	updatedAt, _ := time.Parse(time.DateTime, template.GetUpdatedAt())
	jsonData := json.RawMessage(template.GetJsonData())
	return &do.Template{
		NamespaceModel: do.NamespaceModel{
			Namespace: template.GetNamespace(),
			BaseModel: do.BaseModel{
				ID:        template.GetId(),
				UID:       snowflake.ParseInt64(template.GetUid()),
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
		},
		Name:     template.GetName(),
		App:      vobj.TemplateApp(template.GetApp()),
		JSONData: jsonData,
		Status:   vobj.GlobalStatus(template.GetStatus()),
	}
}

// CreateTemplate implements repository.Template.
func (t *templateRepositoryImpl) CreateTemplate(ctx context.Context, req *do.Template) error {
	return merr.ErrorParamsNotSupportFileConfig()
}

// UpdateTemplate implements repository.Template.
func (t *templateRepositoryImpl) UpdateTemplate(ctx context.Context, req *do.Template) error {
	return merr.ErrorParamsNotSupportFileConfig()
}

// UpdateTemplateStatus implements repository.Template.
func (t *templateRepositoryImpl) UpdateTemplateStatus(ctx context.Context, uid snowflake.ID, status vobj.GlobalStatus) error {
	return merr.ErrorParamsNotSupportFileConfig()
}

// DeleteTemplate implements repository.Template.
func (t *templateRepositoryImpl) DeleteTemplate(ctx context.Context, uid snowflake.ID) error {
	return merr.ErrorParamsNotSupportFileConfig()
}

// GetTemplate implements repository.Template.
func (t *templateRepositoryImpl) GetTemplate(ctx context.Context, uid snowflake.ID) (*do.Template, error) {
	namespace := middler.GetNamespace(ctx)
	templateWithUID, ok := t.templatesWithUID.Get(namespace)
	if !ok {
		return nil, merr.ErrorNotFound("template not found")
	}
	template, ok := templateWithUID.Get(uid)
	if !ok {
		return nil, merr.ErrorNotFound("template not found")
	}
	return template, nil
}

// GetTemplateByName implements repository.Template.
func (t *templateRepositoryImpl) GetTemplateByName(ctx context.Context, name string) (*do.Template, error) {
	namespace := middler.GetNamespace(ctx)
	templateWithName, ok := t.templatesWithName.Get(namespace)
	if !ok {
		return nil, merr.ErrorNotFound("template not found")
	}
	template, ok := templateWithName.Get(name)
	if !ok {
		return nil, merr.ErrorNotFound("template not found")
	}
	return template, nil
}

// ListTemplate implements repository.Template.
func (t *templateRepositoryImpl) ListTemplate(ctx context.Context, req *bo.ListTemplateBo) (*bo.PageResponseBo[*do.Template], error) {
	namespace := middler.GetNamespace(ctx)
	templateWithUID, ok := t.templatesWithUID.Get(namespace)
	if !ok {
		pageRequestBo := bo.NewPageRequestBo(req.Page, req.PageSize)
		pageRequestBo.WithTotal(0)
		req.PageRequestBo = pageRequestBo
		return bo.NewPageResponseBo(req.PageRequestBo, []*do.Template{}), nil
	}
	templates := make([]*do.Template, 0, templateWithUID.Len())
	for _, template := range templateWithUID.Values() {
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
	sort.Slice(templates, func(i, j int) bool {
		return templates[i].CreatedAt.After(templates[j].CreatedAt)
	})
	return bo.NewPageResponseBo(req.PageRequestBo, templates), nil
}

// SelectTemplate implements repository.Template.
func (t *templateRepositoryImpl) SelectTemplate(ctx context.Context, req *bo.SelectTemplateBo) (*bo.SelectTemplateResult, error) {
	namespace := middler.GetNamespace(ctx)
	templateWithUID, ok := t.templatesWithUID.Get(namespace)
	if !ok {
		return &bo.SelectTemplateResult{
			Items:   []*do.Template{},
			Total:   0,
			LastUID: 0,
		}, nil
	}
	templates := make([]*do.Template, 0, templateWithUID.Len())
	for _, template := range templateWithUID.Values() {
		if strutil.IsNotEmpty(req.Keyword) && !strings.Contains(template.Name, req.Keyword) {
			continue
		}
		if req.Status.Exist() && !req.Status.IsUnknown() && template.Status != req.Status {
			continue
		}
		if req.App.Exist() && !req.App.IsUnknown() && template.App != req.App {
			continue
		}
		if req.LastUID > 0 && template.UID >= req.LastUID {
			continue
		}
		templates = append(templates, template)
	}
	total := int64(len(templates))
	sort.Slice(templates, func(i, j int) bool {
		return templates[i].UID > templates[j].UID
	})
	if int32(len(templates)) > req.Limit {
		templates = templates[:req.Limit]
	}
	var lastUID snowflake.ID
	if len(templates) > 0 {
		lastUID = templates[len(templates)-1].UID
	}
	return &bo.SelectTemplateResult{
		Items:   templates,
		Total:   total,
		LastUID: lastUID,
	}, nil
}
