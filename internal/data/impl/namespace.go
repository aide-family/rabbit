package impl

import (
	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/internal/data"
	"github.com/aide-family/rabbit/internal/data/impl/dbimpl"
	"github.com/aide-family/rabbit/internal/data/impl/fileimpl"
)

func NewNamespaceRepository(d *data.Data) repository.Namespace {
	newRepo := fileimpl.NewNamespaceRepository
	if d.UseDatabase() {
		newRepo = dbimpl.NewNamespaceRepository
	}
	return newRepo(d)
}
