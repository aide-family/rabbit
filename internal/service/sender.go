package service

import (
	"context"

	"github.com/aide-family/rabbit/internal/biz"
	"github.com/aide-family/rabbit/internal/biz/bo"
	apiv1 "github.com/aide-family/rabbit/pkg/api/v1"
)

func NewSenderService(emailBiz *biz.Email) *SenderService {
	return &SenderService{
		emailBiz: emailBiz,
	}
}

type SenderService struct {
	apiv1.UnimplementedSenderServer

	emailBiz *biz.Email
}

func (s *SenderService) SendEmail(ctx context.Context, req *apiv1.SendEmailRequest) (*apiv1.SendReply, error) {
	sendEmailBo := bo.NewSendEmailBo(req)
	if err := s.emailBiz.AppendEmailMessage(ctx, sendEmailBo); err != nil {
		return nil, err
	}
	return &apiv1.SendReply{}, nil
}

func (s *SenderService) SendEmailWithTemplate(ctx context.Context, req *apiv1.SendEmailWithTemplateRequest) (*apiv1.SendReply, error) {
	return &apiv1.SendReply{}, nil
}

func (s *SenderService) SendWebhook(ctx context.Context, req *apiv1.SendWebhookRequest) (*apiv1.SendReply, error) {
	return &apiv1.SendReply{}, nil
}

func (s *SenderService) SendWebhookWithTemplate(ctx context.Context, req *apiv1.SendWebhookWithTemplateRequest) (*apiv1.SendReply, error) {
	return &apiv1.SendReply{}, nil
}
