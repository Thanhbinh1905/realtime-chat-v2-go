package handler

import (
	"context"

	authpb "github.com/Thanhbinh1905/realtime-chat-v2-go/auth-service/api/auth/v1"
	"github.com/Thanhbinh1905/realtime-chat-v2-go/auth-service/internal/model"
	"github.com/Thanhbinh1905/realtime-chat-v2-go/auth-service/internal/mq"
	"github.com/Thanhbinh1905/realtime-chat-v2-go/auth-service/internal/service"
	"github.com/Thanhbinh1905/realtime-chat-v2-go/shared/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authHandler struct {
	service service.Service
	authpb.UnimplementedAuthServiceServer
	publisher *mq.RabbitPublisher
}

type AuthHandler interface {
	Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.AuthResponse, error)
	Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.AuthResponse, error)
}

func NewAuthServiceServer(s service.Service, pub *mq.RabbitPublisher) authpb.AuthServiceServer {
	return &authHandler{
		service:   s,
		publisher: pub,
	}
}

func (h *authHandler) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.AuthResponse, error) {
	input := &model.RegisterInput{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}
	resp, err := h.service.Register(ctx, input)
	if err != nil {
		logger.Log.Error("user registration failed", zap.String("email", req.GetEmail()), zap.Error(err))
		return nil, status.Errorf(codes.Internal, "register failed: %v", err)
	}

	// 📨 Publish event to RabbitMQ
	event := mq.UserSignUpEvent{
		UserID: resp.RegisterRequest.ID,
		Email:  resp.RegisterRequest.Email,
	}

	if err := h.publisher.PublishUserSignedUp(event); err != nil {
		logger.Log.Warn("failed to publish signup event", zap.String("user_id", resp.RegisterRequest.ID.String()), zap.Error(err))
		// optional: return error or just log
	}

	logger.Log.Info("user registered successfully", zap.String("email", req.GetEmail()))
	return &authpb.AuthResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, nil
}

func (h *authHandler) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.AuthResponse, error) {
	input := &model.LoginInput{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}
	resp, err := h.service.Login(ctx, input)
	if err != nil {
		logger.Log.Error("user login failed", zap.String("email", req.GetEmail()), zap.Error(err))
		return nil, status.Errorf(codes.Unauthenticated, "login failed: %v", err)
	}
	logger.Log.Info("user logged in successfully", zap.String("email", req.GetEmail()))
	return &authpb.AuthResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, nil
}
