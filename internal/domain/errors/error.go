package errors

import (
	"errors"
	"fmt"
)

type ErrorType string

const (
	// Валидации

	ErrorTypeValidation ErrorType = "VALIDATION_ERROR"

	// БД

	ErrorTypeDatabase     ErrorType = "DATABASE_ERROR"
	ErrorTypeNotFound     ErrorType = "NOT_FOUND"
	ErrorTypeAlreadyExist ErrorType = "ALREADY_EXIST_ERROR"

	// Авторизация

	ErrorTypeUnauthorize ErrorType = "UNAUTHORIZED"
	ErrorTypeForbidden   ErrorType = "FORBIDDEN"

	// Серверные

	ErrorTypeInternal ErrorType = "INTERNAL"
)

type DomainError struct {
	Type    ErrorType
	Message string
	Err     error
	Meta    map[string]interface{}
}

func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func (e *DomainError) Unwrap() error {
	return e.Err
}

func (e *DomainError) WithMeta(key string, value interface{}) *DomainError {
	if e.Meta == nil {
		e.Meta = make(map[string]interface{})
	}
	e.Meta[key] = value
	return e
}

// Конструкторы

func NewValidationError(message string) *DomainError {
	return &DomainError{
		Type:    ErrorTypeValidation,
		Message: message,
	}
}

// Конструкторы БД

func NewNotFoundError(entity, identifier string) *DomainError {
	return &DomainError{
		Type:    ErrorTypeNotFound,
		Message: fmt.Sprintf("%s %s not found", entity, identifier),
		Meta:    map[string]interface{}{"entity": entity, "identifier": identifier},
	}
}

func NewAlreadyExistsError(entity, field, value string) *DomainError {
	return &DomainError{
		Type:    ErrorTypeAlreadyExist,
		Message: fmt.Sprintf("%s with %s already exists", entity, field),
		Meta:    map[string]interface{}{"entity": entity, "field": field, "value": value},
	}
}

func NewDatabaseError(message string, err error) *DomainError {
	return &DomainError{
		Type:    ErrorTypeDatabase,
		Message: message,
		Err:     err,
	}
}

// Конструкторы авторизации

func NewUnauthorizedError(message string) *DomainError {
	return &DomainError{
		Type:    ErrorTypeUnauthorize,
		Message: message,
	}
}

func NewForbiddenError(message string) *DomainError {
	return &DomainError{
		Type:    ErrorTypeForbidden,
		Message: message,
	}
}

func NewInternalError(message string, err error) *DomainError {
	return &DomainError{
		Type:    ErrorTypeInternal,
		Message: message,
		Err:     err,
	}
}

func IsDomainError(err error) bool {
	var domainError *DomainError
	return errors.As(err, &domainError)
}

func GetDomainError(err error) *DomainError {
	var domainError *DomainError
	if errors.As(err, &domainError) {
		return domainError
	}
	return nil
}
