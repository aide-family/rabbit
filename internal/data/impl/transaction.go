package impl

import (
	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/internal/data"
	"github.com/aide-family/rabbit/internal/data/impl/dbimpl"
	"github.com/aide-family/rabbit/internal/data/impl/fileimpl"
)

func NewTransactionRepository(d *data.Data) repository.Transaction {
	newRepo := fileimpl.NewTransactionRepository
	if d.UseDatabase() {
		newRepo = dbimpl.NewTransactionRepository
	}
	return newRepo(d)
}
