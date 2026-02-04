// Package gormimpl is the implementation of the gorm repository for the namespace service.
package gormimpl

import (
	"context"
	"errors"

	"github.com/aide-family/magicbox/hello"
	"github.com/aide-family/magicbox/pointer"
	"github.com/aide-family/magicbox/safety"
	"github.com/aide-family/magicbox/strutil"
	"github.com/bwmarrin/snowflake"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"gorm.io/gen/field"
	"gorm.io/gorm"

	"github.com/aide-family/rabbit/pkg/config"
	"github.com/aide-family/rabbit/pkg/connect"
	"github.com/aide-family/rabbit/pkg/contextx"
	domain "github.com/aide-family/rabbit/pkg/domain"
	namespacev1 "github.com/aide-family/rabbit/pkg/domain/namespace/v1"
	"github.com/aide-family/rabbit/pkg/domain/namespace/v1/gormimpl/model"
	"github.com/aide-family/rabbit/pkg/domain/namespace/v1/gormimpl/query"
	"github.com/aide-family/rabbit/pkg/enum"
	"github.com/aide-family/rabbit/pkg/merr"
)

func init() {
	domain.RegisterNamespaceV1Factory(config.DomainConfig_GORM, NewGormRepository)
}

type Field interface {
	Desc() field.Expr
	Asc() field.Expr
}

func NewGormRepository(c *config.DomainConfig) (namespacev1.Repository, func() error, error) {
	ormConfig := &config.ORMConfig{}
	if pointer.IsNotNil(c.GetOptions()) {
		if err := anypb.UnmarshalTo(c.GetOptions(), ormConfig, proto.UnmarshalOptions{Merge: true}); err != nil {
			return nil, nil, merr.ErrorInternalServer("unmarshal orm config failed: %v", err)
		}
	}
	db, close, err := connect.NewDB(ormConfig)
	if err != nil {
		return nil, nil, err
	}
	query.SetDefault(db)
	fields := safety.NewMap(map[namespacev1.Field]Field{
		namespacev1.Field_ID:         query.Namespace.ID,
		namespacev1.Field_UID:        query.Namespace.UID,
		namespacev1.Field_NAME:       query.Namespace.Name,
		namespacev1.Field_METADATA:   query.Namespace.Metadata,
		namespacev1.Field_STATUS:     query.Namespace.Status,
		namespacev1.Field_CREATED_AT: query.Namespace.CreatedAt,
		namespacev1.Field_UPDATED_AT: query.Namespace.UpdatedAt,
		namespacev1.Field_DELETED_AT: query.Namespace.DeletedAt,
		namespacev1.Field_CREATOR:    query.Namespace.Creator,
	})
	node, err := snowflake.NewNode(hello.NodeID())
	if err != nil {
		return nil, nil, err
	}
	return &gormRepository{repoConfig: c, db: db, fields: fields, node: node}, close, nil
}

type gormRepository struct {
	repoConfig *config.DomainConfig
	db         *gorm.DB
	fields     *safety.Map[namespacev1.Field, Field]
	node       *snowflake.Node
}

func (g *gormRepository) getField(name namespacev1.Field) Field {
	field, ok := g.fields.Get(name)
	if !ok {
		return query.Namespace.CreatedAt
	}
	return field
}

// CreateNamespace implements [namespacev1.Repository].
func (g *gormRepository) CreateNamespace(ctx context.Context, req *namespacev1.CreateNamespaceRequest) (*namespacev1.NamespaceModel, error) {
	namespaceDo := &model.Namespace{
		Name:     req.Name,
		Metadata: safety.NewMap(req.Metadata),
		Status:   uint8(req.Status),
	}
	namespaceDo.WithCreator(contextx.GetUserUID(ctx))
	mutation := query.Namespace
	if err := mutation.WithContext(ctx).Create(namespaceDo); err != nil {
		return nil, merr.ErrorInternalServer("create namespace failed: %v", err)
	}
	return g.GetNamespaceByName(ctx, &namespacev1.GetNamespaceByNameRequest{Name: req.Name})
}

// DeleteNamespace implements [namespacev1.Repository].
func (g *gormRepository) DeleteNamespace(ctx context.Context, req *namespacev1.DeleteNamespaceRequest) (*namespacev1.ResultInfo, error) {
	mutation := query.Use(g.db)
	result, err := mutation.Namespace.WithContext(ctx).Where(mutation.Namespace.UID.Eq(req.Uid)).Delete()
	if err != nil {
		return nil, merr.ErrorInternalServer("delete namespace failed: %v", err)
	}
	return convertResultInfo(&result), nil
}

// GetNamespace implements [namespacev1.Repository].
func (g *gormRepository) GetNamespace(ctx context.Context, req *namespacev1.GetNamespaceRequest) (*namespacev1.NamespaceModel, error) {
	mutation := query.Use(g.db)
	queryNamespace, err := mutation.Namespace.WithContext(ctx).Where(mutation.Namespace.UID.Eq(req.Uid)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, merr.ErrorNotFound("namespace %d not found", req.Uid)
		}
		return nil, err
	}
	return ConvertNamespaceModel(queryNamespace), nil
}

// GetNamespaceByName implements [namespacev1.Repository].
func (g *gormRepository) GetNamespaceByName(ctx context.Context, req *namespacev1.GetNamespaceByNameRequest) (*namespacev1.NamespaceModel, error) {
	mutation := query.Use(g.db)
	queryNamespace, err := mutation.Namespace.WithContext(ctx).Where(mutation.Namespace.Name.Eq(req.Name)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, merr.ErrorNotFound("namespace %s not found", req.Name)
		}
		return nil, err
	}
	return ConvertNamespaceModel(queryNamespace), nil
}

