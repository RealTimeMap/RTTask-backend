package model

import "gorm.io/gorm"

type File struct {
	gorm.Model

	Name     string `gorm:"not null"`
	Path     string `gorm:"not null"`
	Size     int64
	MimeType string

	UploaderID uint
	Uploader   User `gorm:"foreignKey:UploaderID;constraint:-"`

	EntityType string `gorm:"index"`
	EntityID   uint   `gorm:"index"`
}
