package service

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/aide-family/rabbit/internal/biz"
	apiv1 "github.com/aide-family/rabbit/pkg/api/v1"
)

func NewHealthService(healthBiz *biz.Health) *HealthService {
	return &HealthService{
		healthBiz: healthBiz,
	}
}

type HealthService struct {
	apiv1.UnimplementedHealthServer

	healthBiz *biz.Health
}

func (s *HealthService) HealthCheck(ctx context.Context, req *apiv1.HealthCheckRequest) (*apiv1.HealthCheckReply, error) {
	return &apiv1.HealthCheckReply{
		Status:    "OK",
		Message:   "Rabbit is running",
		Timestamp: timestamppb.Now(),
	}, nil
}
