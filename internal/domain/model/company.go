package model

import "gorm.io/gorm"

type Company struct {
	gorm.Model
	Name        string `gorm:"unique"`
	Description string `gorm:"type:text"`

	Avatar *File `gorm:"type:jsonb;serializer:json"`
}
