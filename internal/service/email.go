package service

import (
	"context"

	"github.com/aide-family/rabbit/internal/biz"
	"github.com/aide-family/rabbit/internal/biz/bo"
	apiv1 "github.com/aide-family/rabbit/pkg/api/v1"
)

type EmailService struct {
	apiv1.UnimplementedEmailServer

	emailBiz *biz.EmailBiz
}

func NewEmailService(emailBiz *biz.EmailBiz) *EmailService {
	return &EmailService{
		emailBiz: emailBiz,
	}
}

func (s *EmailService) SendEmail(ctx context.Context, req *apiv1.SendEmailRequest) (*apiv1.SendEmailReply, error) {
	params := &bo.SendEmailBo{
		Namespace:   req.Namespace,
		Subject:     req.Subject,
		Body:        req.Body,
		To:          req.To,
		Cc:          req.Cc,
		ContentType: req.ContentType,
		Headers:     req.Headers,
	}

	if err := s.emailBiz.SendEmail(ctx, params); err != nil {
		return nil, err
	}
	return &apiv1.SendEmailReply{
		Message: "Email sent successfully",
	}, nil
}
