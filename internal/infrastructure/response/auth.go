package response

import (
	"time"

	"github.com/google/uuid"
	"github.com/your-org/go-base/internal/domain/usecase/dto"
)

type AuthUser struct {
	ID        uuid.UUID  `json:"id"`
	Email     string     `json:"email"`
	Phone     *string    `json:"phone,omitempty"`
	FullName  string     `json:"full_name"`
	Role      int        `json:"role"`
	Avatar    *string    `json:"avatar,omitempty"`
	Birthdate *time.Time `json:"birthdate,omitempty"`
	Status    string     `json:"status"`
	LastLogin *time.Time `json:"last_login,omitempty"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

func NewAuthUser(result dto.UserResult) AuthUser {
	return AuthUser{
		ID:        result.ID,
		Email:     result.Email,
		Phone:     result.Phone,
		FullName:  result.FullName,
		Role:      result.Role,
		Avatar:    result.Avatar,
		Birthdate: result.Birthdate,
		Status:    result.Status,
		LastLogin: result.LastLogin,
	}
}

func NewLoginResponse(result dto.LoginResult) LoginResponse {
	return LoginResponse{
		AccessToken: result.AccessToken,
	}
}
