package model

import "gorm.io/gorm"

type InviteLink struct {
	gorm.Model
	Token string `gorm:"type:varchar(255);not null;unique;index"`
}
