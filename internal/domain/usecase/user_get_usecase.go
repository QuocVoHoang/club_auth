package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/your-org/go-base/internal/domain/entity"
	"github.com/your-org/go-base/internal/domain/repository"
	"github.com/your-org/go-base/internal/domain/usecase/dto"
	"github.com/your-org/go-base/pkg/http_error"
)

type getCurrentUserUsecase struct {
	userRepo repository.IUserRepo
}

func NewGetCurrentUserUsecase(userRepo repository.IUserRepo) GetCurrentUserUsecase {
	return &getCurrentUserUsecase{
		userRepo: userRepo,
	}
}

func (uc *getCurrentUserUsecase) Do(
	ctx context.Context,
	req dto.GetCurrentUserRequest,
) (*dto.GetCurrentUserResult, error) {
	if req.UserID == uuid.Nil {
		return nil, http_error.UnauthorizedError("missing authenticated user")
	}

	user, err := uc.userRepo.Get(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, http_error.NotFoundError("resource not found")
	}

	if user.Status != entity.UserStatusActive {
		return nil, http_error.UnauthorizedError("user is inactive")
	}

	userResult := buildUserResult(*user)
	return &userResult, nil
}
