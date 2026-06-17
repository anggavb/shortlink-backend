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
	Id             int        `json:"id,omitempty"`
	Fullname       string     `json:"fullname,omitempty"`
	Email          string     `json:"email,omitempty"`
	FirstName      string     `json:"first_name,omitempty"`
	LastName       string     `json:"last_name,omitempty"`
	Role           string     `json:"role,omitempty"`
	Workplace      string     `json:"workplace,omitempty"`
	IsMember       bool       `json:"is_member,omitempty"`
	IsReceiveEmail bool       `json:"is_receive_email,omitempty"`
	Photo          string     `json:"photo,omitempty"`
	VerifiedAt     *time.Time `json:"verified_at,omitempty"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
}
