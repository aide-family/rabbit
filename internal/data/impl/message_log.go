package impl

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aide-family/magicbox/pointer"
	"github.com/aide-family/magicbox/strutil"
	"github.com/bwmarrin/snowflake"
	"gorm.io/gorm"

	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/do"
	"github.com/aide-family/rabbit/internal/biz/do/query"
	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/internal/biz/vobj"
	"github.com/aide-family/rabbit/internal/data"
	"github.com/aide-family/rabbit/pkg/middler"
)

func NewMessageLogRepository(d *data.Data) repository.MessageLog {
	return &messageLogRepositoryImpl{d: d}
}

type messageLogRepositoryImpl struct {
	d *data.Data
}

// CreateMessageLog implements repository.MessageLog.
func (m *messageLogRepositoryImpl) CreateMessageLog(ctx context.Context, req *do.MessageLog) error {
	namespace := middler.GetNamespace(ctx)
	messageLog := m.d.BizQuery(namespace).MessageLog
	wrappers := messageLog.WithContext(ctx)
	return wrappers.Create(req)
}

// ListMessageLog implements repository.MessageLog.
func (m *messageLogRepositoryImpl) ListMessageLog(ctx context.Context, req *bo.ListMessageLogBo) (*bo.PageResponseBo[*do.MessageLog], error) {
	namespace := middler.GetNamespace(ctx)

	if req.StartAt.IsZero() {
		req.StartAt = time.Now().AddDate(0, 0, -7)
	}
	if req.EndAt.IsZero() {
		req.EndAt = time.Now()
	}
	bizDB := m.d.BizDB(namespace)

	tableNames := do.GenMessageLogTableNames(bizDB, namespace, req.StartAt, req.EndAt)
	if len(tableNames) == 0 {
		return bo.NewPageResponseBo[*do.MessageLog](req.PageRequestBo, nil), nil
	}
	tables := make([]any, 0, len(tableNames))
	unionAllSQL := make([]string, 0, len(tableNames))
	for _, tableName := range tableNames {
		tables = append(tables, bizDB.Table(tableName))
		unionAllSQL = append(unionAllSQL, "?")
	}
	aliasTable := "message_logs"
	wrappers := bizDB.Table(fmt.Sprintf("(%s) as %s", strings.Join(unionAllSQL, " UNION ALL "), aliasTable), tables...).WithContext(ctx)
	messageLog := m.d.BizQuery(namespace).MessageLog.As(aliasTable)

	wrappers = wrappers.Where(messageLog.SendAt.Gte(req.StartAt))
	wrappers = wrappers.Where(messageLog.SendAt.Lte(req.EndAt))
	wrappers = wrappers.Where(messageLog.Namespace.Eq(namespace))

	if strutil.IsNotEmpty(req.Keyword) {
		wrappers = wrappers.Where(messageLog.Message.Like("%" + req.Keyword + "%"))
	}
	if req.Status.Exist() && !req.Status.IsUnknown() {
		wrappers = wrappers.Where(messageLog.Status.Eq(req.Status.GetValue()))
	}
	if req.Type.Exist() && !req.Type.IsUnknown() {
		wrappers = wrappers.Where(messageLog.Type.Eq(req.Type.GetValue()))
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
	return bo.NewPageResponseBo(req.PageRequestBo, messageLogs), nil
}

// GetMessageLog implements repository.MessageLog.
func (m *messageLogRepositoryImpl) GetMessageLog(ctx context.Context, uid snowflake.ID) (*do.MessageLog, error) {
	namespace := middler.GetNamespace(ctx)
	sendAt := time.UnixMilli(uid.Int64())
	tableName := do.GenMessageLogTableName(namespace, sendAt)
	bizDB := m.d.BizDB(namespace)
	if !do.HasTable(bizDB, tableName) {
		return nil, gorm.ErrRecordNotFound
	}

	bizQuery := query.Use(bizDB.Table(tableName))
	messageLog := bizQuery.MessageLog
	wrappers := messageLog.WithContext(ctx)
	wrappers = wrappers.Where(messageLog.UID.Eq(uid.Int64()))
	wrappers = wrappers.Where(messageLog.SendAt.Eq(sendAt))
	wrappers = wrappers.Where(messageLog.Namespace.Eq(namespace))
	return wrappers.First()
}

// UpdateMessageLogStatus implements repository.MessageLog.
func (m *messageLogRepositoryImpl) UpdateMessageLogStatus(ctx context.Context, uid snowflake.ID, status vobj.MessageStatus) error {
	namespace := middler.GetNamespace(ctx)
	bizDB := m.d.BizDB(namespace)
	sendAt := time.UnixMilli(uid.Int64())
	tableName := do.GenMessageLogTableName(namespace, sendAt)
	if !do.HasTable(bizDB, tableName) {
		return gorm.ErrRecordNotFound
	}

	messageLog := query.Use(bizDB.Table(tableName)).MessageLog
	wrappers := messageLog.WithContext(ctx)
	wrappers = wrappers.Where(messageLog.UID.Eq(uid.Int64()))
	wrappers = wrappers.Where(messageLog.SendAt.Eq(sendAt))
	wrappers = wrappers.Where(messageLog.Namespace.Eq(namespace))
	_, err := wrappers.Update(messageLog.Status, status)
	return err
}
