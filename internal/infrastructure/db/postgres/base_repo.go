package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type baseRepo[E any] struct {
	db *gorm.DB
}

func newRepo[E any](db *gorm.DB) *baseRepo[E] {
	return &baseRepo[E]{db: db}
}

// Create inserts a new entity. The caller is responsible for generating the ID
// (e.g., uuid.New()) before calling this method, or relying on GORM hooks.
func (r *baseRepo[E]) Create(ctx context.Context, entity *E) error {
	return r.DBForContext(ctx).Create(entity).Error
}

// Get retrieves a single entity by its primary key (UUID id).
// Returns nil if no record is found (gorm.ErrRecordNotFound is swallowed).
func (r *baseRepo[E]) Get(ctx context.Context, id uuid.UUID) (*E, error) {
	var entity E
	err := r.DBForContext(ctx).Where("id = ?", id).First(&entity).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

// List retrieves entities with cursor-based pagination.
// Pass cursor=nil to start from the beginning.
// Pass limit to control page size.
func (r *baseRepo[E]) List(ctx context.Context, cursor *uuid.UUID, limit int) ([]E, error) {
	var entities []E

	if limit < 1 {
		limit = 20
	}

	q := r.DBForContext(ctx).Order("id ASC").Limit(limit + 1)
	if cursor != nil {
		q = q.Where("id > ?", cursor)
	}

	if err := q.Find(&entities).Error; err != nil {
		return nil, err
	}

	return entities, nil
}

// Count returns the total number of non-deleted rows for the entity type.
func (r *baseRepo[E]) Count(ctx context.Context) (int64, error) {
	var total int64
	if err := r.DBForContext(ctx).Model(new(E)).Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

// Update saves all fields of the entity (including zero values).
func (r *baseRepo[E]) Update(ctx context.Context, entity *E) error {
	return r.DBForContext(ctx).Save(entity).Error
}

// Delete performs a soft-delete by setting deleted_at. GORM's DeletedAt
// field on the entity is used automatically.
func (r *baseRepo[E]) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DBForContext(ctx).Where("id = ?", id).Delete(new(E)).Error
}

func (r *baseRepo[E]) DBForContext(ctx context.Context) *gorm.DB {
	if tx := txFromContext(ctx); tx != nil {
		return tx.WithContext(ctx)
	}

	return r.db.WithContext(ctx)
}