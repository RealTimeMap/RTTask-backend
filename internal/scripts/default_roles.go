package scripts

import "rttask/internal/domain/model/rbac"

// DefaultRoles содержит определения стандартных системных ролей
var DefaultRoles = []rbac.Role{
	{
		Name:     "admin",
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
		Name:     "manager",
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
		Name:     "user",
		IsSystem: true,
		IsActive: true,
		Permissions: []rbac.Permission{
			// Работа с задачами
			rbac.TaskCreate,
			rbac.TaskUpdate,
			rbac.TaskView,
			rbac.TaskList,
			rbac.TaskChangeStatus,

			// Просмотр пользователей
			rbac.UserView,
			rbac.UserList,

			// Просмотр компаний
			rbac.CompanyView,
			rbac.CompanyList,

			// Работа с комментариями
			rbac.CommentCreate,
			rbac.CommentView,
			rbac.CommentUpdate,
			rbac.CommentDelete,
		},
	},
	{
		Name:     "viewer",
		IsSystem: true,
		IsActive: true,
		Permissions: []rbac.Permission{
			// Только просмотр задач
			rbac.TaskView,
			rbac.TaskList,

			// Просмотр пользователей
			rbac.UserView,
			rbac.UserList,

			// Просмотр компаний
			rbac.CompanyView,
			rbac.CompanyList,

			// Просмотр комментариев
			rbac.CommentView,
		},
	},
}
