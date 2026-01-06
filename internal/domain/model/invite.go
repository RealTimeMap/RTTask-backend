package model

import (
	"rttask/internal/domain/model/rbac"

	"gorm.io/gorm"
)

type InviteLink struct {
	gorm.Model
	UserID      uint
	User        User
	Token       string      `gorm:"type:varchar(255);not null;unique;index"`
	Description *string     `gorm:"type:varchar(255);"`
	Roles       []rbac.Role `gorm:"many2many:invite_links_roles;"`
}
