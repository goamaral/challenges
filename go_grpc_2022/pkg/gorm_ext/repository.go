package gorm_ext

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type contextKey string

const txCtxKey contextKey = "tx"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{db}
}

func (r Repository) NewQuery(ctx context.Context, conds ...clause.Expression) *gorm.DB {
	db, ok := ctx.Value(txCtxKey).(*gorm.DB)
	if !ok || db == nil {
		db = r.db
	}
	return db.WithContext(ctx).Clauses(conds...)
}

func (r Repository) NewTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	return ExtractError(
		r.db.Transaction(func(tx *gorm.DB) error {
			return fn(context.WithValue(ctx, txCtxKey, tx)) // TODO: Use TxContext
		}),
	)
}
