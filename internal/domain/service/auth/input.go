package auth

import "rttask/internal/domain/valueobject"

type RegisterInput struct {
	Email      valueobject.Email
	FirstName  string
	LastName   string
	Password   valueobject.Password
	InviteLink string
}

type LoginInput struct {
	Email    valueobject.Email
	Password valueobject.Password
}
