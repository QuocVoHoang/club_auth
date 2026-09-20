package usecase

import (
	"context"

	"github.com/your-org/go-base/internal/domain/usecase/dto"
)

type RegisterUsecase interface {
	Do(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResult, error)
}

type LoginUsecase interface {
	Do(ctx context.Context, req dto.LoginRequest) (*dto.LoginResult, error)
}

type UpdateCurrentUserUsecase interface {
	Do(ctx context.Context, req dto.UpdateCurrentUserRequest) (*dto.UpdateCurrentUserResult, error)
}

type GetCurrentUserUsecase interface {
	Do(ctx context.Context, req dto.GetCurrentUserRequest) (*dto.GetCurrentUserResult, error)
}
