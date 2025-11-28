package biz

import (
	"context"

	"github.com/bwmarrin/snowflake"
	klog "github.com/go-kratos/kratos/v2/log"

	"github.com/aide-family/rabbit/internal/biz/repository"
)

func NewEventBus(
	messageRepo repository.Message,
	helper *klog.Helper,
) *EventBus {
	return &EventBus{
		messageRepo: messageRepo,
		helper:      klog.NewHelper(klog.With(helper.Logger(), "biz", "eventBus")),
	}
}

type EventBus struct {
	helper      *klog.Helper
	messageRepo repository.Message
}

func (e *EventBus) appendMessage(ctx context.Context, messageUID snowflake.ID) error {
	return e.messageRepo.AppendMessage(ctx, messageUID)
}

func (e *EventBus) Start(ctx context.Context) error {
	return e.messageRepo.Start(ctx)
}

func (e *EventBus) Stop(ctx context.Context) error {
	return e.messageRepo.Stop(ctx)
}
