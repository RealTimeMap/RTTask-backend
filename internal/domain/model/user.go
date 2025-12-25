package model

import (
	"fmt"
	"rttask/internal/domain/model/rbac"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName      string
	LastName       string
	Email          string
	HashedPassword string
	Roles          []rbac.Role `gorm:"many2many:users_roles;"`
	Companies      []Company   `gorm:"many2many:users_companies;"`
	AvatarID       *uint
	Avatar         *File `gorm:"foreignKey:AvatarID"`
}

func (u *User) FullName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}
