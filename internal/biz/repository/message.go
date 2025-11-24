package repository

import (
	"context"

	"github.com/bwmarrin/snowflake"
)

type MessageBus interface {
	AppendMessage(ctx context.Context, messageUID snowflake.ID) error
	SendMessage(ctx context.Context, messageUID snowflake.ID) error
	Stop(ctx context.Context)
	Start(ctx context.Context)
}
