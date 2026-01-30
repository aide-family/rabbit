package impl

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aide-family/magicbox/plugin/cache"
	"github.com/aide-family/magicbox/pointer"
	"github.com/bwmarrin/snowflake"
	klog "github.com/go-kratos/kratos/v2/log"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/internal/data"
	"github.com/aide-family/rabbit/internal/data/impl/convert"
	"github.com/aide-family/rabbit/internal/data/impl/do"
	"github.com/aide-family/rabbit/internal/data/impl/query"
	"github.com/aide-family/rabbit/pkg/contextx"
	"github.com/aide-family/rabbit/pkg/enum"
	"github.com/aide-family/rabbit/pkg/merr"
)

func NewMessageLogRepository(d *data.Data) repository.MessageLog {
	query.SetDefault(d.DB())
	return &messageLogRepository{Data: d}
}

type messageLogRepository struct {
	*data.Data
}

// GetMessageLog implements [repository.MessageLog].
func (m *messageLogRepository) GetMessageLog(ctx context.Context, uid snowflake.ID) (*bo.MessageLogItemBo, error) {
	return m.getMessageLog(ctx, uid)
}

// GetAllMessageLogs implements [repository.MessageLog].
func (m *messageLogRepository) GetAllMessageLogs(ctx context.Context, status enum.MessageStatus) ([]*bo.MessageLogItemBo, error) {
	namespace := contextx.GetNamespaceUID(ctx)
	tableName := do.GenMessageLogTableName(namespace, time.Now())
	if _, err := m.Cache().Get(ctx, cache.K(tableName)); err != nil && !do.HasTable(m.DB(), tableName) {
		return []*bo.MessageLogItemBo{}, nil
	}

	bizQuery := query.Use(m.DB().Table(tableName))
	messageLog := bizQuery.MessageLog
	messageLogTable := messageLog.As(tableName)
	wrappers := messageLog.WithContext(ctx)
	wheres := []gen.Condition{
		messageLogTable.Status.Eq(int32(status)),
	}
	wrappers = wrappers.Where(wheres...)
	messageLogs, err := wrappers.Order(messageLog.CreatedAt.Asc()).Find()
	if err != nil {
		return nil, err
	}
	messageLogItems := make([]*bo.MessageLogItemBo, 0, len(messageLogs))
	for _, messageLog := range messageLogs {
		messageLogItems = append(messageLogItems, convert.ToMessageLogItemBo(messageLog))
	}
	return messageLogItems, nil
}

// GetMessageLogWithLock implements [repository.MessageLog].
func (m *messageLogRepository) GetMessageLogWithLock(ctx context.Context, uid snowflake.ID) (*bo.MessageLogItemBo, error) {
	return m.getMessageLog(ctx, uid, clause.Locking{Strength: "UPDATE"})
}

func (m *messageLogRepository) getMessageLog(ctx context.Context, uid snowflake.ID, clauses ...clause.Expression) (*bo.MessageLogItemBo, error) {
	namespace := contextx.GetNamespaceUID(ctx)
	tableName := do.GenMessageLogTableName(namespace, time.UnixMilli(uid.Time()))
	if _, err := m.Cache().Get(ctx, cache.K(tableName)); err != nil && !do.HasTable(m.DB(), tableName) {
		return nil, gorm.ErrRecordNotFound
	}

	bizQuery := query.Use(m.DB().Table(tableName))
	messageLog := bizQuery.MessageLog
	messageLogTable := messageLog.As(tableName)
	wrappers := messageLog.WithContext(ctx)
	wheres := []gen.Condition{
		messageLogTable.UID.Eq(uid.Int64()),
		messageLogTable.NamespaceUID.Eq(namespace.Int64()),
	}
	wrappers = wrappers.Where(wheres...).Clauses(clauses...)
	messageLogDo, err := wrappers.First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, merr.ErrorNotFound("message log %d not found", uid.Int64())
		}
		return nil, err
	}
	return convert.ToMessageLogItemBo(messageLogDo), nil
}

// ListMessageLog implements [repository.MessageLog].
func (m *messageLogRepository) ListMessageLog(ctx context.Context, req *bo.ListMessageLogBo) (*bo.PageResponseBo[*bo.MessageLogItemBo], error) {
	namespace := contextx.GetNamespaceUID(ctx)

	if req.StartAt.IsZero() {
		req.StartAt = time.Now().AddDate(0, 0, -7)
	}
	if req.EndAt.IsZero() {
		req.EndAt = time.Now()
	}
	bizDB := m.DB()

	tableNames := do.GenMessageLogTableNames(bizDB, namespace, req.StartAt, req.EndAt)
	if len(tableNames) == 0 {
		return bo.NewPageResponseBo[*bo.MessageLogItemBo](req.PageRequestBo, nil), nil
	}

	tables := make([]any, 0, len(tableNames))
	unionAllSQL := make([]string, 0, len(tableNames))
	for _, tableName := range tableNames {
		tables = append(tables, bizDB.Table(tableName))
		unionAllSQL = append(unionAllSQL, "?")
	}
	wrappers := bizDB.WithContext(ctx)
	if len(tableNames) > 1 {
		wrappers = wrappers.Table(fmt.Sprintf("(%s) as %s", strings.Join(unionAllSQL, " UNION ALL "), do.TableNameMessageLog), tables...)
	} else {
		wrappers = wrappers.Table(fmt.Sprintf("%s as %s", tableNames[0], do.TableNameMessageLog))
	}

	bizQuery := query.Use(m.DB())
	messageLog := bizQuery.MessageLog

	wrappers = wrappers.Where(messageLog.SendAt.Gte(req.StartAt))
	wrappers = wrappers.Where(messageLog.SendAt.Lte(req.EndAt))
	wrappers = wrappers.Where(messageLog.NamespaceUID.Eq(namespace.Int64()))

	if req.Status > enum.MessageStatus_MessageStatus_UNKNOWN {
		wrappers = wrappers.Where(messageLog.Status.Eq(int32(req.Status)))
	}
	if req.MessageType > enum.MessageType_MessageType_UNKNOWN {
		wrappers = wrappers.Where(messageLog.Type.Eq(int32(req.MessageType)))
	}
	if pointer.IsNotNil(req.PageRequestBo) {
		var total int64
		if err := wrappers.Count(&total).Error; err != nil {
			return nil, err
		}
		req.WithTotal(total)
		wrappers = wrappers.Limit(req.Limit()).Offset(req.Offset())
	}
	var messageLogs []*do.MessageLog
	if err := wrappers.Order(messageLog.CreatedAt.Desc()).Find(&messageLogs).Error; err != nil {
		return nil, err
	}
	messageLogItems := make([]*bo.MessageLogItemBo, 0, len(messageLogs))
	for _, messageLog := range messageLogs {
		messageLogItems = append(messageLogItems, convert.ToMessageLogItemBo(messageLog))
	}
	return bo.NewPageResponseBo(req.PageRequestBo, messageLogItems), nil
}

