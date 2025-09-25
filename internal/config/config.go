package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv          string
	Port            string
	Database        DatabaseConfig
	JWT             JWTConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	URL      string
}

type JWTConfig struct {
	Secret           string
	AccessTokenExp   time.Duration
	RefreshTokenExp  time.Duration
}

func Load() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}

	accessTokenStr := getEnv("ACCESS_TOKEN_EXP", "15m")
	accessTokenExp, err := time.ParseDuration(accessTokenStr)
	if err != nil {
		return nil, fmt.Errorf("invalid ACCESS_TOKEN_EXP: %w", err)
	}

	refreshTokenStr := getEnv("REFRESH_TOKEN_EXP", "168h")
	refreshTokenExp, err := time.ParseDuration(refreshTokenStr)
	if err != nil {
		return nil, fmt.Errorf("invalid REFRESH_TOKEN_EXP: %w", err)
	}

	cfg := &Config{
		AppEnv: getEnv("APP_ENV", "development"),
		Port:   getEnv("PORT", "8080"),
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     dbPort,
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASS", "postgres"),
			Name:     getEnv("DB_NAME", "go_login"),
			URL:      getEnv("DATABASE_URL", ""),
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", ""),
			AccessTokenExp:  accessTokenExp,
			RefreshTokenExp: refreshTokenExp,
		},
	}

	// Validate required fields
	if cfg.Database.URL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	if cfg.JWT.Secret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
