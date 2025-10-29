package service

import (
	"context"

	"github.com/aide-family/rabbit/internal/biz"
	"github.com/aide-family/rabbit/internal/biz/bo"
	apiv1 "github.com/aide-family/rabbit/pkg/api/v1"
)

func NewWebhookService(webhookConfigBiz *biz.WebhookConfig) *WebhookService {
	return &WebhookService{
		webhookConfigBiz: webhookConfigBiz,
	}
}

type WebhookService struct {
	apiv1.UnimplementedWebhookServer
	webhookConfigBiz *biz.WebhookConfig
}

func (s *WebhookService) CreateWebhook(ctx context.Context, req *apiv1.CreateWebhookRequest) (*apiv1.CreateWebhookReply, error) {
	createBo := bo.NewCreateWebhookBo(req)
	if err := s.webhookConfigBiz.CreateWebhook(ctx, createBo); err != nil {
		return nil, err
	}
	return &apiv1.CreateWebhookReply{}, nil
}

func (s *WebhookService) UpdateWebhook(ctx context.Context, req *apiv1.UpdateWebhookRequest) (*apiv1.UpdateWebhookReply, error) {
	updateBo := bo.NewUpdateWebhookBo(req)
	if err := s.webhookConfigBiz.UpdateWebhook(ctx, updateBo); err != nil {
		return nil, err
	}
	return &apiv1.UpdateWebhookReply{}, nil
}

func (s *WebhookService) UpdateWebhookStatus(ctx context.Context, req *apiv1.UpdateWebhookStatusRequest) (*apiv1.UpdateWebhookStatusReply, error) {
	updateBo := bo.NewUpdateWebhookStatusBo(req)
	if err := s.webhookConfigBiz.UpdateWebhookStatus(ctx, updateBo); err != nil {
		return nil, err
	}
	return &apiv1.UpdateWebhookStatusReply{}, nil
}

func (s *WebhookService) DeleteWebhook(ctx context.Context, req *apiv1.DeleteWebhookRequest) (*apiv1.DeleteWebhookReply, error) {
	if err := s.webhookConfigBiz.DeleteWebhook(ctx, req.Uid); err != nil {
		return nil, err
	}
	return &apiv1.DeleteWebhookReply{}, nil
}

func (s *WebhookService) GetWebhook(ctx context.Context, req *apiv1.GetWebhookRequest) (*apiv1.WebhookItem, error) {
	webhookBo, err := s.webhookConfigBiz.GetWebhook(ctx, req.Uid)
	if err != nil {
		return nil, err
	}
	return webhookBo.ToAPIV1WebhookItem(), nil
}

func (s *WebhookService) ListWebhook(ctx context.Context, req *apiv1.ListWebhookRequest) (*apiv1.ListWebhookReply, error) {
	listBo := bo.NewListWebhookBo(req)
	pageResponseBo, err := s.webhookConfigBiz.ListWebhook(ctx, listBo)
	if err != nil {
		return nil, err
	}
	return bo.ToAPIV1ListWebhookReply(pageResponseBo), nil
}

func (s *WebhookService) CreateWebhookTemplate(ctx context.Context, req *apiv1.CreateWebhookTemplateRequest) (*apiv1.CreateWebhookTemplateReply, error) {
	createBo := bo.NewCreateWebhookTemplateBo(req)
	if err := s.webhookConfigBiz.CreateWebhookTemplate(ctx, createBo); err != nil {
		return nil, err
	}
	return &apiv1.CreateWebhookTemplateReply{}, nil
}

func (s *WebhookService) UpdateWebhookTemplate(ctx context.Context, req *apiv1.UpdateWebhookTemplateRequest) (*apiv1.UpdateWebhookTemplateReply, error) {
	updateBo := bo.NewUpdateWebhookTemplateBo(req)
	if err := s.webhookConfigBiz.UpdateWebhookTemplate(ctx, updateBo); err != nil {
		return nil, err
	}
	return &apiv1.UpdateWebhookTemplateReply{}, nil
}

func (s *WebhookService) UpdateWebhookTemplateStatus(ctx context.Context, req *apiv1.UpdateWebhookTemplateStatusRequest) (*apiv1.UpdateWebhookTemplateStatusReply, error) {
	updateBo := bo.NewUpdateWebhookTemplateStatusBo(req)
	if err := s.webhookConfigBiz.UpdateWebhookTemplateStatus(ctx, updateBo); err != nil {
		return nil, err
	}
	return &apiv1.UpdateWebhookTemplateStatusReply{}, nil
}

func (s *WebhookService) DeleteWebhookTemplate(ctx context.Context, req *apiv1.DeleteWebhookTemplateRequest) (*apiv1.DeleteWebhookTemplateReply, error) {
	if err := s.webhookConfigBiz.DeleteWebhookTemplate(ctx, req.Uid); err != nil {
		return nil, err
	}
	return &apiv1.DeleteWebhookTemplateReply{}, nil
}

func (s *WebhookService) GetWebhookTemplate(ctx context.Context, req *apiv1.GetWebhookTemplateRequest) (*apiv1.WebhookTemplateItem, error) {
	templateBo, err := s.webhookConfigBiz.GetWebhookTemplate(ctx, req.Uid)
	if err != nil {
		return nil, err
	}
	return templateBo.ToAPIV1WebhookTemplateItem(), nil
}

func (s *WebhookService) ListWebhookTemplate(ctx context.Context, req *apiv1.ListWebhookTemplateRequest) (*apiv1.ListWebhookTemplateReply, error) {
	listBo := bo.NewListWebhookTemplateBo(req)
	pageResponseBo, err := s.webhookConfigBiz.ListWebhookTemplate(ctx, listBo)
	if err != nil {
		return nil, err
	}
	return bo.ToAPIV1ListWebhookTemplateReply(pageResponseBo), nil
}
