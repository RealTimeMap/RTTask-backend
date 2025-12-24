package model

import "gorm.io/gorm"

type TaskStatus struct {
	gorm.Model
	Name     string `gorm:"unique"`
	IsActive bool   `gorm:"default:true"`
}
