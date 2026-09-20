package usecase

import (
	"net/mail"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/your-org/go-base/internal/domain/entity"
	"github.com/your-org/go-base/internal/domain/usecase/dto"
)

const ()

func buildUserResult(user entity.User) dto.UserResult {
	return dto.UserResult{
		ID:        user.ID,
		Email:     user.Email,
		Phone:     user.Phone,
		FullName:  user.FullName,
		Role:      user.Role,
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

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func comparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
