package handler

import (
	"context"

	authpb "github.com/Thanhbinh1905/realtime-chat-v2-go/auth-service/api/auth/v1"
	"github.com/Thanhbinh1905/realtime-chat-v2-go/auth-service/internal/model"
	"github.com/Thanhbinh1905/realtime-chat-v2-go/auth-service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authHandler struct {
	service service.Service
	authpb.UnimplementedAuthServiceServer
}

type AuthHandler interface {
	Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.AuthResponse, error)
	Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.AuthResponse, error)
}

func NewAuthServiceServer(s service.Service) authpb.AuthServiceServer {
	return &authHandler{
		service: s,
	}
}

func (h *authHandler) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.AuthResponse, error) {
	input := &model.RegisterInput{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}
	resp, err := h.service.Register(ctx, input)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "register failed: %v", err)
	}
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
		return nil, status.Errorf(codes.Unauthenticated, "login failed: %v", err)
	}
	return &authpb.AuthResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, nil
}
