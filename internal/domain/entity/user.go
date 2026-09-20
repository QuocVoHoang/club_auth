package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	UserStatusActive   = "active"
	UserStatusInactive = "inactive"

	UserRoleSuperAdmin = 1
	UserRoleAdmin      = 2
	UserRoleUser       = 3
)

type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	Email        string         `gorm:"type:citext;column:email"`
	Phone        *string        `gorm:"column:phone"`
	FullName     string         `gorm:"column:full_name"`
	Role         int            `gorm:"column:role"`
	Password     string         `gorm:"column:password"`
	PasswordSalt string         `gorm:"column:password_salt"`
	Avatar       *string        `gorm:"column:avatar"`
	Birthdate    *time.Time     `gorm:"column:birthdate"`
	Status       string         `gorm:"column:status"`
	LastLogin    *time.Time     `gorm:"column:last_login"`
	CreatedAt    time.Time      `gorm:"column:created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index"`
	CreatedBy    *uuid.UUID     `gorm:"column:created_by"`
	UpdatedBy    *uuid.UUID     `gorm:"column:updated_by"`
	DeletedBy    *uuid.UUID     `gorm:"column:deleted_by"`
}
