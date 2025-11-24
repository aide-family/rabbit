package dbimpl

import (
	"context"

	"github.com/aide-family/magicbox/pointer"
	"github.com/aide-family/magicbox/strutil"
	"github.com/bwmarrin/snowflake"

	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/do"
	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/internal/data"
	"github.com/aide-family/rabbit/pkg/merr"
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
func (n *namespaceRepositoryImpl) CreateNamespace(ctx context.Context, req *do.Namespace) error {
	namespaceDO := n.d.MainQuery(ctx).Namespace
	wrappers := namespaceDO.WithContext(ctx)
	return wrappers.Create(req)
}

// UpdateNamespace implements repository.Namespace.
func (n *namespaceRepositoryImpl) UpdateNamespace(ctx context.Context, req *do.Namespace) error {
	namespaceDO := n.d.MainQuery(ctx).Namespace
	wrappers := namespaceDO.WithContext(ctx).Where(namespaceDO.UID.Eq(req.UID.Int64()))
	_, err := wrappers.Updates(req)
	return err
}

// UpdateNamespaceStatus implements repository.Namespace.
func (n *namespaceRepositoryImpl) UpdateNamespaceStatus(ctx context.Context, req *bo.UpdateNamespaceStatusBo) error {
	namespaceDO := n.d.MainQuery(ctx).Namespace
	wrappers := namespaceDO.WithContext(ctx).Where(namespaceDO.UID.Eq(req.UID.Int64()))
	_, err := wrappers.Update(namespaceDO.Status, req.Status)
	return err
}

// DeleteNamespace implements repository.Namespace.
func (n *namespaceRepositoryImpl) DeleteNamespace(ctx context.Context, uid snowflake.ID) error {
	namespaceDO := n.d.MainQuery(ctx).Namespace
	wrappers := namespaceDO.WithContext(ctx).Where(namespaceDO.UID.Eq(uid.Int64()))
	_, err := wrappers.Delete()
	return err
}

// GetNamespace implements repository.Namespace.
func (n *namespaceRepositoryImpl) GetNamespace(ctx context.Context, uid snowflake.ID) (*do.Namespace, error) {
	namespaceDO := n.d.MainQuery(ctx).Namespace
	wrappers := namespaceDO.WithContext(ctx).Where(namespaceDO.UID.Eq(uid.Int64()))
	namespaceDo, err := wrappers.First()
	if err != nil {
		if merr.IsNotFound(err) {
			return nil, merr.ErrorNotFound("namespace %s not found", uid)
		}
		return nil, err
	}
	return namespaceDo, nil
}

// ListNamespace implements repository.Namespace.
func (n *namespaceRepositoryImpl) ListNamespace(ctx context.Context, req *bo.ListNamespaceBo) (*bo.PageResponseBo[*do.Namespace], error) {
	namespaceDO := n.d.MainQuery(ctx).Namespace
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

// SelectNamespace implements repository.Namespace.
func (n *namespaceRepositoryImpl) SelectNamespace(ctx context.Context, req *bo.SelectNamespaceBo) (*bo.SelectNamespaceResult, error) {
	namespaceDO := n.d.MainQuery(ctx).Namespace
	wrappers := namespaceDO.WithContext(ctx)

	if strutil.IsNotEmpty(req.Keyword) {
		wrappers = wrappers.Where(namespaceDO.Name.Like("%" + req.Keyword + "%"))
	}
	if req.Status.Exist() && !req.Status.IsUnknown() {
		wrappers = wrappers.Where(namespaceDO.Status.Eq(req.Status.GetValue()))
	}

	// 获取总数
	total, err := wrappers.Count()
	if err != nil {
		return nil, err
	}

	// 游标分页：如果提供了lastUID，则查询UID小于lastUID的记录
	if req.LastUID > 0 {
		wrappers = wrappers.Where(namespaceDO.UID.Lt(req.LastUID.Int64()))
	}

	// 限制返回数量
	wrappers = wrappers.Limit(int(req.Limit))

	// 按UID倒序排列（snowflake ID按时间生成，与CreatedAt一致）
	namespaces, err := wrappers.Order(namespaceDO.UID.Desc()).Find()
	if err != nil {
		return nil, err
	}

	// 获取最后一个UID，用于下次分页
	var lastUID snowflake.ID
	if len(namespaces) > 0 {
		lastUID = namespaces[len(namespaces)-1].UID
	}

	return &bo.SelectNamespaceResult{
		Items:   namespaces,
		Total:   total,
		LastUID: lastUID,
	}, nil
}

// GetNamespaceByName implements repository.Namespace.
func (n *namespaceRepositoryImpl) GetNamespaceByName(ctx context.Context, name string) (*do.Namespace, error) {
	namespaceDO := n.d.MainQuery(ctx).Namespace
	wrappers := namespaceDO.WithContext(ctx).Where(namespaceDO.Name.Eq(name))
	namespaceDo, err := wrappers.First()
	if err != nil {
		if merr.IsNotFound(err) {
			return nil, merr.ErrorNotFound("namespace %s not found", name)
		}
		return nil, err
	}
	return namespaceDo, nil
}
