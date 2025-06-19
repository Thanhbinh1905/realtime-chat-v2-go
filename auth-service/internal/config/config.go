package config

import (
	"os"

	"github.com/Thanhbinh1905/realtime-chat-v2-go/shared/logger"
)

type Config struct {
	DatabaseURL string `mapstructure:"DATABASE_URL"`
	JWTSecret   string `mapstructure:"JWT_SECRET"`

	RAbbitMQURL string `mapstructure:"RABBITMQ_URL"`
}

func LoadConfig() *Config {
	dbURL := os.Getenv("DATABASE_URL")
	jwt := os.Getenv("JWT_SECRET")
	rabbitMQURL := os.Getenv("RABBITMQ_URL")

	if dbURL == "" || jwt == "" || rabbitMQURL == "" {
		logger.Log.Fatal("Missing required environment variables")
	}

	return &Config{
		DatabaseURL: dbURL,
		JWTSecret:   jwt,
		RAbbitMQURL: rabbitMQURL,
	}
}
