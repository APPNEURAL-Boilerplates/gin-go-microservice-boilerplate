package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/example/gin-go-microservice-boilerplate/internal/common"
)

const requestIDHeader = "X-Request-Id"

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(requestIDHeader)
		if requestID == "" {
			requestID = common.NewID("req")
		}

		c.Set(common.RequestIDKey, requestID)
		c.Writer.Header().Set(requestIDHeader, requestID)
		c.Next()
	}
}
