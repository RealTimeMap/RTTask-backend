package dto

// REQUEST

type LoginRequest struct {
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type RegisterRequest struct {
	Email      string `form:"email" binding:"required"`
	FirstName  string `form:"firstName" binding:"required"`
	LastName   string `form:"lastName" binding:"required"`
	Password   string `form:"password" binding:"required"`
	InviteLink string `form:"inviteLink" binding:"required"`
}

// RESPONSE
