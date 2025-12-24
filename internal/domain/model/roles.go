package model

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name        string                 `gorm:"unique"`
	IsActive    bool                   `gorm:"default:true"`
	Permissions map[string]interface{} `gorm:"serializer:json"`
}
