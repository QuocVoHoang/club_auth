package dto

import (
	"time"

	"github.com/google/uuid"
)

type RegisterRequest struct {
	Email     string
	Phone     *string
	FullName  string
	Role      int
	Password  string
	Avatar    *string
	Birthdate *string
}

type LoginRequest struct {
	Email    string
	Password string
}

type UpdateCurrentUserRequest struct {
	UserID    uuid.UUID
	FullName  *string
	Phone     *string
	Avatar    *string
	Birthdate *string
}

type GetCurrentUserRequest struct {
	UserID uuid.UUID
}

type UserResult struct {
	ID        uuid.UUID
	Email     string
	Phone     *string
	FullName  string
	Role      int
	Avatar    *string
	Birthdate *time.Time
	Status    string
	LastLogin *time.Time
}

type RegisterResult = UserResult

type GetCurrentUserResult = UserResult

type UpdateCurrentUserResult = UserResult

type LoginResult struct {
	AccessToken string
}
