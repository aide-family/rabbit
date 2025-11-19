package biz

import (
	"context"
	"errors"

	"github.com/bwmarrin/snowflake"
	klog "github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"

	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/pkg/merr"
)

func NewTemplate(
	templateRepo repository.Template,
	helper *klog.Helper,
) *Template {
	return &Template{
		templateRepo: templateRepo,
		helper:       klog.NewHelper(klog.With(helper.Logger(), "biz", "template")),
	}
}

type Template struct {
	helper       *klog.Helper
	templateRepo repository.Template
}

func (t *Template) CreateTemplate(ctx context.Context, req *bo.CreateTemplateBo) error {
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
	if err := t.templateRepo.UpdateTemplateStatus(ctx, req.UID, req.Status); err != nil {
		t.helper.Errorw("msg", "update template status failed", "error", err, "uid", req.UID)
		return merr.ErrorInternal("update template status %s failed", req.UID)
	}
	return nil
}

func (t *Template) DeleteTemplate(ctx context.Context, uid snowflake.ID) error {
	if err := t.templateRepo.DeleteTemplate(ctx, uid); err != nil {
		t.helper.Errorw("msg", "delete template failed", "error", err, "uid", uid)
		return merr.ErrorInternal("delete template %s failed", uid)
	}
	return nil
}

func (t *Template) GetTemplate(ctx context.Context, uid snowflake.ID) (*bo.TemplateItemBo, error) {
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

func (t *Template) ListTemplate(ctx context.Context, req *bo.ListTemplateBo) (*bo.PageResponseBo[*bo.TemplateItemBo], error) {
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

func (t *Template) SelectTemplate(ctx context.Context, req *bo.SelectTemplateBo) (*bo.SelectTemplateBoResult, error) {
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
