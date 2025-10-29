package biz

import (
	"context"

	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/repository"
)

func NewEmail(messageLogRepo repository.MessageLog) *Email {
	return &Email{
		messageLogRepo: messageLogRepo,
	}
}

type Email struct {
	messageLogRepo repository.MessageLog
}

func (e *Email) SendEmail(ctx context.Context, req *bo.SendEmailBo) error {
	return nil
}
