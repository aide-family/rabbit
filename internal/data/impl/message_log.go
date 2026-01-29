package impl

import (
	"context"

	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/internal/data"
	"github.com/aide-family/rabbit/pkg/enum"
	"github.com/bwmarrin/snowflake"
)

func NewMessageLogRepository(d *data.Data) repository.MessageLog {
	return &messageLogRepository{}
}

type messageLogRepository struct{}

// GetMessageLog implements [repository.MessageLog].
func (m *messageLogRepository) GetMessageLog(ctx context.Context, uid snowflake.ID) (*bo.MessageLogItemBo, error) {
	panic("unimplemented")
}

// GetMessageLogWithLock implements [repository.MessageLog].
func (m *messageLogRepository) GetMessageLogWithLock(ctx context.Context, uid snowflake.ID) (*bo.MessageLogItemBo, error) {
	panic("unimplemented")
}

// ListMessageLog implements [repository.MessageLog].
func (m *messageLogRepository) ListMessageLog(ctx context.Context, req *bo.ListMessageLogBo) (*bo.PageResponseBo[*bo.MessageLogItemBo], error) {
	panic("unimplemented")
}

// UpdateMessageLogStatusIf implements [repository.MessageLog].
func (m *messageLogRepository) UpdateMessageLogStatusIf(ctx context.Context, uid snowflake.ID, oldStatus enum.MessageStatus, newStatus enum.MessageStatus) (bool, error) {
	panic("unimplemented")
}

// CreateMessageLog implements [repository.MessageLog].
func (m *messageLogRepository) CreateMessageLog(ctx context.Context, messageLog *bo.MessageLogItemBo) error {
	panic("unimplemented")
}
