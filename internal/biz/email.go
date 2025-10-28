package biz

import (
	"context"

	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/repository"
)

func NewEmail(messageLogRepository repository.MessageLog) *Email {
	return &Email{
		messageLogRepository: messageLogRepository,
	}
}

type Email struct {
	messageLogRepository repository.MessageLog
}

func (e *Email) SendEmail(ctx context.Context, req *bo.SendEmailBo) error {
	return nil
}
