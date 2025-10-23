package service

import (
	"context"

	"github.com/aide-family/magicbox/merr"
	"github.com/aide-family/magicbox/strutil"
	"github.com/aide-family/rabbit/internal/biz"
	apiv1 "github.com/aide-family/rabbit/pkg/api/v1"
	"github.com/aide-family/rabbit/pkg/middler"
)

type SenderService struct {
	apiv1.UnimplementedSenderServer

	emailBiz *biz.Email
}

func NewSenderService(emailBiz *biz.Email) *SenderService {
	return &SenderService{
		emailBiz: emailBiz,
	}
}

func (s *SenderService) SendEmail(ctx context.Context, req *apiv1.SendEmailRequest) (*apiv1.SendReply, error) {
	namespace := middler.GetNamespace(ctx)
	if strutil.IsEmpty(namespace) {
		return nil, merr.ErrorBadRequest("namespace is required")
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
