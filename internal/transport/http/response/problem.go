package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProblemDetail struct {
	Type     string                 `json:"type"`
	Title    string                 `json:"title"`
	Status   int                    `json:"status"`
	Detail   string                 `json:"detail,omitempty"`
	Instance string                 `json:"instance,omitempty"`
	TraceID  string                 `json:"traceId,omitempty"`
	Meta     map[string]interface{} `json:"meta,omitempty"`
}

func NewProblemDetail(status int, title, detail string) *ProblemDetail {
	return &ProblemDetail{
		Type:   getProblemType(status),
		Title:  title,
		Status: status,
		Detail: detail,
	}
}

// WithTraceID добавляет trace ID
func (p *ProblemDetail) WithTraceID(traceID string) *ProblemDetail {
	p.TraceID = traceID
	return p
}

// WithInstance добавляет URI где произошла ошибка
func (p *ProblemDetail) WithInstance(instance string) *ProblemDetail {
	p.Instance = instance
	return p
}

// WithMeta добавляет дополнительные данные
func (p *ProblemDetail) WithMeta(meta map[string]interface{}) *ProblemDetail {
	p.Meta = meta
	return p
}

func (p *ProblemDetail) Send(c *gin.Context) {
	c.Header("Content-Type", "application/problem+json")
	c.JSON(p.Status, p)
}

func getProblemType(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "https://realtimemap.ru/rttask/problems/bad-request"
	case http.StatusUnauthorized:
		return "https://realtimemap.ru/rttask/problems/unauthorized"
	case http.StatusForbidden:
		return "https://realtimemap.ru/rttask/problems/forbidden"
	case http.StatusNotFound:
		return "https://realtimemap.ru/rttask/problems/not-found"
	case http.StatusConflict:
		return "https://realtimemap.ru/rttask/problems/conflict"
	case http.StatusInternalServerError:
		return "https://realtimemap.ru/rttask/problems/internal-error"
	default:
		return "https://realtimemap.ru/rttask/problems/unknown"
	}
}
