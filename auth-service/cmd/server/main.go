package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	authpb "github.com/Thanhbinh1905/realtime-chat-v2-go/auth-service/api/auth/v1"
	"github.com/Thanhbinh1905/realtime-chat-v2-go/auth-service/internal/config"
	"github.com/Thanhbinh1905/realtime-chat-v2-go/auth-service/internal/db"
	"github.com/Thanhbinh1905/realtime-chat-v2-go/auth-service/internal/handler"
	"github.com/Thanhbinh1905/realtime-chat-v2-go/auth-service/internal/repository"
	"github.com/Thanhbinh1905/realtime-chat-v2-go/auth-service/internal/service"
	"github.com/Thanhbinh1905/realtime-chat-v2-go/auth-service/internal/utils/auth"
	"github.com/Thanhbinh1905/realtime-chat-v2-go/auth-service/internal/utils/hasher"
	"github.com/Thanhbinh1905/realtime-chat-v2-go/shared/logger"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	// Khởi tạo logger, config, kết nối DB
	logger.InitLogger(true)
	cfg := config.LoadConfig()
	if err := db.Connect(cfg.DatabaseURL); err != nil {
		logger.Log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Tạo repository, service, handler
	repo := repository.NewRepository(db.Pool)
	service := service.NewService(repo, auth.NewTokenMaker(cfg.JWTSecret), hasher.NewHasher())
	handler := handler.NewAuthServiceServer(service)

	// 1. Chạy gRPC server trên port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Log.Fatal("failed to listen on port 50051", zap.Error(err))
	}
	grpcServer := grpc.NewServer()
	authpb.RegisterAuthServiceServer(grpcServer, handler)

	go func() {
		logger.Log.Info("gRPC server started on port 50051")
		if err := grpcServer.Serve(lis); err != nil {
			logger.Log.Fatal("failed to serve gRPC", zap.Error(err))
		}
	}()

	// 2. Tạo gRPC Gateway mux (proxy HTTP → gRPC)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err = authpb.RegisterAuthServiceHandlerFromEndpoint(ctx, gwmux, "localhost:50051", opts)
	if err != nil {
		logger.Log.Fatal("failed to start HTTP gateway", zap.Error(err))
	}

	// 3. Khởi tạo HTTP server
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: gwmux,
	}

	go func() {
		logger.Log.Info("HTTP gateway started on port 8080")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("failed to serve HTTP gateway", zap.Error(err))
		}
	}()

	// 4. Đợi signal để graceful shutdown (Ctrl+C)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	logger.Log.Info("Shutting down servers...")

	// Graceful stop gRPC server
	grpcServer.GracefulStop()

	// Graceful shutdown HTTP server với timeout 5s
	ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelTimeout()
	if err := httpServer.Shutdown(ctxTimeout); err != nil {
		logger.Log.Error("HTTP server shutdown error", zap.Error(err))
	}

	logger.Log.Info("Servers stopped successfully")
}
