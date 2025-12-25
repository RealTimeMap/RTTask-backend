package rbac

import (
	"slices"

	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name        string       `gorm:"unique"`
	Permissions []Permission `gorm:"serializer:json"`
	IsSystem    bool         `gorm:"default:false"`
	IsActive    bool         `gorm:"default:true"`
}

// HasPermission проверяет есть ли у пользователя права доступа
func (r *Role) HasPermission(p Permission) bool {
	return slices.Contains(r.Permissions, p)
}

// HasAnyPermission проверяет есть ли у пользователя хоть какие либо права
func (r *Role) HasAnyPermission(permissions ...Permission) bool {
	for _, permission := range permissions {
		if r.HasPermission(permission) {
			return true
		}
	}
	return false
}
