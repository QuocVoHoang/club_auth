package middleware

import (
	"github.com/gin-gonic/gin"
)

// Headers add secure headers
func Headers(c *gin.Context) {
	c.Writer.Header().Add("Cache-Control", "no-store")

	c.Next()
}
