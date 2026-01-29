package gormimpl

import (
	"encoding/json"

	"gorm.io/gen"

	namespacev1 "github.com/aide-family/rabbit/pkg/domain/namespace/v1"
	"github.com/aide-family/rabbit/pkg/domain/namespace/v1/gormimpl/model"
	"github.com/aide-family/rabbit/pkg/enum"
)

func ConvertNamespaceModel(namespaceDo *model.Namespace) *namespacev1.NamespaceModel {
	return &namespacev1.NamespaceModel{
		Id:        namespaceDo.ID,
		Uid:       namespaceDo.UID.Int64(),
		Name:      namespaceDo.Name,
		Metadata:  namespaceDo.Metadata.Map(),
		Status:    enum.GlobalStatus(namespaceDo.Status),
		CreatedAt: namespaceDo.CreatedAt.Unix(),
		UpdatedAt: namespaceDo.UpdatedAt.Unix(),
		DeletedAt: namespaceDo.DeletedAt.Time.Unix(),
		Creator:   namespaceDo.Creator.Int64(),
	}
}

func ConvertNamespaceItemSelect(namespaceDo *model.Namespace) *namespacev1.NamespaceItemSelect {
	metadata, err := json.Marshal(namespaceDo.Metadata.Map())
	if err != nil {
		return nil
	}
	return &namespacev1.NamespaceItemSelect{
		Value:    namespaceDo.UID.Int64(),
		Label:    namespaceDo.Name,
		Disabled: namespaceDo.DeletedAt.Valid || namespaceDo.Status != uint8(enum.GlobalStatus_ENABLED),
		Tooltip:  string(metadata),
	}
}

func convertResultInfo(result *gen.ResultInfo) *namespacev1.ResultInfo {
	var errStr string
	if err := result.Error; err != nil {
		errStr = err.Error()
	}
	return &namespacev1.ResultInfo{
		RowsAffected: result.RowsAffected,
		Error:        errStr,
	}
}
