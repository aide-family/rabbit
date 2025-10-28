// Package bo is the business logic object
package bo

import (
	"time"

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
	*PageRequestBo[*EmailConfigItemBo]
	Keyword string
	Status  vobj.GlobalStatus
}

func NewListEmailConfigBo(req *apiv1.ListEmailConfigRequest) *ListEmailConfigBo {
	return &ListEmailConfigBo{
		PageRequestBo: NewPageRequestBo[*EmailConfigItemBo](req.Page, req.PageSize),
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

type ListEmailTemplateBo struct {
	*PageRequestBo[EmailTemplateItemBo]
	Keyword string
	Status  vobj.GlobalStatus
}
