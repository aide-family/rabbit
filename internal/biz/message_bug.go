package biz

import (
	"context"

	"github.com/bwmarrin/snowflake"
	klog "github.com/go-kratos/kratos/v2/log"

	"github.com/aide-family/rabbit/internal/biz/repository"
)

func NewMessageBus(
	messageBusRepo repository.MessageBus,
	helper *klog.Helper,
) *MessageBus {
	return &MessageBus{
		messageBusRepo: messageBusRepo,
		helper:         klog.NewHelper(klog.With(helper.Logger(), "biz", "messageBus")),
	}
}

type MessageBus struct {
	helper         *klog.Helper
	messageBusRepo repository.MessageBus
}

func (m *MessageBus) appendMessage(ctx context.Context, messageUID snowflake.ID) error {
	return m.messageBusRepo.AppendMessage(ctx, messageUID)
}
