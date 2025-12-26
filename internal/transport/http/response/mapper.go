package response

import (
	"net/http"
	domainerrors "rttask/internal/domain/errors"

	"github.com/gin-gonic/gin"
)

type ErrorMapper struct{}

func NewErrorMapper() *ErrorMapper {
	return &ErrorMapper{}
}

func (m *ErrorMapper) MapError(c *gin.Context, err error) *ProblemDetail {
	traceID := GetTraceID(c)
	instance := c.Request.URL.Path

	domainErr := domainerrors.GetDomainError(err)
	if domainErr == nil {
		return NewProblemDetail(
			http.StatusInternalServerError,
			"Internal Server Error",
			"An unexpected error occurred",
		).WithTraceID(traceID).WithInstance(instance)
	}

	status := m.mapDomainErrorToHTTPStatus(domainErr.Type)
	title := m.getTitle(domainErr.Type)

	problem := NewProblemDetail(status, title, domainErr.Message).
		WithTraceID(traceID).
		WithInstance(instance)

	if domainErr.Meta != nil && len(domainErr.Meta) > 0 {
		problem.WithMeta(domainErr.Meta)
	}

	return problem
}

func (m *ErrorMapper) mapDomainErrorToHTTPStatus(errorType domainerrors.ErrorType) int {
	switch errorType {
	case domainerrors.ErrorTypeValidation:
		return http.StatusBadRequest // 400
	case domainerrors.ErrorTypeNotFound:
		return http.StatusNotFound // 404
	case domainerrors.ErrorTypeAlreadyExist:
		return http.StatusConflict // 409
	case domainerrors.ErrorTypeUnauthorize:
		return http.StatusUnauthorized // 401
	case domainerrors.ErrorTypeForbidden:
		return http.StatusForbidden // 403
	default:
		return http.StatusInternalServerError // 500
	}
}

func (m *ErrorMapper) getTitle(errorType domainerrors.ErrorType) string {
	switch errorType {
	case domainerrors.ErrorTypeValidation:
		return "Validation Error"
	case domainerrors.ErrorTypeNotFound:
		return "Resource Not Found"
	case domainerrors.ErrorTypeAlreadyExist:
		return "Resource Already Exists"
	case domainerrors.ErrorTypeUnauthorize:
		return "Unauthorized"
	case domainerrors.ErrorTypeForbidden:
		return "Forbidden"
	case domainerrors.ErrorTypeInternal:
		return "Internal Server Error"
	default:
		return "Unknown Error"
	}
}

func GetTraceID(c *gin.Context) string {
	if traceID, ok := c.Get("traceID"); ok {
		return traceID.(string)
	}
	return ""
}

func GetUserID(c *gin.Context) uint {
	if userID, exists := c.Get("userID"); exists {
		return userID.(uint)
	}
	return 0
}
