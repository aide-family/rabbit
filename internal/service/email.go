package service

import (
	"context"

	"github.com/aide-family/rabbit/internal/biz"
	"github.com/aide-family/rabbit/internal/biz/bo"
	pb "github.com/aide-family/rabbit/pkg/api/v1"
)

func NewEmailService(emailConfigBiz *biz.EmailConfig) *EmailService {
	return &EmailService{
		emailConfigBiz: emailConfigBiz,
	}
}

type EmailService struct {
	pb.UnimplementedEmailServer

	emailConfigBiz *biz.EmailConfig
}

func (s *EmailService) CreateEmailConfig(ctx context.Context, req *pb.CreateEmailConfigRequest) (*pb.CreateEmailConfigReply, error) {
	createEmailConfigBo := bo.NewCreateEmailConfigBo(req)
	if err := s.emailConfigBiz.CreateEmailConfig(ctx, createEmailConfigBo); err != nil {
		return nil, err
	}
	return &pb.CreateEmailConfigReply{}, nil
}

func (s *EmailService) UpdateEmailConfig(ctx context.Context, req *pb.UpdateEmailConfigRequest) (*pb.UpdateEmailConfigReply, error) {
	updateEmailConfigBo := bo.NewUpdateEmailConfigBo(req)
	if err := s.emailConfigBiz.UpdateEmailConfig(ctx, updateEmailConfigBo); err != nil {
		return nil, err
	}
	return &pb.UpdateEmailConfigReply{}, nil
}

func (s *EmailService) UpdateEmailConfigStatus(ctx context.Context, req *pb.UpdateEmailConfigStatusRequest) (*pb.UpdateEmailConfigStatusReply, error) {
	updateEmailConfigStatusBo := bo.NewUpdateEmailConfigStatusBo(req)
	if err := s.emailConfigBiz.UpdateEmailConfigStatus(ctx, updateEmailConfigStatusBo); err != nil {
		return nil, err
	}
	return &pb.UpdateEmailConfigStatusReply{}, nil
}

func (s *EmailService) DeleteEmailConfig(ctx context.Context, req *pb.DeleteEmailConfigRequest) (*pb.DeleteEmailConfigReply, error) {
	if err := s.emailConfigBiz.DeleteEmailConfig(ctx, req.Uid); err != nil {
		return nil, err
	}
	return &pb.DeleteEmailConfigReply{}, nil
}

func (s *EmailService) GetEmailConfig(ctx context.Context, req *pb.GetEmailConfigRequest) (*pb.EmailConfigItem, error) {
	getEmailConfigBo, err := s.emailConfigBiz.GetEmailConfig(ctx, req.Uid)
	if err != nil {
		return nil, err
	}
	return getEmailConfigBo.ToAPIV1EmailConfigItem(), nil
}

func (s *EmailService) ListEmailConfig(ctx context.Context, req *pb.ListEmailConfigRequest) (*pb.ListEmailConfigReply, error) {
	emailConfigListPageRequestBo := bo.NewListEmailConfigBo(req)
	emailConfigListPageResponseBo, err := s.emailConfigBiz.ListEmailConfig(ctx, emailConfigListPageRequestBo)
	if err != nil {
		return nil, err
	}
	return bo.ToAPIV1ListEmailConfigReply(emailConfigListPageResponseBo), nil
}

func (s *EmailService) CreateEmailTemplate(ctx context.Context, req *pb.CreateEmailTemplateRequest) (*pb.CreateEmailTemplateReply, error) {
	return &pb.CreateEmailTemplateReply{}, nil
}

func (s *EmailService) UpdateEmailTemplate(ctx context.Context, req *pb.UpdateEmailTemplateRequest) (*pb.UpdateEmailTemplateReply, error) {
	return &pb.UpdateEmailTemplateReply{}, nil
}

func (s *EmailService) UpdateEmailTemplateStatus(ctx context.Context, req *pb.UpdateEmailTemplateStatusRequest) (*pb.UpdateEmailTemplateStatusReply, error) {
	return &pb.UpdateEmailTemplateStatusReply{}, nil
}

func (s *EmailService) DeleteEmailTemplate(ctx context.Context, req *pb.DeleteEmailTemplateRequest) (*pb.DeleteEmailTemplateReply, error) {
	return &pb.DeleteEmailTemplateReply{}, nil
}

func (s *EmailService) GetEmailTemplate(ctx context.Context, req *pb.GetEmailTemplateRequest) (*pb.EmailTemplateItem, error) {
	return &pb.EmailTemplateItem{}, nil
}

func (s *EmailService) ListEmailTemplate(ctx context.Context, req *pb.ListEmailTemplateRequest) (*pb.ListEmailTemplateReply, error) {
	return &pb.ListEmailTemplateReply{}, nil
}
