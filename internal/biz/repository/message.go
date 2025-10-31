package repository

import (
	"context"

	"github.com/aide-family/rabbit/internal/biz/do"
)

type MessageBus interface {
	AppendMessage(ctx context.Context, message *do.MessageLog) error
}
