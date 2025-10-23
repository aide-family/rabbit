package impl

import (
	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/internal/data"
)

func NewHealthRepository(d *data.Data) repository.Health {
	return &healthRepositoryImpl{
		d: d,
	}
}

type healthRepositoryImpl struct {
	d *data.Data
}

// Readiness implements repository.Health.
func (h *healthRepositoryImpl) Readiness() error {
	panic("unimplemented")
}
