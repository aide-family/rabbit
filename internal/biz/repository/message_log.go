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
}
