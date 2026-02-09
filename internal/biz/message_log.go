package biz

import (
	"context"
	"slices"

	"github.com/aide-family/magicbox/enum"
	"github.com/aide-family/magicbox/merr"
	"github.com/bwmarrin/snowflake"
	klog "github.com/go-kratos/kratos/v2/log"

	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/repository"
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
		return nil, merr.ErrorInternalServer("list message log failed")
	}
	return pageResponseBo, nil
}

func (m *MessageLog) GetMessageLog(ctx context.Context, uid snowflake.ID) (*bo.MessageLogItemBo, error) {
	messageLogBo, err := m.messageLogRepo.GetMessageLog(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return nil, err
		}
		m.helper.Errorw("msg", "get message log failed", "error", err, "uid", uid)
		return nil, merr.ErrorInternalServer("get message log failed")
	}
	return messageLogBo, nil
}

func (m *MessageLog) RetryMessage(ctx context.Context, uid snowflake.ID) error {
	messageLog, err := m.messageLogRepo.GetMessageLogWithLock(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return err
		}
		m.helper.Errorw("msg", "get message log failed", "error", err, "uid", uid)
		return merr.ErrorInternalServer("get message log failed")
	}
	if slices.Contains([]enum.MessageStatus{enum.MessageStatus_SENT, enum.MessageStatus_SENDING, enum.MessageStatus_CANCELLED}, messageLog.Status) {
		m.helper.Debugw("msg", "message already sent or sending or cancelled", "uid", uid, "status", messageLog.Status)
		return nil
	}
	if err := m.jobBiz.AppendMessage(ctx, uid); err != nil {
		m.helper.Errorw("msg", "append message failed", "error", err, "uid", uid)
		return merr.ErrorInternalServer("append message failed")
	}
	if err := m.messageLogRepo.MessageLogRetryIncrement(ctx, uid); err != nil {
		m.helper.Warnw("msg", "increment message retry failed", "error", err, "uid", uid)
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
		return merr.ErrorInternalServer("get message log failed")
	}
	if slices.Contains([]enum.MessageStatus{enum.MessageStatus_SENT, enum.MessageStatus_CANCELLED}, messageLog.Status) {
		return merr.ErrorParams("message already sent or cancelled")
	}
	success, err := m.messageLogRepo.UpdateMessageLogStatusIf(ctx, uid, messageLog.Status, enum.MessageStatus_CANCELLED)
	if err != nil {
		m.helper.Errorw("msg", "update message status to cancelled failed", "error", err, "uid", uid)
		return merr.ErrorInternalServer("cancel message failed")
	}
	if !success {
		m.helper.Warnw("msg", "message status is not sending, message cancelled failed", "uid", uid)
		return merr.ErrorNotFound("cancel message failed, the status of this message has changed.")
	}
	return nil
}

func (m *MessageLog) createMessageLog(ctx context.Context, messageLog *bo.MessageLogItemBo) error {
	return m.messageLogRepo.CreateMessageLog(ctx, messageLog)
}
