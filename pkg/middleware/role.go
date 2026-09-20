package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/your-org/go-base/internal/infrastructure/context"
)

// RequireRoles allows a request to continue only when the authenticated user
// has one of the supplied roles. AuthRequired must run before this middleware.
func RequireRoles(allowedRoles ...int) gin.HandlerFunc {
	allowed := make(map[int]struct{}, len(allowedRoles))
	for _, role := range allowedRoles {
		allowed[role] = struct{}{}
	}

	return func(ctx *gin.Context) {
		role := context.GetRole(ctx.Request.Context())
		if _, ok := allowed[role]; !ok {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    http.StatusForbidden,
				"message": "insufficient permissions",
			})
			return
		}
		ctx.Next()
	}
}
