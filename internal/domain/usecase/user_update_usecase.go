package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/your-org/go-base/internal/domain/entity"
	"github.com/your-org/go-base/internal/domain/repository"
	"github.com/your-org/go-base/internal/domain/usecase/dto"
	"github.com/your-org/go-base/pkg/http_error"
)

type updateCurrentUserUsecase struct {
	userRepo repository.IUserRepo
}

func NewUpdateCurrentUserUsecase(
	userRepo repository.IUserRepo,
) UpdateCurrentUserUsecase {
	return &updateCurrentUserUsecase{
		userRepo: userRepo,
	}
}

func (uc *updateCurrentUserUsecase) Do(
	ctx context.Context,
	req dto.UpdateCurrentUserRequest,
) (*dto.UpdateCurrentUserResult, error) {
	if err := validateUpdateCurrentUserRequest(req); err != nil {
		return nil, err
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

	applyUpdateCurrentUserRequest(user, req)
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	userResult := buildUserResult(*user)
	return &userResult, nil
}

func validateUpdateCurrentUserRequest(req dto.UpdateCurrentUserRequest) error {
	if req.UserID == uuid.Nil {
		return http_error.UnauthorizedError("missing authenticated user")
	}
	if req.Birthdate != nil {
		if _, err := parseOptionalDate(req.Birthdate); err != nil {
			return http_error.BadRequestError("birthdate must use YYYY-MM-DD format")
		}
	}

	return nil
}

func applyUpdateCurrentUserRequest(user *entity.User, req dto.UpdateCurrentUserRequest) {
	if req.FullName != nil {
		if fullName := normalizeOptional(req.FullName); fullName != nil {
			user.FullName = *fullName
		}
	}
	if req.Phone != nil {
		user.Phone = normalizeOptional(req.Phone)
	}
	if req.Avatar != nil {
		user.Avatar = normalizeOptional(req.Avatar)
	}
	if req.Birthdate != nil {
		birthdate, _ := parseOptionalDate(req.Birthdate)
		user.Birthdate = birthdate
	}
}
