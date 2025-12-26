package valueobject

import (
	"regexp"
	domainerrors "rttask/internal/domain/errors"
	"strings"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

type Email struct {
	value string
}

func NewEmail(value string) (Email, error) {
	normalizedValue := strings.TrimSpace(strings.ToLower(value))

	if !emailRegex.MatchString(normalizedValue) {
		return Email{}, domainerrors.NewValidationError("invalid email address")
	}
	return Email{value: value}, nil
}

func (e *Email) String() string {
	return e.value
}

func (e *Email) Domain() string {
	domainParts := strings.Split(e.value, "@")
	if len(domainParts) == 2 {
		return domainParts[1]
	}
	return ""
}
