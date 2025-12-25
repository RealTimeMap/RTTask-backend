package model

import "gorm.io/gorm"

type InviteLink struct {
	gorm.Model
	Token string
}
