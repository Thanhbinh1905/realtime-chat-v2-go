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
	// Init logger, config, DB
	logger.InitLogger(true)
	cfg := config.LoadConfig()
	if err := db.Connect(cfg.DatabaseURL); err != nil {
		logger.Log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Init repo, service, handler
	repo := repository.NewRepository(db.Pool)
	svc := service.NewService(repo, auth.NewTokenMaker(cfg.JWTSecret), hasher.NewHasher())
	h := handler.NewAuthServiceServer(svc)

	// gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Log.Fatal("failed to listen on port 50051", zap.Error(err))
	}
	grpcServer := grpc.NewServer()
	authpb.RegisterAuthServiceServer(grpcServer, h)

	go func() {
		logger.Log.Info("gRPC server started on port 50051")
		if err := grpcServer.Serve(lis); err != nil {
			logger.Log.Fatal("failed to serve gRPC", zap.Error(err))
		}
	}()

	// gRPC Gateway
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err = authpb.RegisterAuthServiceHandlerFromEndpoint(ctx, gwmux, "localhost:50051", opts)
	if err != nil {
		logger.Log.Fatal("failed to register handler", zap.Error(err))
	}

	// ⚡ Prefix /api/v1/auth/*
	prefix := "/api/v1/auth"
	mux := http.NewServeMux()
	mux.Handle(prefix+"/", http.StripPrefix(prefix, gwmux))

	// Dummy /health endpoint for warm-up
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// HTTP server
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		logger.Log.Info("HTTP gateway started on port 8080 with prefix " + prefix)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("HTTP gateway error", zap.Error(err))
		}
	}()

	go func() {
		timeout := time.After(30 * time.Second)
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-timeout:
				logger.Log.Error("warm-up timeout: /health not ready")
				return
			case <-ticker.C:
				resp, err := http.Get("http://localhost:8080/health")
				if err == nil && resp.StatusCode == http.StatusOK {
					logger.Log.Info("/health warm-up successful")
					resp.Body.Close()
					return
				}
				if resp != nil {
					resp.Body.Close()
				}
			}
		}
	}()

	// Shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	logger.Log.Info("Shutting down servers...")

	grpcServer.GracefulStop()

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Log.Error("HTTP server shutdown error", zap.Error(err))
	}

	logger.Log.Info("Servers stopped successfully")
}
