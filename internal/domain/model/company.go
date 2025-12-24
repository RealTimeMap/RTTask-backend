package model

import "gorm.io/gorm"

type Company struct {
	gorm.Model
	Name        string `gorm:"unique"`
	Description string `gorm:"type:text"`

	AvatarID *uint
	Avatar   *File `gorm:"foreignKey:AvatarID"`
}
