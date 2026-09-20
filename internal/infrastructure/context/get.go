package context

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetUserID(ctx *gin.Context) uuid.UUID {
	if userID, ok := ctx.Get(keyUserID.String()); ok {
		if userID, ok := userID.(uuid.UUID); ok {
			return userID
		}
	}
	return uuid.Nil
}

func GetToken(ctx context.Context) string {
	if token, ok := ctx.Value(keyToken).(string); ok {
		return token
	}
	return ""
}

func GetEmail(ctx context.Context) string {
	if email, ok := ctx.Value(keyEmail).(string); ok {
		return email
	}
	return ""
}

func GetRole(ctx context.Context) int {
	if role, ok := ctx.Value(keyRole).(int); ok {
		return role
	}
	return 0
}
