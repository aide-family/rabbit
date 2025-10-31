package server

import (
	"context"

	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"

	"github.com/aide-family/rabbit/internal/biz/repository"
)

var _ transport.Server = (*EventBus)(nil)

func NewEventBus(messageBusRepo repository.MessageBus, helper *klog.Helper) *EventBus {
	return &EventBus{
		messageBusRepo: messageBusRepo,
		helper:         helper,
	}
}

type EventBus struct {
	messageBusRepo repository.MessageBus
	helper         *klog.Helper
}

// Start implements transport.Server.
func (e *EventBus) Start(context.Context) error {
	e.messageBusRepo.Start()
	e.helper.Infow("msg", "[EventBus] started")
	return nil
}

// Stop implements transport.Server.
func (e *EventBus) Stop(context.Context) error {
	e.messageBusRepo.Stop()
	e.helper.Infow("msg", "[EventBus] stopped")
	return nil
}
