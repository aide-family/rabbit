package service

import (
	"context"

	"github.com/aide-family/rabbit/internal/biz"
)

func NewJobService(jobBiz *biz.Job) *JobService {
	return &JobService{
		jobBiz: jobBiz,
	}
}

type JobService struct {
	jobBiz *biz.Job
}

func (s *JobService) Start(ctx context.Context) error {
	return s.jobBiz.Start(ctx)
}

func (s *JobService) Stop(ctx context.Context) error {
	return s.jobBiz.Stop(ctx)
}
