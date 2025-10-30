// Package bo is the business logic object
package bo

import (
	"time"

	"github.com/aide-family/magicbox/safety"
	"github.com/aide-family/magicbox/strutil"

	"github.com/aide-family/rabbit/internal/biz/do"
	"github.com/aide-family/rabbit/internal/biz/vobj"
	apiv1 "github.com/aide-family/rabbit/pkg/api/v1"
	"github.com/aide-family/rabbit/pkg/enum"
)

type SendEmailBo struct {
	Namespace   string
	Subject     string
	Body        string
	To          []string
	Cc          []string
	ContentType string
	Headers     []string
}

type CreateEmailConfigBo struct {
	Name     string
	Host     string
	Port     int32
	Username string
	Password string
}

func (c *CreateEmailConfigBo) ToDoEmailConfig() *do.EmailConfig {
	return &do.EmailConfig{
		Name:     c.Name,
		Host:     c.Host,
		Port:     c.Port,
		Username: c.Username,
		Password: strutil.EncryptString(c.Password),
	}
}

func NewCreateEmailConfigBo(req *apiv1.CreateEmailConfigRequest) *CreateEmailConfigBo {
	return &CreateEmailConfigBo{
		Name:     req.Name,
		Host:     req.Host,
		Port:     req.Port,
		Username: req.Username,
		Password: req.Password,
	}
}

type UpdateEmailConfigBo struct {
	UID string
	CreateEmailConfigBo
}

func (c *UpdateEmailConfigBo) ToDoEmailConfig() *do.EmailConfig {
	return &do.EmailConfig{
		NamespaceModel: do.NamespaceModel{
			UID: c.UID,
		},
		Name:     c.Name,
		Host:     c.Host,
		Port:     c.Port,
		Username: c.Username,
		Password: strutil.EncryptString(c.Password),
	}
}

func NewUpdateEmailConfigBo(req *apiv1.UpdateEmailConfigRequest) *UpdateEmailConfigBo {
	return &UpdateEmailConfigBo{
		UID: req.Uid,
		CreateEmailConfigBo: CreateEmailConfigBo{
			Name:     req.Name,
			Host:     req.Host,
			Port:     req.Port,
			Username: req.Username,
			Password: req.Password,
		},
	}
}

type UpdateEmailConfigStatusBo struct {
	UID    string
	Status vobj.GlobalStatus
}

func NewUpdateEmailConfigStatusBo(req *apiv1.UpdateEmailConfigStatusRequest) *UpdateEmailConfigStatusBo {
	return &UpdateEmailConfigStatusBo{
		UID:    req.Uid,
		Status: vobj.GlobalStatus(req.Status),
	}
}

type ListEmailConfigBo struct {
	*PageRequestBo
	Keyword string
	Status  vobj.GlobalStatus
}

func NewListEmailConfigBo(req *apiv1.ListEmailConfigRequest) *ListEmailConfigBo {
	return &ListEmailConfigBo{
		PageRequestBo: NewPageRequestBo(req.Page, req.PageSize),
		Keyword:       req.Keyword,
		Status:        vobj.GlobalStatus(req.Status),
	}
}

func ToAPIV1ListEmailConfigReply(pageResponseBo *PageResponseBo[*EmailConfigItemBo]) *apiv1.ListEmailConfigReply {
	items := make([]*apiv1.EmailConfigItem, 0, len(pageResponseBo.GetItems()))
	for _, item := range pageResponseBo.GetItems() {
		items = append(items, item.ToAPIV1EmailConfigItem())
	}
	return &apiv1.ListEmailConfigReply{
		Items:    items,
		Total:    pageResponseBo.GetTotal(),
		Page:     pageResponseBo.GetPage(),
		PageSize: pageResponseBo.GetPageSize(),
	}
}

type EmailConfigItemBo struct {
	UID       string
	Name      string
	Host      string
	Port      int32
	Username  string
	Password  string
	Status    vobj.GlobalStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewEmailConfigItemBo(doEmailConfig *do.EmailConfig) *EmailConfigItemBo {
	return &EmailConfigItemBo{
		UID:       doEmailConfig.UID,
		Name:      doEmailConfig.Name,
		Host:      doEmailConfig.Host,
		Port:      doEmailConfig.Port,
		Username:  doEmailConfig.Username,
		Password:  string(doEmailConfig.Password),
		Status:    doEmailConfig.Status,
		CreatedAt: doEmailConfig.CreatedAt,
		UpdatedAt: doEmailConfig.UpdatedAt,
	}
}

func (b *EmailConfigItemBo) ToAPIV1EmailConfigItem() *apiv1.EmailConfigItem {
	return &apiv1.EmailConfigItem{
		Uid:       b.UID,
		Name:      b.Name,
		Host:      b.Host,
		Port:      b.Port,
		Username:  b.Username,
		Password:  b.Password,
		Status:    enum.GlobalStatus(b.Status),
		CreatedAt: b.CreatedAt.Format(time.DateTime),
		UpdatedAt: b.UpdatedAt.Format(time.DateTime),
	}
}

type CreateEmailTemplateBo struct {
	Name        string
	Subject     string
	Body        string
	ContentType string
	Headers     map[string]string
}

func (c *CreateEmailTemplateBo) ToDoEmailTemplate() *do.EmailTemplate {
	return &do.EmailTemplate{
		Name:        c.Name,
		Subject:     c.Subject,
		Body:        c.Body,
		ContentType: c.ContentType,
		Headers:     safety.NewMap(c.Headers),
	}
}

// NewCreateEmailTemplateBo 已被移除，请使用 template 统一接口
// 使用 bo.CreateEmailTemplateBoFromEmail 将旧的 BO 转换为新的统一 Template BO

type UpdateEmailTemplateBo struct {
	UID string
	CreateEmailTemplateBo
}

func (c *UpdateEmailTemplateBo) ToDoEmailTemplate() *do.EmailTemplate {
	return &do.EmailTemplate{
		NamespaceModel: do.NamespaceModel{
			UID: c.UID,
		},
		Name:        c.Name,
		Subject:     c.Subject,
		Body:        c.Body,
		ContentType: c.ContentType,
		Headers:     safety.NewMap(c.Headers),
	}
}

// NewUpdateEmailTemplateBo 已被移除，请使用 template 统一接口

type UpdateEmailTemplateStatusBo struct {
	UID    string
	Status vobj.GlobalStatus
}

// NewUpdateEmailTemplateStatusBo 已被移除，请使用 template 统一接口

type EmailTemplateItemBo struct {
	UID         string
	Subject     string
	Body        string
	ContentType string
	Headers     map[string]string
	Status      vobj.GlobalStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewEmailTemplateItemBo(doEmailTemplate *do.EmailTemplate) *EmailTemplateItemBo {
	return &EmailTemplateItemBo{
		UID:         doEmailTemplate.UID,
		Subject:     doEmailTemplate.Subject,
		Body:        doEmailTemplate.Body,
		ContentType: doEmailTemplate.ContentType,
		Headers:     doEmailTemplate.Headers.Map(),
		Status:      doEmailTemplate.Status,
		CreatedAt:   doEmailTemplate.CreatedAt,
		UpdatedAt:   doEmailTemplate.UpdatedAt,
	}
}

// ToAPIV1EmailTemplateItem 已被移除，请使用 template 统一接口

type ListEmailTemplateBo struct {
	*PageRequestBo
	Keyword string
	Status  vobj.GlobalStatus
}

// NewListEmailTemplateBo 已被移除，请使用 template 统一接口

// ToAPIV1ListEmailTemplateReply 已被移除，请使用 template 统一接口
