package service

import (
	"context"
	"time"

	"github.com/aide-family/rabbit/internal/biz"
	"github.com/aide-family/rabbit/internal/biz/bo"
	apiv1 "github.com/aide-family/rabbit/pkg/api/v1"
)

func NewMessageLogService(messageLogBiz *biz.MessageLog) *MessageLogService {
	return &MessageLogService{
		messageLogBiz: messageLogBiz,
	}
}

type MessageLogService struct {
	apiv1.UnimplementedMessageLogServer
	messageLogBiz *biz.MessageLog
}

func (s *MessageLogService) RetryMessage(ctx context.Context, req *apiv1.RetryMessageLogRequest) (*apiv1.RetryMessageLogReply, error) {
	sendAt := time.Unix(req.SendAtUnix, 0)
	err := s.messageLogBiz.RetryMessage(ctx, req.Uid, sendAt)
	if err != nil {
		return nil, err
	}
	return &apiv1.RetryMessageLogReply{}, nil
}

func (s *MessageLogService) CancelMessage(ctx context.Context, req *apiv1.CancelMessageLogRequest) (*apiv1.CancelMessageLogReply, error) {
	sendAt := time.Unix(req.SendAtUnix, 0)
	err := s.messageLogBiz.CancelMessage(ctx, req.Uid, sendAt)
	if err != nil {
		return nil, err
	}
	return &apiv1.CancelMessageLogReply{}, nil
}

func (s *MessageLogService) GetMessageLog(ctx context.Context, req *apiv1.GetMessageLogRequest) (*apiv1.MessageLogItem, error) {
	sendAt := time.Unix(req.SendAtUnix, 0)
	messageLogBo, err := s.messageLogBiz.GetMessageLog(ctx, req.Uid, sendAt)
	if err != nil {
		return nil, err
	}
	return messageLogBo.ToAPIV1MessageLogItem(), nil
}

func (s *MessageLogService) ListMessageLog(ctx context.Context, req *apiv1.ListMessageLogRequest) (*apiv1.ListMessageLogReply, error) {
	listBo := bo.NewListMessageLogBo(req)
	pageResponseBo, err := s.messageLogBiz.ListMessageLog(ctx, listBo)
	if err != nil {
		return nil, err
	}
	return bo.ToAPIV1ListMessageLogReply(pageResponseBo), nil
}
