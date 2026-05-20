package api

import "time"

type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	UserID       string `json:"user_id" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

type ProtectedResponse struct {
	Message string `json:"message"`
	UserID  string `json:"user_id"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type RegisterResponse struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
