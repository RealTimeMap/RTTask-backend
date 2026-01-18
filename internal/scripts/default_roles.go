package scripts

import "rttask/internal/domain/model/rbac"

// DefaultRoles содержит определения стандартных системных ролей
var DefaultRoles = []rbac.Role{
	{
		Name:     "Админ",
		IsSystem: true,
		IsActive: true,
		Permissions: []rbac.Permission{
			// Все права - администратор имеет полный доступ
			rbac.InviteCreate,
			rbac.InviteDelete,
			rbac.InviteList,

			rbac.TaskCreate,
			rbac.TaskDelete,
			rbac.TaskUpdate,
			rbac.TaskView,
			rbac.TaskList,
			rbac.TaskAssign,
			rbac.TaskChangeStatus,

			rbac.RoleCreate,
			rbac.RoleUpdate,
			rbac.RoleDelete,
			rbac.RoleList,
			rbac.RoleAssign,

			rbac.UserView,
			rbac.UserList,
			rbac.UserUpdate,
			rbac.UserDelete,

			rbac.CompanyCreate,
			rbac.CompanyUpdate,
			rbac.CompanyDelete,
			rbac.CompanyView,
			rbac.CompanyList,

			rbac.CommentCreate,
			rbac.CommentView,
			rbac.CommentUpdate,
			rbac.CommentDelete,
		},
	},
	{
		Name:     "Менеджер",
		IsSystem: true,
		IsActive: true,
		Permissions: []rbac.Permission{
			// Управление инвайтами
			rbac.InviteCreate,
			rbac.InviteDelete,
			rbac.InviteList,

			// Полное управление задачами
			rbac.TaskCreate,
			rbac.TaskDelete,
			rbac.TaskUpdate,
			rbac.TaskView,
			rbac.TaskList,
			rbac.TaskAssign,
			rbac.TaskChangeStatus,

			// Просмотр ролей
			rbac.RoleList,

			// Просмотр и частичное управление пользователями
			rbac.UserView,
			rbac.UserList,
			rbac.UserUpdate,

			// Управление компаниями
			rbac.CompanyCreate,
			rbac.CompanyUpdate,
			rbac.CompanyDelete,
			rbac.CompanyView,
			rbac.CompanyList,

			// Управление комментариями
			rbac.CommentCreate,
			rbac.CommentView,
			rbac.CommentUpdate,
			rbac.CommentDelete,
		},
	},
	{
		Name:        "Сотрудник",
		IsSystem:    true,
		IsActive:    true,
		Permissions: []rbac.Permission{},
	},
}
