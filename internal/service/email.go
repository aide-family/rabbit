package service

import (
	"context"

	"github.com/aide-family/rabbit/internal/biz"
	"github.com/aide-family/rabbit/internal/biz/bo"
	apiv1 "github.com/aide-family/rabbit/pkg/api/v1"
)

func NewEmailService(emailConfigBiz *biz.EmailConfig) *EmailService {
	return &EmailService{
		emailConfigBiz: emailConfigBiz,
	}
}

type EmailService struct {
	apiv1.UnimplementedEmailServer

	emailConfigBiz *biz.EmailConfig
}

func (s *EmailService) CreateEmailConfig(ctx context.Context, req *apiv1.CreateEmailConfigRequest) (*apiv1.CreateEmailConfigReply, error) {
	createEmailConfigBo := bo.NewCreateEmailConfigBo(req)
	if err := s.emailConfigBiz.CreateEmailConfig(ctx, createEmailConfigBo); err != nil {
		return nil, err
	}
	return &apiv1.CreateEmailConfigReply{}, nil
}

func (s *EmailService) UpdateEmailConfig(ctx context.Context, req *apiv1.UpdateEmailConfigRequest) (*apiv1.UpdateEmailConfigReply, error) {
	updateEmailConfigBo := bo.NewUpdateEmailConfigBo(req)
	if err := s.emailConfigBiz.UpdateEmailConfig(ctx, updateEmailConfigBo); err != nil {
		return nil, err
	}
	return &apiv1.UpdateEmailConfigReply{}, nil
}

func (s *EmailService) UpdateEmailConfigStatus(ctx context.Context, req *apiv1.UpdateEmailConfigStatusRequest) (*apiv1.UpdateEmailConfigStatusReply, error) {
	updateEmailConfigStatusBo := bo.NewUpdateEmailConfigStatusBo(req)
	if err := s.emailConfigBiz.UpdateEmailConfigStatus(ctx, updateEmailConfigStatusBo); err != nil {
		return nil, err
	}
	return &apiv1.UpdateEmailConfigStatusReply{}, nil
}

func (s *EmailService) DeleteEmailConfig(ctx context.Context, req *apiv1.DeleteEmailConfigRequest) (*apiv1.DeleteEmailConfigReply, error) {
	if err := s.emailConfigBiz.DeleteEmailConfig(ctx, req.Uid); err != nil {
		return nil, err
	}
	return &apiv1.DeleteEmailConfigReply{}, nil
}

func (s *EmailService) GetEmailConfig(ctx context.Context, req *apiv1.GetEmailConfigRequest) (*apiv1.EmailConfigItem, error) {
	getEmailConfigBo, err := s.emailConfigBiz.GetEmailConfig(ctx, req.Uid)
	if err != nil {
		return nil, err
	}
	return getEmailConfigBo.ToAPIV1EmailConfigItem(), nil
}

func (s *EmailService) ListEmailConfig(ctx context.Context, req *apiv1.ListEmailConfigRequest) (*apiv1.ListEmailConfigReply, error) {
	emailConfigListPageRequestBo := bo.NewListEmailConfigBo(req)
	emailConfigListPageResponseBo, err := s.emailConfigBiz.ListEmailConfig(ctx, emailConfigListPageRequestBo)
	if err != nil {
		return nil, err
	}
	return bo.ToAPIV1ListEmailConfigReply(emailConfigListPageResponseBo), nil
}
