// Package bo is the business logic object
package bo

import (
	"context"
	"time"

	"github.com/aide-family/magicbox/serialize"
	"github.com/aide-family/magicbox/strutil"
	"github.com/bwmarrin/snowflake"

	"github.com/aide-family/rabbit/internal/biz/do"
	"github.com/aide-family/rabbit/internal/biz/vobj"
	apiv1 "github.com/aide-family/rabbit/pkg/api/v1"
	"github.com/aide-family/rabbit/pkg/enum"
	"github.com/aide-family/rabbit/pkg/middler"
)

type SendEmailBo struct {
	Namespace   string
	Subject     string
	Body        string
	To          []string
	Cc          []string
	ContentType string
	Headers     map[string]string
}

func (b *SendEmailBo) ToMessageLog() (*do.MessageLog, error) {
	messageBytes, err := serialize.JSONMarshal(b)
	if err != nil {
		return nil, err
	}
	return &do.MessageLog{
		SendAt:  time.Now(),
		Message: string(messageBytes),
		Type:    vobj.MessageTypeEmail,
		Status:  vobj.MessageStatusPending,
	}, nil
}

func NewSendEmailBo(ctx context.Context, req *apiv1.SendEmailRequest) *SendEmailBo {
	return &SendEmailBo{
		Namespace:   middler.GetNamespace(ctx),
		Subject:     req.Subject,
		Body:        req.Body,
		To:          req.To,
		Cc:          req.Cc,
		ContentType: req.ContentType,
		Headers:     req.Headers,
	}
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
	UID snowflake.ID
	CreateEmailConfigBo
}

func (c *UpdateEmailConfigBo) ToDoEmailConfig() *do.EmailConfig {
	emailConfig := &do.EmailConfig{
		Name:     c.Name,
		Host:     c.Host,
		Port:     c.Port,
		Username: c.Username,
		Password: strutil.EncryptString(c.Password),
	}
	emailConfig.WithUID(c.UID)
	return emailConfig
}

func NewUpdateEmailConfigBo(req *apiv1.UpdateEmailConfigRequest) *UpdateEmailConfigBo {
	return &UpdateEmailConfigBo{
		UID: snowflake.ParseInt64(req.Uid),
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
	UID    snowflake.ID
	Status vobj.GlobalStatus
}

func NewUpdateEmailConfigStatusBo(req *apiv1.UpdateEmailConfigStatusRequest) *UpdateEmailConfigStatusBo {
	return &UpdateEmailConfigStatusBo{
		UID:    snowflake.ParseInt64(req.Uid),
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
	UID       snowflake.ID
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
		Uid:       b.UID.Int64(),
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
