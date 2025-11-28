package service

import (
	"context"

	"github.com/aide-family/rabbit/internal/biz"
)

func NewEventBusService(eventBusBiz *biz.EventBus) *EventBusService {
	return &EventBusService{
		eventBusBiz: eventBusBiz,
	}
}

type EventBusService struct {
	eventBusBiz *biz.EventBus
}

func (s *EventBusService) Start(ctx context.Context) error {
	return s.eventBusBiz.Start(ctx)
}

func (s *EventBusService) Stop(ctx context.Context) error {
	return s.eventBusBiz.Stop(ctx)
}
