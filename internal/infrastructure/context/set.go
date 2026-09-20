package context

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func setContextKey(ctx *gin.Context, k key, value interface{}) {
	ctx.Set(k.String(), value)
}

func SetUserID(ctx *gin.Context, userID uuid.UUID) {
	setContextKey(ctx, keyUserID, userID)
}

func SetToken(ctx *gin.Context, token string) {
	childCtx := context.WithValue(ctx.Request.Context(), keyToken, token)
	*ctx.Request = *ctx.Request.WithContext(childCtx)
}

func SetEmail(ctx *gin.Context, email string) {
	childCtx := context.WithValue(ctx.Request.Context(), keyEmail, email)
	*ctx.Request = *ctx.Request.WithContext(childCtx)
}

func SetRole(ctx *gin.Context, role int) {
	childCtx := context.WithValue(ctx.Request.Context(), keyRole, role)
	*ctx.Request = *ctx.Request.WithContext(childCtx)
}
