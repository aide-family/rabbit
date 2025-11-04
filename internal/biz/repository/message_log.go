package repository

import (
	"context"

	"github.com/bwmarrin/snowflake"

	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/do"
	"github.com/aide-family/rabbit/internal/biz/vobj"
)

type MessageLog interface {
	CreateMessageLog(ctx context.Context, messageLog *do.MessageLog) error
	ListMessageLog(ctx context.Context, req *bo.ListMessageLogBo) (*bo.PageResponseBo[*do.MessageLog], error)
	GetMessageLog(ctx context.Context, uid snowflake.ID) (*do.MessageLog, error)
	UpdateMessageLogStatus(ctx context.Context, uid snowflake.ID, status vobj.MessageStatus) error
	// GetMessageLogWithLock 使用 SELECT FOR UPDATE 获取消息日志并加锁，用于分布式锁场景
	GetMessageLogWithLock(ctx context.Context, uid snowflake.ID) (*do.MessageLog, error)
	// UpdateMessageLogStatusIf 条件更新消息状态，只有当前状态匹配时才更新，用于实现 CAS 操作
	UpdateMessageLogStatusIf(ctx context.Context, uid snowflake.ID, oldStatus, newStatus vobj.MessageStatus) (bool, error)
}
