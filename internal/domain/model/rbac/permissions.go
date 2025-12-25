package rbac

import (
	"cmp"
	"slices"
)

type Permission string

type PermissionInfo struct {
	Name        string
	Description string
	Group       string
}

const (
	// Инвайты

	InviteCreate Permission = "invite:create"
	InviteDelete Permission = "invite:delete"
	InviteList   Permission = "invite:list"

	// Задачи

	TaskCreate       Permission = "task:create"
	TaskDelete       Permission = "task:delete"
	TaskUpdate       Permission = "task:update"
	TaskAssign       Permission = "task:assign"
	TaskChangeStatus Permission = "task:changeStatus"

	// Роли

	RoleCreate Permission = "role:create"
	RoleDelete Permission = "role:delete"
	RoleAssign Permission = "role:assign"
)

var PermissionsRegister = map[Permission]PermissionInfo{
	// Инвайты
	InviteCreate: {
		Name:        "invite:create",
		Description: "Создание пригласительных ссылок",
		Group:       "Приглашения",
	},
	InviteDelete: {
		Name:        "invite:delete",
		Description: "Удаление пригласительных ссылок",
		Group:       "Приглашения",
	},
	InviteList: {
		Name:        "invite:list",
		Description: "Просмотр списка пригласительных ссылок",
		Group:       "Приглашения",
	},

	// Задачи
	TaskCreate: {
		Name:        "task:create",
		Description: "Создание новых задач",
		Group:       "Задачи",
	},
	TaskDelete: {
		Name:        "task:delete",
		Description: "Удаление задач",
		Group:       "Задачи",
	},
	TaskUpdate: {
		Name:        "task:update",
		Description: "Редактирование задач",
		Group:       "Задачи",
	},
	TaskAssign: {
		Name:        "task:assign",
		Description: "Назначение задач пользователям",
		Group:       "Задачи",
	},
	TaskChangeStatus: {
		Name:        "task:changeStatus",
		Description: "Изменение статуса задач",
		Group:       "Задачи",
	},

	// Роли
	RoleCreate: {
		Name:        "role:create",
		Description: "Создание новых ролей",
		Group:       "Роли",
	},
	RoleDelete: {
		Name:        "role:delete",
		Description: "Удаление ролей",
		Group:       "Роли",
	},
	RoleAssign: {
		Name:        "role:assign",
		Description: "Назначение ролей пользователям",
		Group:       "Роли",
	},
}

// GetAllPermissions получение всех прав в отсортированном ввиде по группам
func GetAllPermissions() []PermissionInfo {
	result := make([]PermissionInfo, 0, len(PermissionsRegister))
	for _, permission := range PermissionsRegister {
		result = append(result, permission)
	}
	slices.SortFunc(result, func(a, b PermissionInfo) int {
		return cmp.Or(
			cmp.Compare(a.Group, b.Group),
			cmp.Compare(a.Name, b.Name),
		)
	})
	return result
}
