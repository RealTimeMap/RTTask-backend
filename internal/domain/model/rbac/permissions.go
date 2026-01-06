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
	TaskView         Permission = "task:view"
	TaskList         Permission = "task:list"
	TaskAssign       Permission = "task:assign"
	TaskChangeStatus Permission = "task:changeStatus"

	// Роли

	RoleCreate Permission = "role:create"
	RoleUpdate Permission = "role:update"
	RoleDelete Permission = "role:delete"
	RoleList   Permission = "role:list"
	RoleAssign Permission = "role:assign"

	// Пользователи

	UserView   Permission = "user:view"
	UserList   Permission = "user:list"
	UserUpdate Permission = "user:update"
	UserDelete Permission = "user:delete"

	// Компании

	CompanyCreate Permission = "company:create"
	CompanyUpdate Permission = "company:update"
	CompanyDelete Permission = "company:delete"
	CompanyView   Permission = "company:view"
	CompanyList   Permission = "company:list"

	// Комментарии

	CommentCreate Permission = "comment:create"
	CommentView   Permission = "comment:view"
	CommentUpdate Permission = "comment:update"
	CommentDelete Permission = "comment:delete"
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
	TaskView: {
		Name:        "task:view",
		Description: "Просмотр задачи",
		Group:       "Задачи",
	},
	TaskList: {
		Name:        "task:list",
		Description: "Просмотр списка задач",
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
	RoleUpdate: {
		Name:        "role:update",
		Description: "Редактирование ролей",
		Group:       "Роли",
	},
	RoleDelete: {
		Name:        "role:delete",
		Description: "Удаление ролей",
		Group:       "Роли",
	},
	RoleList: {
		Name:        "role:list",
		Description: "Просмотр списка ролей",
		Group:       "Роли",
	},
	RoleAssign: {
		Name:        "role:assign",
		Description: "Назначение ролей пользователям",
		Group:       "Роли",
	},

	// Пользователи
	UserView: {
		Name:        "user:view",
		Description: "Просмотр профиля пользователя",
		Group:       "Пользователи",
	},
	UserList: {
		Name:        "user:list",
		Description: "Просмотр списка пользователей",
		Group:       "Пользователи",
	},
	UserUpdate: {
		Name:        "user:update",
		Description: "Редактирование пользователей",
		Group:       "Пользователи",
	},
	UserDelete: {
		Name:        "user:delete",
		Description: "Удаление пользователей",
		Group:       "Пользователи",
	},

	// Компании
	CompanyCreate: {
		Name:        "company:create",
		Description: "Создание компаний",
		Group:       "Компании",
	},
	CompanyUpdate: {
		Name:        "company:update",
		Description: "Редактирование компаний",
		Group:       "Компании",
	},
	CompanyDelete: {
		Name:        "company:delete",
		Description: "Удаление компаний",
		Group:       "Компании",
	},
	CompanyView: {
		Name:        "company:view",
		Description: "Просмотр компании",
		Group:       "Компании",
	},
	CompanyList: {
		Name:        "company:list",
		Description: "Просмотр списка компаний",
		Group:       "Компании",
	},

	// Комментарии
	CommentCreate: {
		Name:        "comment:create",
		Description: "Создание комментариев",
		Group:       "Комментарии",
	},
	CommentView: {
		Name:        "comment:view",
		Description: "Просмотр комментариев",
		Group:       "Комментарии",
	},
	CommentUpdate: {
		Name:        "comment:update",
		Description: "Редактирование комментариев",
		Group:       "Комментарии",
	},
	CommentDelete: {
		Name:        "comment:delete",
		Description: "Удаление комментариев",
		Group:       "Комментарии",
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
