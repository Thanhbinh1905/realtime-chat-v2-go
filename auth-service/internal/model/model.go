package model

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type RegisterRequest struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
}

type RegisterInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthResponse struct {
	AccessToken  string
	RefreshToken string
}
type RegisterResponse struct {
	AccessToken     string
	RefreshToken    string
	RegisterRequest RegisterRequest
}
