package postgres

import (
	"context"

	"github.com/your-org/go-base/internal/domain/repository"

	"gorm.io/gorm"
)

type transactionContextKey struct{}

type transactionManager struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) repository.TransactionManager {
	return &transactionManager{db: db}
}

func (m *transactionManager) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	if txFromContext(ctx) != nil {
		return fn(ctx)
	}

	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(context.WithValue(ctx, transactionContextKey{}, tx))
	})
}

func txFromContext(ctx context.Context) *gorm.DB {
	tx, _ := ctx.Value(transactionContextKey{}).(*gorm.DB)
	return tx
}
