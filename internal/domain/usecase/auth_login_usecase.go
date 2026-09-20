package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/your-org/go-base/internal/domain/entity"
	"github.com/your-org/go-base/internal/domain/repository"
	"github.com/your-org/go-base/internal/domain/usecase/dto"
	"github.com/your-org/go-base/pkg/http_error"
	middlewarepkg "github.com/your-org/go-base/pkg/middleware"
)

const accessTokenDuration = 24 * time.Hour

type loginUsecase struct {
	userRepo repository.IUserRepo
	jwt      middlewarepkg.JWT
}

func NewLoginUsecase(
	userRepo repository.IUserRepo,
	jwt middlewarepkg.JWT,
) LoginUsecase {
	return &loginUsecase{
		userRepo: userRepo,
		jwt:      jwt,
	}
}

func (uc *loginUsecase) Do(ctx context.Context, req dto.LoginRequest) (*dto.LoginResult, error) {
	if err := validateLoginRequest(req); err != nil {
		return nil, err
	}

	user, err := uc.userRepo.FindByEmail(ctx, normalizeEmail(req.Email))
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, http_error.UnauthorizedError("invalid email or password")
	}

	if user.Status != entity.UserStatusActive {
		return nil, http_error.UnauthorizedError("user is inactive")
	}

	if err := comparePassword(user.Password, req.Password, user.PasswordSalt); err != nil {
		return nil, http_error.UnauthorizedError("invalid email or password")
	}

	if err := uc.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		return nil, fmt.Errorf("update last login: %w", err)
	}

	now := time.Now()
	user.LastLogin = &now

	result, err := buildLoginResult(uc.jwt, *user)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func validateLoginRequest(req dto.LoginRequest) error {
	if normalizeEmail(req.Email) == "" {
		return http_error.BadRequestError("email is required")
	}
	if strings.TrimSpace(req.Password) == "" {
		return http_error.BadRequestError("password is required")
	}
	return nil
}

func buildLoginResult(jwtService middlewarepkg.JWT, user entity.User) (*dto.LoginResult, error) {
	token, err := jwtService.Encrypt(entity.AuthClaims{
		UserID: user.ID.String(),
		Email:  user.Email,
		Role:   user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(accessTokenDuration).Unix(),
			IssuedAt:  time.Now().Unix(),
			Subject:   user.ID.String(),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	return &dto.LoginResult{
		AccessToken: token,
	}, nil
}