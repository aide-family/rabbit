package impl

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aide-family/magicbox/pointer"
	"github.com/aide-family/magicbox/safety"
	"github.com/bwmarrin/snowflake"
	"gorm.io/gen"
	"gorm.io/gorm"

	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/do"
	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/internal/biz/vobj"
	"github.com/aide-family/rabbit/internal/data"
	"github.com/aide-family/rabbit/pkg/middler"
)

func NewMessageLogRepository(d *data.Data) repository.MessageLog {
	return &messageLogRepositoryImpl{
		d:     d,
		cache: safety.NewMap(make(map[string]struct{})),
	}
}

type messageLogRepositoryImpl struct {
	d *data.Data

	cache *safety.Map[string, struct{}]
}

func (m *messageLogRepositoryImpl) getTableName(ctx context.Context, req *do.MessageLog) (string, error) {
	namespace := middler.GetNamespace(ctx)
	req.WithNamespace(namespace)
	tableName := do.GenMessageLogTableName(namespace, req.SendAt)

	if _, ok := m.cache.Get(tableName); ok {
		return tableName, nil
	}
	if bizDB := m.d.BizDB(namespace); !do.HasTable(bizDB, tableName) {
		initModel := &do.MessageLog{}
		oldTableName := initModel.TableName()
		if !do.HasTable(bizDB, oldTableName) {
			if err := bizDB.Migrator().CreateTable(initModel); err != nil {
				return "", err
			}
		}
		if err := bizDB.Migrator().RenameTable(oldTableName, tableName); err != nil {
			return "", err
		}
	}
	m.cache.Set(tableName, struct{}{})

	return tableName, nil
}

// CreateMessageLog implements repository.MessageLog.
func (m *messageLogRepositoryImpl) CreateMessageLog(ctx context.Context, req *do.MessageLog) error {
	tableName, err := m.getTableName(ctx, req)
	if err != nil {
		return err
	}
	messageLog := m.d.BizQueryWithTable(req.Namespace, tableName).MessageLog
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
	for _, tableName := range tableNames {
		m.cache.Set(tableName, struct{}{})
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

	messageLog := m.d.BizQuery(namespace).MessageLog.As(do.TableNameMessageLog)

	wrappers = wrappers.Where(messageLog.SendAt.Gte(req.StartAt))
	wrappers = wrappers.Where(messageLog.SendAt.Lte(req.EndAt))
	wrappers = wrappers.Where(messageLog.Namespace.Eq(namespace))

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
	tableName := do.GenMessageLogTableName(namespace, time.UnixMilli(uid.Int64()))
	if _, ok := m.cache.Get(tableName); !ok && !do.HasTable(m.d.BizDB(namespace), tableName) {
		return nil, gorm.ErrRecordNotFound
	}

	bizQuery := m.d.BizQueryWithTable(namespace, tableName)
	messageLog := bizQuery.MessageLog
	wrappers := messageLog.WithContext(ctx)
	wheres := []gen.Condition{
		messageLog.UID.Eq(uid.Int64()),
		messageLog.Namespace.Eq(namespace),
	}
	wrappers = wrappers.Where(wheres...)
	return wrappers.First()
}

// UpdateMessageLogStatus implements repository.MessageLog.
func (m *messageLogRepositoryImpl) UpdateMessageLogStatus(ctx context.Context, uid snowflake.ID, status vobj.MessageStatus) error {
	namespace := middler.GetNamespace(ctx)
	tableName := do.GenMessageLogTableName(namespace, time.UnixMilli(uid.Int64()))
	if _, ok := m.cache.Get(tableName); !ok && !do.HasTable(m.d.BizDB(namespace), tableName) {
		return gorm.ErrRecordNotFound
	}

	messageLog := m.d.BizQueryWithTable(namespace, tableName).MessageLog
	wrappers := messageLog.WithContext(ctx)
	wheres := []gen.Condition{
		messageLog.UID.Eq(uid.Int64()),
		messageLog.Namespace.Eq(namespace),
	}
	wrappers = wrappers.Where(wheres...)
	_, err := wrappers.Update(messageLog.Status, status)
	return err
}
