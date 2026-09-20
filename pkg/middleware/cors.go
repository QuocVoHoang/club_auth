package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CorsMiddleware return the middleware instance
func CorsMiddleware(allowOriginHosts []string) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: allowOriginHosts,
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Accept",
			"Authorization",
			"Cache-Control",
			"Content-Length",
			"Content-Type",
			"Cookie",
			"Origin",
			"Pragma",
			RequestIDHeader,
			"X-Csrf-Token",
			"X-Requested-With",
		},
		ExposeHeaders:    []string{"Set-Cookie", "Content-Length", "Content-Disposition"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
