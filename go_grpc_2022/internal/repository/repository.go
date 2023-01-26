package repository

import (
	"context"

	"gorm.io/gorm"
)

const contextTransactionKey = "tx"

type AbstractRepository interface {
	RunInTransaction(ctx context.Context, fn func(txCtx context.Context) error) error
}

type abstractRepository struct {
	db *gorm.DB
}

func (r abstractRepository) newQuery(ctx context.Context) *gorm.DB {
	db, ok := ctx.Value(contextTransactionKey).(*gorm.DB)
	if !ok || db == nil {
		return r.db.WithContext(ctx)
	}
	return db.WithContext(ctx)
}

func (r abstractRepository) RunInTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(context.WithValue(ctx, contextTransactionKey, tx))
	})
}
