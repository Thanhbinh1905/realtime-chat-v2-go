package config

import (
	"os"

	"github.com/Thanhbinh1905/realtime-chat-v2-go/shared/logger"
)

type Config struct {
	DatabaseURL string `mapstructure:"DATABASE_URL"`
	JWTSecret   string `mapstructure:"JWT_SECRET"`
}

func LoadConfig() *Config {
	dbURL := os.Getenv("DATABASE_URL")
	jwt := os.Getenv("JWT_SECRET")

	if dbURL == "" || jwt == "" {
		logger.Log.Fatal("Missing required environment variables")
	}

	return &Config{
		DatabaseURL: dbURL,
		JWTSecret:   jwt,
	}
}
