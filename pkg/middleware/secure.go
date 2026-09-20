package middleware

import (
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
)

// Secure return the middleware instance
func Secure() gin.HandlerFunc {
	return secure.New(secure.Config{
		FrameDeny:             true,
		ContentSecurityPolicy: "frame-ancestors 'none'",
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
	})
}
