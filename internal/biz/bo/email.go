// Package bo is the business logic object
package bo

import (
	"net/http"
	"time"

	"github.com/aide-family/magicbox/message/email"
	"github.com/aide-family/magicbox/serialize"
	"github.com/aide-family/magicbox/strutil"
	"github.com/bwmarrin/snowflake"

	"github.com/aide-family/rabbit/internal/biz/do"
	"github.com/aide-family/rabbit/internal/biz/vobj"
	apiv1 "github.com/aide-family/rabbit/pkg/api/v1"
	"github.com/aide-family/rabbit/pkg/enum"
)

var _ email.Config = (*EmailConfigItemBo)(nil)

type SendEmailBo struct {
	UID         snowflake.ID `json:"uid"`
	Subject     string       `json:"subject"`
	Body        string       `json:"body"`
	To          []string     `json:"to"`
	Cc          []string     `json:"cc"`
	ContentType string       `json:"content_type"`
	Headers     http.Header  `json:"headers"`
}

func (b *SendEmailBo) ToMessageLog(emailConfig *do.EmailConfig) (*do.MessageLog, error) {
	messageBytes, err := serialize.JSONMarshal(b)
	if err != nil {
		return nil, err
	}
	emailConfigBytes, err := serialize.JSONMarshal(NewEmailConfigItemBo(emailConfig))
	if err != nil {
		return nil, err
	}
	return &do.MessageLog{
		SendAt:  time.Now(),
		Message: strutil.EncryptString(messageBytes),
		Config:  strutil.EncryptString(emailConfigBytes),
		Type:    vobj.MessageTypeEmail,
		Status:  vobj.MessageStatusPending,
	}, nil
}

func NewSendEmailBo(req *apiv1.SendEmailRequest) *SendEmailBo {
	headers := make(http.Header)
	for key, value := range req.Headers {
		headers.Add(key, value)
	}
	return &SendEmailBo{
		UID:         snowflake.ParseInt64(req.Uid),
		Subject:     req.Subject,
		Body:        req.Body,
		To:          req.To,
		Cc:          req.Cc,
		ContentType: req.ContentType,
		Headers:     headers,
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
	UID       snowflake.ID      `json:"uid"`
	Name      string            `json:"name"`
	Host      string            `json:"host"`
	Port      int32             `json:"port"`
	Username  string            `json:"username"`
	Password  string            `json:"password"`
	Status    vobj.GlobalStatus `json:"status"`
	CreatedAt time.Time         `json:"-"`
	UpdatedAt time.Time         `json:"-"`
}

// GetHost implements email.Config.
func (b *EmailConfigItemBo) GetHost() string {
	return b.Host
}

// GetPassword implements email.Config.
func (b *EmailConfigItemBo) GetPassword() string {
	return b.Password
}

// GetPort implements email.Config.
func (b *EmailConfigItemBo) GetPort() int32 {
	return b.Port
}

// GetUsername implements email.Config.
func (b *EmailConfigItemBo) GetUsername() string {
	return b.Username
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
