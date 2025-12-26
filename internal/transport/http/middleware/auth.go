package middleware

import (
	domainerrors "rttask/internal/domain/errors"
	"rttask/internal/infrastructure/security"
	"rttask/internal/transport/http/response"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	authorizationHeader = "Authorization"
	UserIDKey           = "userID"
	UserEmailKey        = "userEmail"
)

func AuthMiddleware(
	manager security.JWTManager,
	logger *zap.Logger,
	mapper *response.ErrorMapper,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Проверки на валидацию токена, типа и его наличия
		authHeader := c.GetHeader(authorizationHeader)
		if authHeader == "" {
			err := domainerrors.NewUnauthorizedError("authorization header is required")
			logger.Warn("missing authorization header",
				zap.String("traceID", response.GetTraceID(c)),
				zap.String("path", c.Request.URL.Path),
			)
			problem := mapper.MapError(c, err)
			problem.Send(c)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			err := domainerrors.NewUnauthorizedError("invalid authorization header format")
			logger.Warn("invalid auth header format",
				zap.String("traceID", response.GetTraceID(c)),
			)
			problem := mapper.MapError(c, err)
			problem.Send(c)
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := manager.ValidateToken(token)
		if err != nil {
			domainErr := domainerrors.NewUnauthorizedError("invalid or expired token")
			logger.Warn("token validation failed",
				zap.String("traceID", response.GetTraceID(c)),
				zap.Error(err),
			)
			problem := mapper.MapError(c, domainErr)
			problem.Send(c)
			c.Abort()
			return
		}

		if claims.Type != security.AccessToken {
			err := domainerrors.NewUnauthorizedError("invalid token type")
			logger.Warn("wrong token type",
				zap.String("traceID", response.GetTraceID(c)),
				zap.String("tokenType", string(claims.Type)),
			)
			problem := mapper.MapError(c, err)
			problem.Send(c)
			c.Abort()
			return
		}

		// Все проверки прошли

		c.Set(UserIDKey, claims.UserID)
		c.Set(UserEmailKey, claims.Email)

		logger.Debug("request authenticated",
			zap.String("traceID", response.GetTraceID(c)),
			zap.Uint("userID", claims.UserID),
			zap.String("email", claims.Email),
		)

		c.Next()
	}
}
