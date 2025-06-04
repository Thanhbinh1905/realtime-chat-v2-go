package repository

import (
	"context"

	db "github.com/Thanhbinh1905/realtime-chat-v2-go/user-service/internal/db/generated"
	"github.com/google/uuid"
)

type Repository interface {
	// Users
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (db.User, error)
	GetUserByEmail(ctx context.Context, email string) (db.User, error)
	ListUsers(ctx context.Context, arg db.ListUsersParams) ([]db.User, error)

	// Friendships
	SendFriendRequest(ctx context.Context, arg db.SendFriendRequestParams) (db.Friendship, error)
	AcceptFriendRequest(ctx context.Context, arg db.AcceptFriendRequestParams) error
	GetFriends(ctx context.Context, requesterID uuid.UUID) ([]db.User, error)
	RejectFriendRequest(ctx context.Context, arg db.RejectFriendRequestParams) error
}

type repository struct {
	q *db.Queries
}

func NewRepository(q *db.Queries) Repository {
	return &repository{q: q}
}

func (r *repository) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	return r.q.CreateUser(ctx, arg)
}

func (r *repository) GetUserByID(ctx context.Context, id uuid.UUID) (db.User, error) {
	return r.q.GetUserByID(ctx, id)
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	return r.q.GetUserByEmail(ctx, email)
}

func (r *repository) ListUsers(ctx context.Context, arg db.ListUsersParams) ([]db.User, error) {
	return r.q.ListUsers(ctx, arg)
}

// Friendships
func (r *repository) SendFriendRequest(ctx context.Context, arg db.SendFriendRequestParams) (db.Friendship, error) {
	return r.q.SendFriendRequest(ctx, arg)
}

func (r *repository) AcceptFriendRequest(ctx context.Context, arg db.AcceptFriendRequestParams) error {
	return r.q.AcceptFriendRequest(ctx, arg)
}

func (r *repository) RejectFriendRequest(ctx context.Context, arg db.RejectFriendRequestParams) error {
	return r.q.RejectFriendRequest(ctx, arg)
}

func (r *repository) GetFriends(ctx context.Context, requesterID uuid.UUID) ([]db.User, error) {
	return r.q.GetFriends(ctx, requesterID)
}
