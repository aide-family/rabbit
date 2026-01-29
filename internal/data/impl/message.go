package impl

import (
	"context"

	"github.com/bwmarrin/snowflake"

	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/internal/data"
)

func NewMessageRepository(d *data.Data) repository.Message {
	return &messageRepository{}
}

type messageRepository struct{}

// AppendMessage implements [repository.Message].
func (m *messageRepository) AppendMessage(ctx context.Context, messageUID snowflake.ID) error {
	panic("unimplemented")
}

// SendMessage implements [repository.Message].
func (m *messageRepository) SendMessage(ctx context.Context, messageUID snowflake.ID) error {
	panic("unimplemented")
}

// Start implements [repository.Message].
func (m *messageRepository) Start(ctx context.Context) error {
	panic("unimplemented")
}

// Stop implements [repository.Message].
func (m *messageRepository) Stop(ctx context.Context) error {
	panic("unimplemented")
}
