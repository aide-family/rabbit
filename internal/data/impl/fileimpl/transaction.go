package fileimpl

import (
	"context"

	"github.com/aide-family/rabbit/internal/biz/repository"
	"github.com/aide-family/rabbit/internal/data"
)

func NewTransactionRepository(d *data.Data) repository.Transaction {
	return &transactionRepositoryImpl{
		d: d,
	}
}

type transactionRepositoryImpl struct {
	d *data.Data
}

func (t *transactionRepositoryImpl) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}
