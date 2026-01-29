// Package model is the model package for the namespace service.
package model

import "github.com/aide-family/rabbit/pkg/enum"

type NamespaceModel struct {
	ID        uint32            `json:"id" yaml:"id"`
	UID       int64             `json:"uid" yaml:"uid"`
	Name      string            `json:"name" yaml:"name"`
	Metadata  map[string]string `json:"metadata" yaml:"metadata"`
	Status    enum.GlobalStatus `json:"status" yaml:"status"`
	CreatedAt int64             `json:"createdAt" yaml:"createdAt"`
	UpdatedAt int64             `json:"updatedAt" yaml:"updatedAt"`
	DeletedAt int64             `json:"deletedAt" yaml:"deletedAt"`
	Creator   int64             `json:"creator" yaml:"creator"`
}
