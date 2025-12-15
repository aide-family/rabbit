package biz

import (
	"context"

	"github.com/bwmarrin/snowflake"
	klog "github.com/go-kratos/kratos/v2/log"

	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/do"
	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/internal/biz/vobj"
	"github.com/aide-family/rabbit/pkg/merr"
)

func NewMessageLog(
	messageLogRepo repository.MessageLog,
	jobBiz *Job,
	helper *klog.Helper,
) *MessageLog {
	return &MessageLog{
		messageLogRepo: messageLogRepo,
		jobBiz:         jobBiz,
		helper:         klog.NewHelper(klog.With(helper.Logger(), "biz", "messageLog")),
	}
}

type MessageLog struct {
	helper         *klog.Helper
	messageLogRepo repository.MessageLog
	jobBiz         *Job
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
		if merr.IsNotFound(err) {
			return nil, err
		}
		m.helper.Errorw("msg", "get message log failed", "error", err, "uid", uid)
		return nil, merr.ErrorInternal("get message log failed")
	}
	return bo.NewMessageLogItemBo(messageLogDO), nil
}

func (m *MessageLog) RetryMessage(ctx context.Context, uid snowflake.ID) error {
	messageLog, err := m.messageLogRepo.GetMessageLogWithLock(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return err
		}
		m.helper.Errorw("msg", "get message log failed", "error", err, "uid", uid)
		return merr.ErrorInternal("get message log failed")
	}
	if messageLog.Status.IsSent() || messageLog.Status.IsSending() || messageLog.Status.IsCancelled() {
		return nil
	}
	if err := m.jobBiz.appendMessage(ctx, uid); err != nil {
		m.helper.Errorw("msg", "append message failed", "error", err, "uid", uid)
		return merr.ErrorInternal("append message failed")
	}
	return nil
}

func (m *MessageLog) CancelMessage(ctx context.Context, uid snowflake.ID) error {
	messageLog, err := m.messageLogRepo.GetMessageLogWithLock(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return err
		}
		m.helper.Errorw("msg", "get message log failed", "error", err, "uid", uid)
		return merr.ErrorInternal("get message log failed")
	}
	if messageLog.Status.IsSent() || messageLog.Status.IsCancelled() {
		return merr.ErrorNotFound("message already sent or cancelled")
	}
	success, err := m.messageLogRepo.UpdateMessageLogStatusIf(ctx, uid, messageLog.Status, vobj.MessageStatusCancelled)
	if err != nil {
		m.helper.Errorw("msg", "update message status to cancelled failed", "error", err, "uid", uid)
		return merr.ErrorInternal("cancel message failed")
	}
	if !success {
		m.helper.Warnw("msg", "message status is not sending, message cancelled failed", "uid", uid)
		return merr.ErrorNotFound("cancel message failed, the status of this message has changed.")
	}
	return nil
}

func (m *MessageLog) createMessageLog(ctx context.Context, messageLog *do.MessageLog) error {
	// TODO 区分数据存储实现， 如果未启用数据库，则使用文件存储
	return m.messageLogRepo.CreateMessageLog(ctx, messageLog)
}
