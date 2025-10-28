// Package bo is the business logic object
package bo

import (
	"time"

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
		UID:      c.UID,
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
		Headers:     c.Headers,
	}
}

func NewCreateEmailTemplateBo(req *apiv1.CreateEmailTemplateRequest) *CreateEmailTemplateBo {
	return &CreateEmailTemplateBo{
		Name:        req.Name,
		Subject:     req.Subject,
		Body:        req.Body,
		ContentType: req.ContentType,
		Headers:     req.Headers,
	}
}

type UpdateEmailTemplateBo struct {
	UID string
	CreateEmailTemplateBo
}

func (c *UpdateEmailTemplateBo) ToDoEmailTemplate() *do.EmailTemplate {
	return &do.EmailTemplate{
		UID:         c.UID,
		Name:        c.Name,
		Subject:     c.Subject,
		Body:        c.Body,
		ContentType: c.ContentType,
		Headers:     c.Headers,
	}
}

func NewUpdateEmailTemplateBo(req *apiv1.UpdateEmailTemplateRequest) *UpdateEmailTemplateBo {
	return &UpdateEmailTemplateBo{
		UID: req.Uid,
		CreateEmailTemplateBo: CreateEmailTemplateBo{
			Name:        req.Name,
			Subject:     req.Subject,
			Body:        req.Body,
			ContentType: req.ContentType,
			Headers:     req.Headers,
		},
	}
}

type UpdateEmailTemplateStatusBo struct {
	UID    string
	Status vobj.GlobalStatus
}

func NewUpdateEmailTemplateStatusBo(req *apiv1.UpdateEmailTemplateStatusRequest) *UpdateEmailTemplateStatusBo {
	return &UpdateEmailTemplateStatusBo{
		UID:    req.Uid,
		Status: vobj.GlobalStatus(req.Status),
	}
}

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
		Headers:     doEmailTemplate.Headers,
		Status:      doEmailTemplate.Status,
		CreatedAt:   doEmailTemplate.CreatedAt,
		UpdatedAt:   doEmailTemplate.UpdatedAt,
	}
}

func (b *EmailTemplateItemBo) ToAPIV1EmailTemplateItem() *apiv1.EmailTemplateItem {
	return &apiv1.EmailTemplateItem{
		Uid:         b.UID,
		Subject:     b.Subject,
		Body:        b.Body,
		ContentType: b.ContentType,
		Headers:     b.Headers,
		Status:      enum.GlobalStatus(b.Status),
		CreatedAt:   b.CreatedAt.Format(time.DateTime),
		UpdatedAt:   b.UpdatedAt.Format(time.DateTime),
	}
}

type ListEmailTemplateBo struct {
	*PageRequestBo
	Keyword string
	Status  vobj.GlobalStatus
}

func NewListEmailTemplateBo(req *apiv1.ListEmailTemplateRequest) *ListEmailTemplateBo {
	return &ListEmailTemplateBo{
		PageRequestBo: NewPageRequestBo(req.Page, req.PageSize),
		Keyword:       req.Keyword,
		Status:        vobj.GlobalStatus(req.Status),
	}
}

func ToAPIV1ListEmailTemplateReply(pageResponseBo *PageResponseBo[*EmailTemplateItemBo]) *apiv1.ListEmailTemplateReply {
	items := make([]*apiv1.EmailTemplateItem, 0, len(pageResponseBo.GetItems()))
	for _, item := range pageResponseBo.GetItems() {
		items = append(items, item.ToAPIV1EmailTemplateItem())
	}
	return &apiv1.ListEmailTemplateReply{
		Items:    items,
		Total:    pageResponseBo.GetTotal(),
		Page:     pageResponseBo.GetPage(),
		PageSize: pageResponseBo.GetPageSize(),
	}
}
