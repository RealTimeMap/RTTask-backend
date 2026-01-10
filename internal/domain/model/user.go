package model

import (
	"fmt"
	"rttask/internal/domain/model/rbac"
	"slices"

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
	Avatar         *File `gorm:"type:jsonb;serializer:json"`
}

func (u *User) FullName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}
func (u *User) GetPermissions() []rbac.Permission {
	seen := make(map[rbac.Permission]struct{})
	result := make([]rbac.Permission, 0)

	for _, role := range u.Roles {
		if !role.IsActive {
			continue
		}
		for _, p := range role.Permissions {
			if _, exists := seen[p]; !exists {
				seen[p] = struct{}{}
				result = append(result, p)
			}
		}
	}
	return result
}

func (u *User) Can(p rbac.Permission) bool {
	for _, role := range u.Roles {
		if !role.IsActive {
			continue
		}
		if slices.Contains(role.Permissions, p) {
			return true
		}
	}
	return false
}

func (u *User) CanAll(permissions ...rbac.Permission) bool {
	for _, p := range permissions {
		if !u.Can(p) {
			return false
		}
	}
	return true
}

func (u *User) CanAny(permissions ...rbac.Permission) bool {
	for _, p := range permissions {
		if u.Can(p) {
			return true
		}
	}
	return false
}

func (u *User) HasRole(roleName string) bool {
	for _, role := range u.Roles {
		if role.Name == roleName && role.IsActive {
			return true
		}
	}
	return false
}
