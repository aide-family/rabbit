package data

import (
	"context"

	"gorm.io/gorm"
)

type transactionKey struct{}

type TransactionValue struct {
	DB        *gorm.DB
	Namespace string
	IsMain    bool
}

func withTransaction(ctx context.Context, db *gorm.DB, namespace string, isMain bool) context.Context {
	return context.WithValue(ctx, transactionKey{}, TransactionValue{
		DB:        db,
		Namespace: namespace,
		IsMain:    isMain,
	})
}

func WithMainTransaction(ctx context.Context, db *gorm.DB) context.Context {
	return withTransaction(ctx, db, "", true)
}

func WithBizTransaction(ctx context.Context, db *gorm.DB, namespace string) context.Context {
	return withTransaction(ctx, db, namespace, false)
}

func getTransaction(ctx context.Context, namespace string, isMain bool) (TransactionValue, bool) {
	value, ok := ctx.Value(transactionKey{}).(TransactionValue)
	return value, ok && value.Namespace == namespace && value.IsMain == isMain
}

func GetMainTransaction(ctx context.Context) (TransactionValue, bool) {
	return getTransaction(ctx, "", true)
}

func GetBizTransaction(ctx context.Context, namespace string) (TransactionValue, bool) {
	return getTransaction(ctx, namespace, false)
}
