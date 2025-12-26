package dto

import "rttask/internal/domain/model"

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
