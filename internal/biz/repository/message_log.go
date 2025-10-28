package repository

import (
	"context"

	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/do"
)

type MessageLog interface {
	CreateMessageLog(ctx context.Context, messageLog *do.MessageLog) error
	ListMessageLog(ctx context.Context, req *bo.ListMessageLogBo) (*bo.PageResponseBo[*do.MessageLog], error)
}
