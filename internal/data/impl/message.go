package impl

import (
	"context"

	"github.com/aide-family/rabbit/internal/biz/do"
	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/internal/data"
)

func NewMessageBus(d *data.Data) repository.MessageBus {
	return &messageBusImpl{d: d}
}

type messageBusImpl struct {
	d *data.Data
}

// AppendMessage implements repository.MessageBus.
func (m *messageBusImpl) AppendMessage(ctx context.Context, message *do.MessageLog) error {
	return nil
}
