package impl

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aide-family/magicbox/pointer"
	"github.com/aide-family/magicbox/strutil"
	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/do"
	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/internal/data"
	"github.com/aide-family/rabbit/pkg/middler"
	"github.com/google/uuid"
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
	req.UID = uuid.New().String()
	req.WithNamespace(namespace).WithCreator(ctx)
	return wrappers.Create(req)
}

// ListMessageLog implements repository.MessageLog.
func (m *messageLogRepositoryImpl) ListMessageLog(ctx context.Context, req *bo.ListMessageLogBo) (*bo.PageResponseBo[*do.MessageLog], error) {
	namespace := middler.GetNamespace(ctx)
	messageLog := m.d.BizQuery(namespace).MessageLog

	if req.StartAt.IsZero() {
		req.StartAt = time.Now().AddDate(0, 0, -7)
	}
	if req.EndAt.IsZero() {
		req.EndAt = time.Now()
	}
	bizDB := m.d.BizDB(namespace)

	tableNames := do.GenMessageLogTableNames(bizDB, namespace, req.StartAt, req.EndAt)
	tables := make([]any, 0, len(tableNames))
	unionAllSQL := make([]string, 0, len(tableNames))
	for _, tableName := range tableNames {
		tables = append(tables, bizDB.Table(tableName))
		unionAllSQL = append(unionAllSQL, "?")
	}
	wrappers := bizDB.Table(fmt.Sprintf("(%s)", strings.Join(unionAllSQL, " UNION ALL ")), tables...).WithContext(ctx)
	wrappers = wrappers.Where(messageLog.SendAt.Gte(req.StartAt))
	wrappers = wrappers.Where(messageLog.SendAt.Lte(req.EndAt))
	wrappers = wrappers.Where(messageLog.Namespace.Eq(namespace))

	if strutil.IsNotEmpty(req.Keyword) {
		wrappers = wrappers.Where(messageLog.Message.Like("%" + req.Keyword + "%"))
	}
	if req.Status.Exist() {
		wrappers = wrappers.Where(messageLog.Status.Eq(req.Status.GetValue()))
	}
	if req.Type.Exist() {
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
