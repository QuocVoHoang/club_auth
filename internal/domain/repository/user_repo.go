package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/your-org/go-base/internal/domain/entity"
)

type IUserRepo interface {
	IBaseRepo[entity.User]

	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	UpdateLastLogin(ctx context.Context, userID uuid.UUID) error
}
