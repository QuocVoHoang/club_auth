package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	infracontext "github.com/your-org/go-base/internal/infrastructure/context"
)

type AuthClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   int    `json:"role"`
	jwt.StandardClaims
}

func AuthRequired(jwtService JWT) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, ok := bearerToken(ctx.GetHeader("Authorization"))
		if !ok {
			renderUnauthorized(ctx, "missing or invalid authorization header")
			return
		}

		claims := &AuthClaims{}
		if err := jwtService.Decrypt(token, claims, false); err != nil {
			renderUnauthorized(ctx, "invalid access token")
			return
		}

		userID, err := uuid.Parse(claims.UserID)
		if err != nil || userID == uuid.Nil {
			renderUnauthorized(ctx, "invalid access token")
			return
		}

		infracontext.SetUserID(ctx, userID)
		infracontext.SetToken(ctx, token)
		infracontext.SetEmail(ctx, claims.Email)
		infracontext.SetRole(ctx, claims.Role)

		ctx.Next()
	}
}

func bearerToken(header string) (string, bool) {
	parts := strings.Fields(header)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", false
	}

	return parts[1], true
}

func renderUnauthorized(ctx *gin.Context, message string) {
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"code":    http.StatusUnauthorized,
		"message": message,
	})
}
