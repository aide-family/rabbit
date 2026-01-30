package repository

import (
	"context"

	"github.com/bwmarrin/snowflake"

	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/pkg/enum"
)

type MessageLog interface {
	CreateMessageLog(ctx context.Context, messageLog *bo.MessageLogItemBo) error
	ListMessageLog(ctx context.Context, req *bo.ListMessageLogBo) (*bo.PageResponseBo[*bo.MessageLogItemBo], error)
	GetMessageLog(ctx context.Context, uid snowflake.ID) (*bo.MessageLogItemBo, error)
	GetAllMessageLogs(ctx context.Context, status enum.MessageStatus) ([]*bo.MessageLogItemBo, error)
	// GetMessageLogWithLock 使用 SELECT FOR UPDATE 获取消息日志并加锁，用于分布式锁场景
	GetMessageLogWithLock(ctx context.Context, uid snowflake.ID) (*bo.MessageLogItemBo, error)
	// UpdateMessageLogStatusIf 条件更新消息状态，只有当前状态匹配时才更新，用于实现 CAS 操作
	UpdateMessageLogStatusIf(ctx context.Context, uid snowflake.ID, oldStatus, newStatus enum.MessageStatus) (bool, error)
	// UpdateMessageLogLastErrorIf 条件更新消息最后错误，只有当前状态匹配时才更新，用于实现 CAS 操作
	UpdateMessageLogLastErrorIf(ctx context.Context, uid snowflake.ID, oldStatus enum.MessageStatus, lastError string) (bool, error)
}