// UpdateMessageLogStatusIf implements [repository.MessageLog].
func (m *messageLogRepository) UpdateMessageLogStatusIf(ctx context.Context, uid snowflake.ID, oldStatus enum.MessageStatus, newStatus enum.MessageStatus) (bool, error) {
	namespace := contextx.GetNamespaceUID(ctx)
	tableName := do.GenMessageLogTableName(namespace, time.UnixMilli(uid.Time()))
	if _, err := m.Cache().Get(ctx, cache.K(tableName)); err != nil && !do.HasTable(m.DB(), tableName) {
		return false, merr.ErrorNotFound("message log %d not found", uid.Int64())
	}

	bizQuery := query.Use(m.DB().Table(tableName))
	messageLog := bizQuery.MessageLog
	messageLogTable := messageLog.As(tableName)
	wrappers := messageLog.WithContext(ctx)
	wheres := []gen.Condition{
		messageLogTable.UID.Eq(uid.Int64()),
		messageLogTable.NamespaceUID.Eq(namespace.Int64()),
		messageLogTable.Status.Eq(int32(oldStatus)),
	}
	wrappers = wrappers.Where(wheres...)
	result, err := wrappers.Update(messageLogTable.Status, newStatus)
	if err != nil {
		return false, err
	}
	return result.RowsAffected == 1, nil
}

// UpdateMessageLogLastErrorIf implements [repository.MessageLog].
func (m *messageLogRepository) UpdateMessageLogLastErrorIf(ctx context.Context, uid snowflake.ID, oldStatus enum.MessageStatus, lastError string) (bool, error) {
	namespace := contextx.GetNamespaceUID(ctx)
	tableName := do.GenMessageLogTableName(namespace, time.UnixMilli(uid.Time()))
	if _, err := m.Cache().Get(ctx, cache.K(tableName)); err != nil && !do.HasTable(m.DB(), tableName) {
		return false, merr.ErrorNotFound("message log %d not found", uid.Int64())
	}

	bizQuery := query.Use(m.DB().Table(tableName))
	messageLog := bizQuery.MessageLog
	messageLogTable := messageLog.As(tableName)
	wrappers := messageLog.WithContext(ctx)
	wheres := []gen.Condition{
		messageLogTable.UID.Eq(uid.Int64()),
		messageLogTable.NamespaceUID.Eq(namespace.Int64()),
		messageLogTable.Status.Eq(int32(oldStatus)),
	}
	wrappers = wrappers.Where(wheres...)
	columns := []field.AssignExpr{
		messageLogTable.LastError.Value(lastError),
		messageLogTable.Status.Value(int32(enum.MessageStatus_FAILED)),
	}
	result, err := wrappers.UpdateColumnSimple(columns...)
	if err != nil {
		return false, err
	}
	return result.RowsAffected == 1, nil
}

// CreateMessageLog implements [repository.MessageLog].
func (m *messageLogRepository) CreateMessageLog(ctx context.Context, req *bo.MessageLogItemBo) error {
	tableName, err := m.getTableName(ctx, req)
	if err != nil {
		return err
	}
	bizQuery := query.Use(m.DB().Table(tableName))
	messageLog := bizQuery.MessageLog
	mutation := messageLog.WithContext(ctx)
	return mutation.Create(convert.ToMessageLogDO(ctx, req))
}

func (m *messageLogRepository) getTableName(ctx context.Context, req *bo.MessageLogItemBo) (string, error) {
	namespace := contextx.GetNamespaceUID(ctx)
	tableName := do.GenMessageLogTableName(namespace, req.SendAt)

	if _, err := m.Cache().Get(ctx, cache.K(tableName)); err == nil && do.HasTable(m.DB(), tableName) {
		return tableName, nil
	}
	if !do.HasTable(m.DB(), tableName) {
		initModel := &do.MessageLog{}
		oldTableName := initModel.TableName()
		if !do.HasTable(m.DB(), oldTableName) {
			if err := m.DB().Migrator().CreateTable(initModel); err != nil {
				return "", err
			}
		}
		if err := m.DB().Migrator().RenameTable(oldTableName, tableName); err != nil {
			return "", err
		}
	}
	if err := m.Cache().Set(ctx, cache.K(tableName), "", 0); err != nil {
		klog.Context(ctx).Warnw("msg", "set cache failed", "error", err, "tableName", tableName)
	}

	return tableName, nil
}
