package mq

import "github.com/google/uuid"

type UserSignUpEvent struct {
	UserID uuid.UUID `json:"id"`
	Email  string    `json:"email"`
}
