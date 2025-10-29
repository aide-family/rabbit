package repository

import (
	"context"
	"time"

	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/do"
	"github.com/aide-family/rabbit/internal/biz/vobj"
)

type MessageLog interface {
	CreateMessageLog(ctx context.Context, messageLog *do.MessageLog) error
	ListMessageLog(ctx context.Context, req *bo.ListMessageLogBo) (*bo.PageResponseBo[*do.MessageLog], error)
	GetMessageLog(ctx context.Context, uid string, sendAt time.Time) (*do.MessageLog, error)
	UpdateMessageLogStatus(ctx context.Context, uid string, sendAt time.Time, status vobj.MessageStatus) error
}
