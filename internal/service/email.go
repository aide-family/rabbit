package service

import (
	"context"

	pb "github.com/aide-family/rabbit/pkg/api/v1"
)

type EmailService struct {
	pb.UnimplementedEmailServer
}

func NewEmailService() *EmailService {
	return &EmailService{}
}

func (s *EmailService) CreateEmailConfig(ctx context.Context, req *pb.CreateEmailConfigRequest) (*pb.CreateEmailConfigReply, error) {
	return &pb.CreateEmailConfigReply{}, nil
}
func (s *EmailService) UpdateEmailConfig(ctx context.Context, req *pb.UpdateEmailConfigRequest) (*pb.UpdateEmailConfigReply, error) {
	return &pb.UpdateEmailConfigReply{}, nil
}
func (s *EmailService) DeleteEmailConfig(ctx context.Context, req *pb.DeleteEmailConfigRequest) (*pb.DeleteEmailConfigReply, error) {
	return &pb.DeleteEmailConfigReply{}, nil
}
func (s *EmailService) GetEmailConfig(ctx context.Context, req *pb.GetEmailConfigRequest) (*pb.EmailConfigItem, error) {
	return &pb.EmailConfigItem{}, nil
}
func (s *EmailService) ListEmailConfig(ctx context.Context, req *pb.ListEmailConfigRequest) (*pb.ListEmailConfigReply, error) {
	return &pb.ListEmailConfigReply{}, nil
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
