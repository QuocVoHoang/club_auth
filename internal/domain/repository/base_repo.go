package repository

import (
	"context"

	"github.com/google/uuid"
)

//go:generate mockgen -source=base_repo.go -destination=mock/mock_base_repo.go -package=mock

// IBaseRepo defines generic CRUD operations for any entity type E.
// Concrete repo interfaces embed this to inherit standard operations
// and add domain-specific query methods.
type IBaseRepo[E any] interface {
	Create(ctx context.Context, entity *E) error
	Get(ctx context.Context, id uuid.UUID) (*E, error)
	List(ctx context.Context, cursor *uuid.UUID, limit int) ([]E, error)
	Count(ctx context.Context) (int64, error)
	Update(ctx context.Context, entity *E) error
	Delete(ctx context.Context, id uuid.UUID) error
}