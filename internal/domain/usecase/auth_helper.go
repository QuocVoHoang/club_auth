package usecase

import (
	"crypto/rand"
	"encoding/hex"
	"net/mail"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/your-org/go-base/internal/domain/entity"
	"github.com/your-org/go-base/internal/domain/usecase/dto"
)

const (
	passwordSaltSize = 16
)

func buildUserResult(user entity.User) dto.UserResult {
	return dto.UserResult{
		ID:        user.ID,
		Email:     user.Email,
		Phone:     user.Phone,
		FullName:  user.FullName,
		Role:      user.Role,
		Avatar:    user.Avatar,
		Birthdate: user.Birthdate,
		Status:    user.Status,
		LastLogin: user.LastLogin,
	}
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func normalizeRequired(value string) string {
	return strings.TrimSpace(value)
}

func normalizeOptional(value *string) *string {
	if value == nil {
		return nil
	}

	normalized := strings.TrimSpace(*value)
	if normalized == "" {
		return nil
	}

	return &normalized
}

func validEmail(email string) bool {
	address, err := mail.ParseAddress(email)
	return err == nil && address.Address == email
}

func validRole(role int) bool {
	switch role {
	case entity.UserRoleSuperAdmin, entity.UserRoleAdmin, entity.UserRoleUser:
		return true
	default:
		return false
	}
}

func parseOptionalDate(value *string) (*time.Time, error) {
	normalized := normalizeOptional(value)
	if normalized == nil {
		return nil, nil
	}

	parsed, err := time.Parse(time.DateOnly, *normalized)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}

func generateSalt() (string, error) {
	buffer := make([]byte, passwordSaltSize)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}

	return hex.EncodeToString(buffer), nil
}

func hashPassword(password, salt string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func comparePassword(hashedPassword, password, salt string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password+salt))
}
