package service

import (
	"context"

	apiv1 "github.com/aide-family/rabbit/pkg/api/v1"
)

type EmailService struct {
	apiv1.UnimplementedEmailServer
}

func NewEmailService() *EmailService {
	return &EmailService{}
}

func (s *EmailService) SendEmail(ctx context.Context, req *apiv1.SendEmailRequest) (*apiv1.SendEmailReply, error) {
	return &apiv1.SendEmailReply{
		Message: "Email sent successfully",
	}, nil
}
