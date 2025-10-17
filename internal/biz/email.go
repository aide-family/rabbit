package biz

import (
	"context"

	"github.com/aide-family/rabbit/internal/biz/bo"
)

type EmailBiz struct{}

func NewEmailBiz() *EmailBiz {
	return &EmailBiz{}
}

func (b *EmailBiz) SendEmail(ctx context.Context, req *bo.SendEmailBo) error {
	return nil
}
