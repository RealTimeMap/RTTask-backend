package valueobject

import (
	domainerrors "rttask/internal/domain/errors"
	"unicode"
)

type Password struct {
	value string
}

func NewPassword(value string) (Password, error) {
	if len(value) < 8 {
		return Password{}, domainerrors.NewValidationError("password must be at least 8 characters")
	}
	hasLetter := false
	hasNumber := false

	for _, char := range value {
		if unicode.IsLetter(char) {
			hasLetter = true
		}
		if unicode.IsNumber(char) {
			hasNumber = true
		}
	}

	if !hasLetter || !hasNumber {
		return Password{}, domainerrors.NewValidationError("password must contain letters and numbers")
	}
	return Password{value}, nil
}

func (p Password) String() string {
	return p.value
}
