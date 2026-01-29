package impl

import (
	"context"
	"time"

	_ "github.com/aide-family/rabbit/pkg/domain/namespace/v1/fileimpl"
	_ "github.com/aide-family/rabbit/pkg/domain/namespace/v1/gormimpl"

	"github.com/bwmarrin/snowflake"

	"github.com/aide-family/rabbit/internal/biz/bo"
	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/internal/conf"
	"github.com/aide-family/rabbit/internal/data"
	"github.com/aide-family/rabbit/pkg/domain"
	namespacev1 "github.com/aide-family/rabbit/pkg/domain/namespace/v1"
	"github.com/aide-family/rabbit/pkg/enum"
	"github.com/aide-family/rabbit/pkg/merr"
)

func NewNamespaceRepository(c *conf.Bootstrap, d *data.Data) (repository.Namespace, error) {
	repoConfig := c.GetNamespaceConfig()
	version := repoConfig.GetVersion()
	driver := repoConfig.GetDriver()
	switch version {
	default:
		factory, ok := domain.GetNamespaceV1Factory(driver)
		if !ok {
			return nil, merr.ErrorInternalServer("namespace repository factory not found")
		}
		repoImpl, close, err := factory(repoConfig)
		if err != nil {
			return nil, err
		}
		d.AppendClose("namespaceRepo", close)
		return &namespaceRepository{repo: repoImpl}, nil
	}
}

type namespaceRepository struct {
	repo namespacev1.Repository
}

// CreateNamespace implements [repository.Namespace].
func (n *namespaceRepository) CreateNamespace(ctx context.Context, req *bo.CreateNamespaceBo) error {
	_, err := n.repo.CreateNamespace(ctx, &namespacev1.CreateNamespaceRequest{
		Name:     req.Name,
		Metadata: req.Metadata,
		Status:   enum.GlobalStatus(req.Status),
	})
	if err != nil {
		return err
	}
	return nil
}

// DeleteNamespace implements [repository.Namespace].
func (n *namespaceRepository) DeleteNamespace(ctx context.Context, uid snowflake.ID) error {
	_, err := n.repo.DeleteNamespace(ctx, &namespacev1.DeleteNamespaceRequest{
		Uid: uid.Int64(),
	})
	if err != nil {
		return err
	}
	return nil
}

// GetNamespace implements [repository.Namespace].
func (n *namespaceRepository) GetNamespace(ctx context.Context, uid snowflake.ID) (*bo.NamespaceItemBo, error) {
	namespaceModel, err := n.repo.GetNamespace(ctx, &namespacev1.GetNamespaceRequest{
		Uid: uid.Int64(),
	})
	if err != nil {
		return nil, err
	}
	return parseNamespaceModel(namespaceModel), nil
}

// GetNamespaceByName implements [repository.Namespace].
func (n *namespaceRepository) GetNamespaceByName(ctx context.Context, name string) (*bo.NamespaceItemBo, error) {
	namespaceModel, err := n.repo.GetNamespaceByName(ctx, &namespacev1.GetNamespaceByNameRequest{
		Name: name,
	})
	if err != nil {
		return nil, err
	}
	return parseNamespaceModel(namespaceModel), nil
}

// ListNamespace implements [repository.Namespace].
func (n *namespaceRepository) ListNamespace(ctx context.Context, req *bo.ListNamespaceBo) (*bo.PageResponseBo[*bo.NamespaceItemBo], error) {
	listNamespaceResponse, err := n.repo.ListNamespace(ctx, &namespacev1.ListNamespaceRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		Keyword:  req.Keyword,
		Status:   enum.GlobalStatus(req.Status),
	})
	if err != nil {
		return nil, err
	}
	items := make([]*bo.NamespaceItemBo, 0, len(listNamespaceResponse.Namespaces))
	for _, namespaceModel := range listNamespaceResponse.Namespaces {
		items = append(items, parseNamespaceModel(namespaceModel))
	}
	req.WithTotal(listNamespaceResponse.Total)
	return bo.NewPageResponseBo(req.PageRequestBo, items), nil
}

// SelectNamespace implements [repository.Namespace].
func (n *namespaceRepository) SelectNamespace(ctx context.Context, req *bo.SelectNamespaceBo) (*bo.SelectNamespaceBoResult, error) {
	selectNamespaceResponse, err := n.repo.SelectNamespace(ctx, &namespacev1.SelectNamespaceRequest{
		Keyword: req.Keyword,
		Limit:   req.Limit,
		LastUID: req.LastUID.Int64(),
		Status:  enum.GlobalStatus(req.Status),
	})
	if err != nil {
		return nil, err
	}
	items := make([]*bo.NamespaceItemSelectBo, 0, len(selectNamespaceResponse.Items))
	for _, namespaceItemSelect := range selectNamespaceResponse.Items {
		items = append(items, parseNamespaceItemSelect(namespaceItemSelect))
	}
	return &bo.SelectNamespaceBoResult{
		Items:   items,
		Total:   selectNamespaceResponse.Total,
		LastUID: snowflake.ParseInt64(selectNamespaceResponse.LastUID),
		HasMore: selectNamespaceResponse.HasMore,
	}, nil
}

// UpdateNamespace implements [repository.Namespace].
func (n *namespaceRepository) UpdateNamespace(ctx context.Context, req *bo.UpdateNamespaceBo) error {
	panic("unimplemented")
}

// UpdateNamespaceStatus implements [repository.Namespace].
func (n *namespaceRepository) UpdateNamespaceStatus(ctx context.Context, req *bo.UpdateNamespaceStatusBo) error {
	panic("unimplemented")
}

func parseNamespaceModel(namespaceModel *namespacev1.NamespaceModel) *bo.NamespaceItemBo {
	return &bo.NamespaceItemBo{
		UID:       snowflake.ParseInt64(namespaceModel.Uid),
		Name:      namespaceModel.Name,
		Metadata:  namespaceModel.Metadata,
		Status:    namespaceModel.Status,
		CreatedAt: time.Unix(namespaceModel.CreatedAt, 0),
		UpdatedAt: time.Unix(namespaceModel.UpdatedAt, 0),
	}
}

func parseNamespaceItemSelect(namespaceItemSelect *namespacev1.NamespaceItemSelect) *bo.NamespaceItemSelectBo {
	return &bo.NamespaceItemSelectBo{
		UID:      snowflake.ParseInt64(namespaceItemSelect.Value),
		Name:     namespaceItemSelect.Label,
		Disabled: namespaceItemSelect.Disabled,
		Tooltip:  namespaceItemSelect.Tooltip,
	}
}
