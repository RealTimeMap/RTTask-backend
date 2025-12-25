package dto

import (
	"cmp"
	"rttask/internal/domain/model/rbac"
	"slices"
)

type PermissionResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PermissionGroupDTO struct {
	Group       string               `json:"group"`
	Permissions []PermissionResponse `json:"permissions"`
}

type PermissionsResponseDTO struct {
	Groups []PermissionGroupDTO `json:"groups"`
}

func NewGroupedPermissions() []PermissionGroupDTO {
	// Собираем по группам
	grouped := make(map[string][]PermissionResponse)

	for _, info := range rbac.PermissionsRegister {
		grouped[info.Group] = append(grouped[info.Group], PermissionResponse{
			Name:        info.Name,
			Description: info.Description,
		})
	}

	// Конвертируем в слайс
	groups := make([]PermissionGroupDTO, 0, len(grouped))
	for group, permissions := range grouped {
		// Сортируем permissions внутри группы
		slices.SortFunc(permissions, func(a, b PermissionResponse) int {
			return cmp.Compare(a.Name, b.Name)
		})

		groups = append(groups, PermissionGroupDTO{
			Group:       group,
			Permissions: permissions,
		})
	}

	// Сортируем группы
	slices.SortFunc(groups, func(a, b PermissionGroupDTO) int {
		return cmp.Compare(a.Group, b.Group)
	})

	return groups
}
