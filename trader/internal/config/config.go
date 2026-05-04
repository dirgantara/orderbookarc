package config

import (
	"os"
)

type Config struct {
	ServerPort         string
	DatabaseURL        string
	JWTSecret          string
	OrderbookURL       string
	OrderbookAPIKeyID  string
	OrderbookAPIKey    string
	OrderbookAPISecret string
}

func Load() *Config {
	return &Config{
		ServerPort:         getEnv("SERVER_PORT", "9000"),
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://orderbook:password@postgres:5432/orderbook?sslmode=disable"),
		JWTSecret:          getEnv("JWT_SECRET", "trader-jwt-secret-change-me"),
		OrderbookURL:       getEnv("ORDERBOOK_URL", "http://server:8080"),
		OrderbookAPIKeyID:  getEnv("ORDERBOOK_API_KEY_ID", ""),
		OrderbookAPIKey:    getEnv("ORDERBOOK_API_KEY", ""),
		OrderbookAPISecret: getEnv("ORDERBOOK_API_SECRET", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
