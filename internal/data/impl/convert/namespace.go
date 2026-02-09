package convert

import (
	"time"

	namespacev1 "github.com/aide-family/magicbox/domain/namespace/v1"
	"github.com/bwmarrin/snowflake"

	"github.com/aide-family/rabbit/internal/biz/bo"
)

func ToNamespaceItemBo(namespaceModel *namespacev1.NamespaceModel) *bo.NamespaceItemBo {
	return &bo.NamespaceItemBo{
		UID:       snowflake.ParseInt64(namespaceModel.UID),
		Name:      namespaceModel.Name,
		Status:    namespaceModel.Status,
		CreatedAt: time.Unix(namespaceModel.CreatedAt, 0),
		UpdatedAt: time.Unix(namespaceModel.UpdatedAt, 0),
	}
}

func ToNamespaceItemSelectBo(namespaceItemSelect *namespacev1.SelectNamespaceItem) *bo.NamespaceItemSelectBo {
	return &bo.NamespaceItemSelectBo{
		Value:    namespaceItemSelect.Value,
		Label:    namespaceItemSelect.Label,
		Disabled: namespaceItemSelect.Disabled,
		Tooltip:  namespaceItemSelect.Tooltip,
	}
}
