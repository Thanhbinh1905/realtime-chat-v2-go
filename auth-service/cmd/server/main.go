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
	"github.com/Thanhbinh1905/realtime-chat-v2-go/auth-service/internal/mq"
	"github.com/Thanhbinh1905/realtime-chat-v2-go/auth-service/internal/repository"
	"github.com/Thanhbinh1905/realtime-chat-v2-go/auth-service/internal/service"
	"github.com/Thanhbinh1905/realtime-chat-v2-go/auth-service/internal/utils/auth"
	"github.com/Thanhbinh1905/realtime-chat-v2-go/auth-service/internal/utils/hasher"
	"github.com/Thanhbinh1905/realtime-chat-v2-go/shared/logger"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Init logger
	logger.InitLogger(true)

	// Load config
	cfg := config.LoadConfig()

	// Connect DB
	if err := db.Connect(cfg.DatabaseURL); err != nil {
		logger.Log.Fatal("failed to connect database", zap.Error(err))
	}
	defer db.Close()

	// Init RabbitMQ
	publisher, err := mq.NewRabbitPublisher(cfg.RAbbitMQURL)
	if err != nil {
		logger.Log.Fatal("failed to create RabbitMQ publisher", zap.Error(err))
	}
	defer publisher.Close()

	// Repo, service, handler
	repo := repository.NewRepository(db.Pool)
	svc := service.NewService(repo, auth.NewTokenMaker(cfg.JWTSecret), hasher.NewHasher())
	h := handler.NewAuthServiceServer(svc, publisher)

	// gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Log.Fatal("failed to listen on port 50051", zap.Error(err))
	}
	grpcServer := grpc.NewServer()
	authpb.RegisterAuthServiceServer(grpcServer, h)

	// run gRPC server
	go func() {
		logger.Log.Info("🚀 gRPC server started on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			logger.Log.Fatal("failed to serve gRPC", zap.Error(err))
		}
	}()

	// gRPC gateway (mux)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gwMux := runtime.NewServeMux()
	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err = authpb.RegisterAuthServiceHandlerFromEndpoint(ctx, gwMux, "localhost:50051", dialOpts)
	if err != nil {
		logger.Log.Fatal("failed to register grpc-gateway handler", zap.Error(err))
	}

	// GIN
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Mount grpc-gateway under /api/v1/auth/*
	r.Any("/api/v1/auth/*any", gin.WrapH(gwMux))

	// Run Gin server
	httpServer := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		logger.Log.Info("🚀 Gin HTTP gateway started on :8080, prefix /api/v1/auth")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("Gin HTTP server error", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("🔥 Shutting down servers...")

	grpcServer.GracefulStop()

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Log.Error("Gin HTTP shutdown error", zap.Error(err))
	}

	logger.Log.Info("✅ Servers stopped cleanly")
}
