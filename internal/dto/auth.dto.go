package dto

import "time"

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email" example:"anggavb8@gmail.com"`
	Password string `json:"password" binding:"required,min=8" example:"secret123"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"anggavb8@gmail.com"`
	Password string `json:"password" binding:"required" example:"secret123"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email" example:"anggavb8@gmail.com"`
}

type ResetPasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required,min=8" example:"secret123"`
}

type AuthResponse struct {
	Token     string    `json:"token"`
	User      UserLogin `json:"user"`
	ExpiresAt time.Time `json:"expires_at"`
}

type UserLogin struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}
