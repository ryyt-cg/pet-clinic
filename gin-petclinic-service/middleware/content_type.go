package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// JSONContentType checks if the Content-Type header is application/json
func JSONContentType() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut {
			if c.GetHeader("Content-Type") != "application/json" {
				c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "Content-Type must be application/json"})
				c.Abort() // Abort the request if content type is not JSON
				return
			}
		}
		c.Next() // Continue to the next handler if content type is valid or method is not POST/PUT
	}
}
