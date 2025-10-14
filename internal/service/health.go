package service

import (
	"context"

	apiv1 "github.com/aide-family/rabbit/api/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type HealthService struct {
	apiv1.UnimplementedHealthServer
}

func NewHealthService() *HealthService {
	return &HealthService{}
}

func (s *HealthService) HealthCheck(ctx context.Context, req *apiv1.HealthCheckRequest) (*apiv1.HealthCheckReply, error) {
	return &apiv1.HealthCheckReply{
		Status:    "OK",
		Message:   "Rabbit is running",
		Timestamp: timestamppb.Now(),
	}, nil
}
