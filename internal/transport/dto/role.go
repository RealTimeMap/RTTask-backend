package dto

import "rttask/internal/domain/model/rbac"

// REQUEST

type RoleRequest struct {
	Name        string   `json:"name" binding:"required"`
	Permissions []string `json:"permissions" binding:"required"`
}

// RESPONSE

type RoleResponse struct {
	ID          uint     `json:"id"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

func NewRoleResponse(role *rbac.Role) RoleResponse {
	var permissions []string
	for _, permission := range role.Permissions {
		permissions = append(permissions, string(permission))
	}
	return RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Permissions: permissions,
	}
}
