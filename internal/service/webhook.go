package service

import (
	"context"

	apiv1 "github.com/aide-family/rabbit/pkg/api/v1"
)

type WebhookService struct {
	apiv1.UnimplementedWebhookServer
}

func NewWebhookService() *WebhookService {
	return &WebhookService{}
}

func (s *WebhookService) CreateWebhook(ctx context.Context, req *apiv1.CreateWebhookRequest) (*apiv1.CreateWebhookReply, error) {
	return &apiv1.CreateWebhookReply{}, nil
}

func (s *WebhookService) UpdateWebhook(ctx context.Context, req *apiv1.UpdateWebhookRequest) (*apiv1.UpdateWebhookReply, error) {
	return &apiv1.UpdateWebhookReply{}, nil
}

func (s *WebhookService) UpdateWebhookStatus(ctx context.Context, req *apiv1.UpdateWebhookStatusRequest) (*apiv1.UpdateWebhookStatusReply, error) {
	return &apiv1.UpdateWebhookStatusReply{}, nil
}

func (s *WebhookService) DeleteWebhook(ctx context.Context, req *apiv1.DeleteWebhookRequest) (*apiv1.DeleteWebhookReply, error) {
	return &apiv1.DeleteWebhookReply{}, nil
}

func (s *WebhookService) GetWebhook(ctx context.Context, req *apiv1.GetWebhookRequest) (*apiv1.WebhookItem, error) {
	return &apiv1.WebhookItem{}, nil
}

func (s *WebhookService) ListWebhook(ctx context.Context, req *apiv1.ListWebhookRequest) (*apiv1.ListWebhookReply, error) {
	return &apiv1.ListWebhookReply{}, nil
}

func (s *WebhookService) CreateWebhookTemplate(ctx context.Context, req *apiv1.CreateWebhookTemplateRequest) (*apiv1.CreateWebhookTemplateReply, error) {
	return &apiv1.CreateWebhookTemplateReply{}, nil
}

func (s *WebhookService) UpdateWebhookTemplate(ctx context.Context, req *apiv1.UpdateWebhookTemplateRequest) (*apiv1.UpdateWebhookTemplateReply, error) {
	return &apiv1.UpdateWebhookTemplateReply{}, nil
}

func (s *WebhookService) UpdateWebhookTemplateStatus(ctx context.Context, req *apiv1.UpdateWebhookTemplateStatusRequest) (*apiv1.UpdateWebhookTemplateStatusReply, error) {
	return &apiv1.UpdateWebhookTemplateStatusReply{}, nil
}

func (s *WebhookService) DeleteWebhookTemplate(ctx context.Context, req *apiv1.DeleteWebhookTemplateRequest) (*apiv1.DeleteWebhookTemplateReply, error) {
	return &apiv1.DeleteWebhookTemplateReply{}, nil
}

func (s *WebhookService) GetWebhookTemplate(ctx context.Context, req *apiv1.GetWebhookTemplateRequest) (*apiv1.WebhookTemplateItem, error) {
	return &apiv1.WebhookTemplateItem{}, nil
}

func (s *WebhookService) ListWebhookTemplate(ctx context.Context, req *apiv1.ListWebhookTemplateRequest) (*apiv1.ListWebhookTemplateReply, error) {
	return &apiv1.ListWebhookTemplateReply{}, nil
}
