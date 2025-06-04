package service

import (
	"context"
	"database/sql"

	db "github.com/Thanhbinh1905/realtime-chat-v2-go/user-service/internal/db/generated"
	"github.com/Thanhbinh1905/realtime-chat-v2-go/user-service/internal/model"
	"github.com/Thanhbinh1905/realtime-chat-v2-go/user-service/internal/repository"
	"github.com/google/uuid"
)

type Service interface {
}

type service struct {
	repo repository.Repository
}

func NewSercice(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateUser(ctx context.Context, input model.CreateUserInput) (*db.User, error) {
	arg := db.CreateUserParams{
		ID:       uuid.New(),
		Email:    input.Email,
		Username: input.Username,
		Avatar:   sql.NullString{String: input.Avatar, Valid: true},
	}
	user, err := s.repo.CreateUser(ctx, arg)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *service) SendFriendRequest(ctx context.Context, from, to uuid.UUID) error {
	_, err := s.repo.SendFriendRequest(ctx, db.SendFriendRequestParams{
		ID:          uuid.New(),
		RequesterID: from,
		AddresseeID: to,
	})
	return err
}

func (s *service) AcceptFriendRequest(ctx context.Context, arg db.AcceptFriendRequestParams) error {
	return s.repo.AcceptFriendRequest(ctx, db.AcceptFriendRequestParams{
		RequesterID: arg.RequesterID,
		AddresseeID: arg.AddresseeID,
	})
}

func (s *service) RejectFriendRequest(ctx context.Context, arg db.RejectFriendRequestParams) error {
	return s.repo.RejectFriendRequest(ctx, db.RejectFriendRequestParams{
		RequesterID: arg.RequesterID,
		AddresseeID: arg.AddresseeID,
	})
}

func (s *service) GetFriends(ctx context.Context, userID uuid.UUID) ([]db.User, error) {
	return s.repo.GetFriends(ctx, userID)
}
