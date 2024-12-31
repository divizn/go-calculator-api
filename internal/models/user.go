package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	UserRole     string    `json:"user_role"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type RegisterUserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
	UserRole string `json:"user_role" validate:"required,oneof=admin user"`
}

type LoginUserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
