package middleware

import (
	"net/http"
	"rttask/internal/transport/http/response"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RecoveryMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()

				traceID := ""
				if id, exists := c.Get("traceID"); exists {
					traceID = id.(string)
				}

				userID := uint(0)
				if id, exists := c.Get("userID"); exists {
					userID = id.(uint)
				}

				logger.Error("panic recovered",
					zap.String("traceID", traceID),
					zap.Uint("userID", userID),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
					zap.Any("error", err),
					zap.String("stack", string(stack)),
				)

				problem := response.NewProblemDetail(
					http.StatusInternalServerError,
					"Internal Server Error",
					"An unexpected error occurred while processing your request",
				).WithTraceID(traceID).WithInstance(c.Request.URL.Path)

				problem.Send(c)
				c.Abort()
			}
		}()

		c.Next()
	}
}