// ListNamespace implements [namespacev1.Repository].
func (g *gormRepository) ListNamespace(ctx context.Context, req *namespacev1.ListNamespaceRequest) (*namespacev1.ListNamespaceResponse, error) {
	mutation := query.Namespace
	wrappers := mutation.WithContext(ctx)
	if strutil.IsNotEmpty(req.Keyword) {
		wrappers = wrappers.Where(mutation.Name.Like("%" + req.Keyword + "%"))
	}
	if req.Status > enum.GlobalStatus_GlobalStatus_UNKNOWN {
		wrappers = wrappers.Where(mutation.Status.Eq(uint8(req.Status)))
	}

	fieldExpr := g.getField(req.OrderBy)
	switch req.Order {
	case namespacev1.Order_DESC:
		wrappers = wrappers.Order(fieldExpr.Desc())
	case namespacev1.Order_ASC:
		wrappers = wrappers.Order(fieldExpr.Asc())
	}

	var total int64
	if req.Page > 0 && req.PageSize > 0 {
		wrappers = wrappers.Limit(int(req.PageSize)).Offset(int((req.Page - 1) * req.PageSize))
		var err error
		total, err = wrappers.Count()
		if err != nil {
			return nil, merr.ErrorInternalServer("count namespace failed: %v", err)
		}
	}
	queryNamespaces, err := wrappers.Find()
	if err != nil {
		return nil, merr.ErrorInternalServer("list namespace failed: %v", err)
	}
	namespaces := make([]*namespacev1.NamespaceModel, 0, len(queryNamespaces))
	for _, queryNamespace := range queryNamespaces {
		namespaces = append(namespaces, ConvertNamespaceModel(queryNamespace))
	}
	return &namespacev1.ListNamespaceResponse{
		Namespaces: namespaces,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
	}, nil
}

// SelectNamespace implements [namespacev1.Repository].
func (g *gormRepository) SelectNamespace(ctx context.Context, req *namespacev1.SelectNamespaceRequest) (*namespacev1.SelectNamespaceResponse, error) {
	mutation := query.Namespace
	wrappers := mutation.WithContext(ctx)
	if strutil.IsNotEmpty(req.Keyword) {
		wrappers = wrappers.Where(mutation.Name.Like("%" + req.Keyword + "%"))
	}
	if req.Status > enum.GlobalStatus_GlobalStatus_UNKNOWN {
		wrappers = wrappers.Where(mutation.Status.Eq(uint8(req.Status)))
	}
	wrappers = wrappers.Limit(int(req.Limit))
	switch req.Order {
	case namespacev1.Order_DESC:
		wrappers = wrappers.Order(mutation.UID.Desc())
	case namespacev1.Order_ASC:
		wrappers = wrappers.Order(mutation.UID.Asc())
	}
	if req.LastUID > 0 {
		switch req.Order {
		case namespacev1.Order_DESC:
			wrappers = wrappers.Where(mutation.UID.Lt(req.LastUID))
		case namespacev1.Order_ASC:
			wrappers = wrappers.Where(mutation.UID.Gt(req.LastUID))
		}
	}
	wrappers = wrappers.Select(mutation.UID, mutation.Name, mutation.Status, mutation.DeletedAt)
	queryNamespaces, err := wrappers.Find()
	if err != nil {
		return nil, merr.ErrorInternalServer("select namespace failed: %v", err)
	}
	total, err := wrappers.Count()
	if err != nil {
		return nil, merr.ErrorInternalServer("count namespace failed: %v", err)
	}
	namespaces := make([]*namespacev1.NamespaceItemSelect, 0, len(queryNamespaces))
	for _, queryNamespace := range queryNamespaces {
		namespaces = append(namespaces, ConvertNamespaceItemSelect(queryNamespace))
	}
	return &namespacev1.SelectNamespaceResponse{
		Items:   namespaces,
		Total:   total,
		LastUID: queryNamespaces[len(queryNamespaces)-1].UID.Int64(),
		HasMore: len(queryNamespaces) == int(req.Limit),
	}, nil
}

// UpdateNamespace implements [namespacev1.Repository].
func (g *gormRepository) UpdateNamespace(ctx context.Context, req *namespacev1.UpdateNamespaceRequest) (*namespacev1.ResultInfo, error) {
	metadata := safety.NewMap(req.Metadata)
	mutation := query.Use(g.db)
	result, err := mutation.Namespace.WithContext(ctx).Where(mutation.Namespace.UID.Eq(req.Uid)).UpdateSimple(mutation.Namespace.Name.Value(req.Name), mutation.Namespace.Metadata.Value(metadata))
	if err != nil {
		return nil, merr.ErrorInternalServer("update namespace failed: %v", err)
	}
	return convertResultInfo(&result), nil
}

// UpdateNamespaceStatus implements [namespacev1.Repository].
func (g *gormRepository) UpdateNamespaceStatus(ctx context.Context, req *namespacev1.UpdateNamespaceStatusRequest) (*namespacev1.ResultInfo, error) {
	mutation := query.Use(g.db)
	result, err := mutation.Namespace.WithContext(ctx).Where(mutation.Namespace.UID.Eq(req.Uid)).UpdateSimple(mutation.Namespace.Status.Value(uint8(req.Status)))
	if err != nil {
		return nil, merr.ErrorInternalServer("update namespace status failed: %v", err)
	}
	return convertResultInfo(&result), nil
}
