package biz

import (
	"context"

	"github.com/bwmarrin/snowflake"
	klog "github.com/go-kratos/kratos/v2/log"

	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/pkg/merr"
)

func NewMessageLog(messageLogRepo repository.MessageLog, helper *klog.Helper) *MessageLog {
	return &MessageLog{
		messageLogRepo: messageLogRepo,
		helper:         klog.NewHelper(klog.With(helper.Logger(), "biz", "messageLog")),
	}
}

type MessageLog struct {
	helper         *klog.Helper
	messageLogRepo repository.MessageLog
}

func (m *MessageLog) ListMessageLog(ctx context.Context, req *bo.ListMessageLogBo) (*bo.PageResponseBo[*bo.MessageLogItemBo], error) {
	pageResponseBo, err := m.messageLogRepo.ListMessageLog(ctx, req)
	if err != nil {
		m.helper.Errorw("msg", "list message log failed", "error", err)
		return nil, merr.ErrorInternal("list message log failed")
	}
	items := make([]*bo.MessageLogItemBo, 0, len(pageResponseBo.GetItems()))
	for _, item := range pageResponseBo.GetItems() {
		items = append(items, bo.NewMessageLogItemBo(item))
	}
	return bo.NewPageResponseBo(pageResponseBo.PageRequestBo, items), nil
}

func (m *MessageLog) GetMessageLog(ctx context.Context, uid snowflake.ID) (*bo.MessageLogItemBo, error) {
	messageLogDO, err := m.messageLogRepo.GetMessageLog(ctx, uid)
	if err != nil {
		m.helper.Errorw("msg", "get message log failed", "error", err, "uid", uid)
		return nil, merr.ErrorInternal("get message log failed")
	}
	return bo.NewMessageLogItemBo(messageLogDO), nil
}

func (m *MessageLog) RetryMessage(ctx context.Context, uid snowflake.ID) error {
	return nil
}

func (m *MessageLog) CancelMessage(ctx context.Context, uid snowflake.ID) error {
	return nil
}
