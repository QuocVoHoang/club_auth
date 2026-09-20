package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Recovery middleware to handle panics and recover gracefully
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var errMsg string
				switch e := err.(type) {
				case string:
					errMsg = e
				case error:
					errMsg = e.Error()
				default:
					errMsg = "unknown error"
				}
				log.Printf("panic recovered: %s", errMsg)
				c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
				c.Abort()
			}
		}()
		c.Next()
	}
}
