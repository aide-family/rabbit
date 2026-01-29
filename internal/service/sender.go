package service

import (
	"context"

	"github.com/aide-family/rabbit/internal/biz"
	"github.com/aide-family/rabbit/internal/biz/bo"
	apiv1 "github.com/aide-family/rabbit/pkg/api/v1"
	"github.com/bwmarrin/snowflake"
)

func NewSenderService(emailBiz *biz.Email, webhookBiz *biz.Webhook, messageBiz *biz.Message) *SenderService {
	return &SenderService{
		emailBiz:   emailBiz,
		webhookBiz: webhookBiz,
		messageBiz: messageBiz,
	}
}

type SenderService struct {
	apiv1.UnimplementedSenderServer

	emailBiz   *biz.Email
	webhookBiz *biz.Webhook
	messageBiz *biz.Message
}

func (s *SenderService) SendMessage(ctx context.Context, req *apiv1.SendMessageRequest) (*apiv1.SendReply, error) {
	if err := s.messageBiz.SendMessage(ctx, snowflake.ParseInt64(req.Uid)); err != nil {
		return nil, err
	}
	return &apiv1.SendReply{}, nil
}

func (s *SenderService) SendEmail(ctx context.Context, req *apiv1.SendEmailRequest) (*apiv1.SendReply, error) {
	sendEmailBo := bo.NewSendEmailBo(req)
	if err := s.emailBiz.AppendEmailMessage(ctx, sendEmailBo); err != nil {
		return nil, err
	}
	return &apiv1.SendReply{}, nil
}

func (s *SenderService) SendEmailWithTemplate(ctx context.Context, req *apiv1.SendEmailWithTemplateRequest) (*apiv1.SendReply, error) {
	sendEmailWithTemplateBo, err := bo.NewSendEmailWithTemplateBo(req)
	if err != nil {
		return nil, err
	}
	if err := s.emailBiz.AppendEmailMessageWithTemplate(ctx, sendEmailWithTemplateBo); err != nil {
		return nil, err
	}
	return &apiv1.SendReply{}, nil
}

func (s *SenderService) SendWebhook(ctx context.Context, req *apiv1.SendWebhookRequest) (*apiv1.SendReply, error) {
	sendWebhookBo := bo.NewSendWebhookBo(req)
	if err := s.webhookBiz.AppendWebhookMessage(ctx, sendWebhookBo); err != nil {
		return nil, err
	}
	return &apiv1.SendReply{}, nil
}

func (s *SenderService) SendWebhookWithTemplate(ctx context.Context, req *apiv1.SendWebhookWithTemplateRequest) (*apiv1.SendReply, error) {
	sendWebhookWithTemplateBo, err := bo.NewSendWebhookWithTemplateBo(req)
	if err != nil {
		return nil, err
	}
	if err := s.webhookBiz.AppendWebhookMessageWithTemplate(ctx, sendWebhookWithTemplateBo); err != nil {
		return nil, err
	}
	return &apiv1.SendReply{}, nil
}
