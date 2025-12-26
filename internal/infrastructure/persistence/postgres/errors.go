package postgres

import (
	"errors"
	domainerrors "rttask/internal/domain/errors"

	"gorm.io/gorm"
)

func MapGormError(err error, entity string) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domainerrors.NewNotFoundError(entity, "")
	}
	return domainerrors.NewDatabaseError("postgres error", err)
}
