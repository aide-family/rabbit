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

func (s *EmailService) CreateEmailTemplate(ctx context.Context, req *apiv1.CreateEmailTemplateRequest) (*apiv1.CreateEmailTemplateReply, error) {
	createEmailTemplateBo := bo.NewCreateEmailTemplateBo(req)
	if err := s.emailConfigBiz.CreateEmailTemplate(ctx, createEmailTemplateBo); err != nil {
		return nil, err
	}
	return &apiv1.CreateEmailTemplateReply{}, nil
}

func (s *EmailService) UpdateEmailTemplate(ctx context.Context, req *apiv1.UpdateEmailTemplateRequest) (*apiv1.UpdateEmailTemplateReply, error) {
	updateEmailTemplateBo := bo.NewUpdateEmailTemplateBo(req)
	if err := s.emailConfigBiz.UpdateEmailTemplate(ctx, updateEmailTemplateBo); err != nil {
		return nil, err
	}
	return &apiv1.UpdateEmailTemplateReply{}, nil
}

func (s *EmailService) UpdateEmailTemplateStatus(ctx context.Context, req *apiv1.UpdateEmailTemplateStatusRequest) (*apiv1.UpdateEmailTemplateStatusReply, error) {
	updateEmailTemplateStatusBo := bo.NewUpdateEmailTemplateStatusBo(req)
	if err := s.emailConfigBiz.UpdateEmailTemplateStatus(ctx, updateEmailTemplateStatusBo); err != nil {
		return nil, err
	}
	return &apiv1.UpdateEmailTemplateStatusReply{}, nil
}

func (s *EmailService) DeleteEmailTemplate(ctx context.Context, req *apiv1.DeleteEmailTemplateRequest) (*apiv1.DeleteEmailTemplateReply, error) {
	if err := s.emailConfigBiz.DeleteEmailTemplate(ctx, req.Uid); err != nil {
		return nil, err
	}
	return &apiv1.DeleteEmailTemplateReply{}, nil
}

func (s *EmailService) GetEmailTemplate(ctx context.Context, req *apiv1.GetEmailTemplateRequest) (*apiv1.EmailTemplateItem, error) {
	getEmailTemplateBo, err := s.emailConfigBiz.GetEmailTemplate(ctx, req.Uid)
	if err != nil {
		return nil, err
	}
	return getEmailTemplateBo.ToAPIV1EmailTemplateItem(), nil
}

func (s *EmailService) ListEmailTemplate(ctx context.Context, req *apiv1.ListEmailTemplateRequest) (*apiv1.ListEmailTemplateReply, error) {
	emailTemplateListPageRequestBo := bo.NewListEmailTemplateBo(req)
	emailTemplateListPageResponseBo, err := s.emailConfigBiz.ListEmailTemplate(ctx, emailTemplateListPageRequestBo)
	if err != nil {
		return nil, err
	}
	return bo.ToAPIV1ListEmailTemplateReply(emailTemplateListPageResponseBo), nil
}
