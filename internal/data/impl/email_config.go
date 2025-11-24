package impl

import (
	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/internal/data"
	"github.com/aide-family/rabbit/internal/data/impl/dbimpl"
	"github.com/aide-family/rabbit/internal/data/impl/fileimpl"
)

func NewEmailConfigRepository(d *data.Data) repository.EmailConfig {
	newRepo := fileimpl.NewEmailConfigRepository
	if d.UseDatabase() {
		newRepo = dbimpl.NewEmailConfigRepository
	}
	return &emailConfigRepositoryImpl{EmailConfig: newRepo(d)}
}

type emailConfigRepositoryImpl struct {
	repository.EmailConfig
}
