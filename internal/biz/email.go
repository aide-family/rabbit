package biz

import (
	"context"

	"github.com/aide-family/rabbit/internal/biz/bo"
)

func NewEmail() *Email {
	return &Email{}
}

type Email struct{}

func (e *Email) SendEmail(ctx context.Context, req *bo.SendEmailBo) error {
	return nil
}
