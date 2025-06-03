package service

import (
	"context"

	"github.com/Thanhbinh1905/realtime-chat-v2-go/auth-service/internal/model"
	"github.com/Thanhbinh1905/realtime-chat-v2-go/auth-service/internal/repository"
	"github.com/Thanhbinh1905/realtime-chat-v2-go/auth-service/internal/utils/auth"
	"github.com/Thanhbinh1905/realtime-chat-v2-go/auth-service/internal/utils/errors"
	"github.com/Thanhbinh1905/realtime-chat-v2-go/auth-service/internal/utils/hasher"
	"github.com/Thanhbinh1905/realtime-chat-v2-go/shared/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/go-playground/validator/v10"
)

type Service interface {
	Register(ctx context.Context, input *model.RegisterInput) (*model.AuthResponse, error)
	Login(ctx context.Context, input *model.LoginInput) (*model.AuthResponse, error)
}

type service struct {
	repo       repository.Repository
	tokenMaker auth.TokenMaker
	hasher     hasher.Hasher
}

func NewService(repo repository.Repository, tokenMaker auth.TokenMaker, hasher hasher.Hasher) Service {
	return &service{
		repo:       repo,
		tokenMaker: tokenMaker,
		hasher:     hasher,
	}
}

func (s *service) Register(ctx context.Context, input *model.RegisterInput) (*model.AuthResponse, error) {

	// Validate input
	if err := validator.New().Struct(input); err != nil {
		logger.Log.Error("validation error: ", zap.Error(err))
		return nil, err
	}

	// Check if email exists
	exists, err := s.repo.CheckEmailExists(ctx, input.Email)
	if err != nil {
		logger.Log.Error("error checking email existence", zap.Error(err))
		return nil, err
	}

	if exists {
		logger.Log.Error("email already exists", zap.String("email", input.Email))
		return nil, errors.ErrEmailExists
	}

	hashed, err := s.hasher.Hash(input.Password)
	if err != nil {
		logger.Log.Error("error hashing password", zap.Error(err))
		return nil, errors.ErrInternalServerError
	}

	// Tạo RegisterRequest
	req := &model.RegisterRequest{
		ID:           uuid.New(),
		Email:        input.Email,
		PasswordHash: hashed,
	}

	err = s.repo.Register(ctx, req)
	if err != nil {
		logger.Log.Error("error registering user", zap.Error(err))
		return nil, err
	}

	accessToken, refreshToken, refreshExpiresAt, err := s.tokenMaker.GenerateTokens(req.ID.String(), req.Email)
	if err != nil {
		logger.Log.Error("error generating tokens", zap.Error(err))
		return nil, errors.ErrInternalServerError
	}

	err = s.repo.SaveRefreshToken(ctx, refreshToken, req.ID.String(), refreshExpiresAt)
	if err != nil {
		logger.Log.Error("error saving refresh token", zap.Error(err))
		return nil, err
	}

	return &model.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *service) Login(ctx context.Context, input *model.LoginInput) (*model.AuthResponse, error) {
	if err := validator.New().Struct(input); err != nil {
		logger.Log.Error("validation error:", zap.Error(err))
		return nil, err
	}

	account, err := s.repo.GetAccountByEmail(ctx, input.Email)

	if err != nil {
		logger.Log.Error("account not found:", zap.Error(err))
		return nil, errors.ErrInvalidCredentials // Tránh leak info
	}

	if account == nil {
		logger.Log.Error("account not found for email", zap.String("email", input.Email))
		return nil, errors.ErrInvalidCredentials
	}

	matched := s.hasher.Compare(account.PasswordHash, input.Password)
	if !matched {
		logger.Log.Error("password mismatch:", zap.Error(err))
		return nil, errors.ErrInvalidCredentials
	}

	accessToken, refreshToken, refreshExpiresAt, err := s.tokenMaker.GenerateTokens(account.ID.String(), input.Email)
	if err != nil {
		logger.Log.Error("error generating tokens", zap.Error(err))
		return nil, errors.ErrInternalServerError
	}

	err = s.repo.SaveRefreshToken(ctx, refreshToken, account.ID.String(), refreshExpiresAt)
	if err != nil {
		logger.Log.Error("error saving refresh token", zap.Error(err))
		return nil, err
	}

	return &model.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}
