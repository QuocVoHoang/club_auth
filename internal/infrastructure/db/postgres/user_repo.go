package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/go-base/internal/domain/entity"
	"github.com/your-org/go-base/internal/domain/repository"

	"gorm.io/gorm"
)

type userRepo struct {
	*baseRepo[entity.User]
}

func NewUserRepository(db *gorm.DB) repository.IUserRepo {
	return &userRepo{
		baseRepo: newRepo[entity.User](db),
	}
}

var _ repository.IUserRepo = (*userRepo)(nil)

func (r *userRepo) Create(ctx context.Context, user *entity.User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	return r.baseRepo.Create(ctx, user)
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.DBForContext(ctx).
		Where("email = ?", email).
		First(&user).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	return r.DBForContext(ctx).
		Model(&entity.User{}).
		Where("id = ?", userID).
		Update("last_login", time.Now()).
		Error
}
