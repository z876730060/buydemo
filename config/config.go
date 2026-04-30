package config

import (
	"os"
)

type Config struct {
	DBPath      string
	JWTSecret   string
	ServerPort  string
	AdminUser   string
	AdminPass   string
}

func Load() *Config {
	return &Config{
		DBPath:     getEnv("DB_PATH", "./data/buydemo.db"),
		JWTSecret:  getEnv("JWT_SECRET", "buy-demo-secret-key-change-in-production"),
		ServerPort: getEnv("SERVER_PORT", ":8080"),
		AdminUser:  getEnv("ADMIN_USER", "admin"),
		AdminPass:  getEnv("ADMIN_PASS", "admin123"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
