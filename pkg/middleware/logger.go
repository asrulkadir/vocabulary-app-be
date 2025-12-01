package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger logs API requests with status, latency, and errors
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path

		// Log with error message if present
		if len(c.Errors) > 0 {
			errorMsg := c.Errors.Last().Error()
			log.Printf("[%d] %s %s %v - Error: %s",
				statusCode,
				method,
				path,
				latency,
				errorMsg,
			)
		} else {
			log.Printf("[%d] %s %s %v",
				statusCode,
				method,
				path,
				latency,
			)
		}
	}
}
