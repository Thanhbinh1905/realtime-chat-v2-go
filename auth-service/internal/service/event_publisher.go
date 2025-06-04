package service

import (
	"context"

	"github.com/google/uuid"
)

type UserCreatedEvent struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	CreateAt string    `json:"created_at"`
}

type EventPublisher interface {
	PublishUserCreated(ctx context.Context, event *UserCreatedEvent) error
}
