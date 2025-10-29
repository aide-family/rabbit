package impl

import (
	"context"

	"github.com/aide-family/magicbox/pointer"
	"github.com/aide-family/magicbox/strutil"
	"gorm.io/gorm/clause"

	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/do"
	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/internal/data"
)

func NewNamespaceRepository(d *data.Data) repository.Namespace {
	return &namespaceRepositoryImpl{
		d: d,
	}
}

type namespaceRepositoryImpl struct {
	d *data.Data
}

// SaveNamespace implements repository.Namespace.
func (n *namespaceRepositoryImpl) SaveNamespace(ctx context.Context, req *do.Namespace) error {
	namespaceDO := n.d.MainQuery().Namespace
	wrappers := namespaceDO.WithContext(ctx)
	assignmentColumns := []string{
		namespaceDO.Metadata.ColumnName().String(),
		namespaceDO.Status.ColumnName().String(),
		namespaceDO.UpdatedAt.ColumnName().String(),
	}
	onConflict := clause.OnConflict{
		DoUpdates: clause.AssignmentColumns(assignmentColumns),
	}
	return wrappers.Clauses(onConflict).Create(req)
}

// UpdateNamespaceStatus implements repository.Namespace.
func (n *namespaceRepositoryImpl) UpdateNamespaceStatus(ctx context.Context, req *bo.UpdateNamespaceStatusBo) error {
	namespaceDO := n.d.MainQuery().Namespace
	wrappers := namespaceDO.WithContext(ctx).Where(namespaceDO.Name.Eq(req.Name))
	_, err := wrappers.Update(namespaceDO.Status, req.Status)
	return err
}

// DeleteNamespace implements repository.Namespace.
func (n *namespaceRepositoryImpl) DeleteNamespace(ctx context.Context, name string) error {
	namespaceDO := n.d.MainQuery().Namespace
	wrappers := namespaceDO.WithContext(ctx).Where(namespaceDO.Name.Eq(name))
	_, err := wrappers.Delete()
	return err
}

// GetNamespace implements repository.Namespace.
func (n *namespaceRepositoryImpl) GetNamespace(ctx context.Context, name string) (*do.Namespace, error) {
	namespaceDO := n.d.MainQuery().Namespace
	wrappers := namespaceDO.WithContext(ctx).Where(namespaceDO.Name.Eq(name))
	return wrappers.First()
}

// ListNamespace implements repository.Namespace.
func (n *namespaceRepositoryImpl) ListNamespace(ctx context.Context, req *bo.ListNamespaceBo) (*bo.PageResponseBo[*do.Namespace], error) {
	namespaceDO := n.d.MainQuery().Namespace
	wrappers := namespaceDO.WithContext(ctx)
	if strutil.IsNotEmpty(req.Keyword) {
		wrappers = wrappers.Where(namespaceDO.Name.Like("%" + req.Keyword + "%"))
	}
	if req.Status.Exist() && !req.Status.IsUnknown() {
		wrappers = wrappers.Where(namespaceDO.Status.Eq(req.Status.GetValue()))
	}
	if pointer.IsNotNil(req.PageRequestBo) {
		total, err := wrappers.Count()
		if err != nil {
			return nil, err
		}
		req.WithTotal(total)
		wrappers = wrappers.Limit(req.Limit()).Offset(req.Offset())
	}
	namespaces, err := wrappers.Order(namespaceDO.CreatedAt.Desc()).Find()
	if err != nil {
		return nil, err
	}
	return bo.NewPageResponseBo(req.PageRequestBo, namespaces), nil
}
