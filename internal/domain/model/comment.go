package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Content string `gorm:"type:text"`
	UserID  uint
	User    User `gorm:"foreignKey:UserID"`
	TaskID  uint
	Task    Task `gorm:"foreignKey:TaskID"`
	// Files
}
