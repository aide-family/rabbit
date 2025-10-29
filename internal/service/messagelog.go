package service

import (
	"context"

	apiv1 "github.com/aide-family/rabbit/pkg/api/v1"
)

type MessageLogService struct {
	apiv1.UnimplementedMessageLogServer
}

func NewMessageLogService() *MessageLogService {
	return &MessageLogService{}
}

func (s *MessageLogService) RetryMessage(ctx context.Context, req *apiv1.RetryMessageLogRequest) (*apiv1.RetryMessageLogReply, error) {
	return &apiv1.RetryMessageLogReply{}, nil
}

func (s *MessageLogService) CancelMessage(ctx context.Context, req *apiv1.CancelMessageLogRequest) (*apiv1.CancelMessageLogReply, error) {
	return &apiv1.CancelMessageLogReply{}, nil
}

func (s *MessageLogService) GetMessageLog(ctx context.Context, req *apiv1.GetMessageLogRequest) (*apiv1.MessageLogItem, error) {
	return &apiv1.MessageLogItem{}, nil
}

func (s *MessageLogService) ListMessageLog(ctx context.Context, req *apiv1.ListMessageLogRequest) (*apiv1.ListMessageLogReply, error) {
	return &apiv1.ListMessageLogReply{}, nil
}
