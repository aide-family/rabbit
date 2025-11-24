package impl

import (
	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/internal/data"
	"github.com/aide-family/rabbit/internal/data/impl/dbimpl"
	"github.com/aide-family/rabbit/internal/data/impl/fileimpl"
)

func NewTemplateRepository(d *data.Data) repository.Template {
	newRepo := fileimpl.NewTemplateRepository
	if d.UseDatabase() {
		newRepo = dbimpl.NewTemplateRepository
	}
	return &templateRepositoryImpl{
		Template: newRepo(d),
	}
}

type templateRepositoryImpl struct {
	repository.Template
}
