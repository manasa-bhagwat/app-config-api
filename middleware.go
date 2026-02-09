package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger logs every incoming HTTP request
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// After request completes
		duration := time.Since(start)

		log.Printf(
			"%s %s | %d | %v",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			duration,
		)
	}
}
