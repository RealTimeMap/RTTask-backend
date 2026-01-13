package dto

import "rttask/internal/domain/model"

// UserResponse - базовый ответ для пользователя (без чувствительных данных)
type UserResponse struct {
	ID       uint   `json:"id"`
	FullName string `json:"fullName"`
	Email    string `json:"email"`
}

func NewUserResponse(user *model.User) UserResponse {
	return UserResponse{
		ID:       user.ID,
		FullName: user.FullName(),
		Email:    user.Email,
	}
}

// UserDetailResponse - полный ответ с avatar и ролями
type UserDetailResponse struct {
	ID        uint          `json:"id"`
	FirstName string        `json:"firstName"`
	LastName  string        `json:"lastName"`
	Email     string        `json:"email"`
	Avatar    *FileResponse `json:"avatar,omitempty"`
	Roles     []string      `json:"roles,omitempty"`
}

func NewUserDetailResponse(user *model.User) UserDetailResponse {
	roles := make([]string, 0, len(user.Roles))
	for _, role := range user.Roles {
		if role.IsActive {
			roles = append(roles, role.Name)
		}
	}

	return UserDetailResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Avatar:    NewFileResponse(user.Avatar),
		Roles:     roles,
	}
}
