package main

import (
	"authorization/internal/config"
	"authorization/internal/handler"
	"authorization/internal/pkg/logger"
	"authorization/internal/server"
	"authorization/internal/service"
	"authorization/internal/store"
	"authorization/internal/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize logger
	if err := logger.Init(cfg.AppEnv); err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer logger.Sync()

	logger.Info("Starting authorization service", zap.String("env", cfg.AppEnv), zap.String("port", cfg.Port))

	// Initialize database
	db, err := initDatabase(cfg.Database.URL)
	if err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}

	// Initialize repositories
	userRepo := store.NewUserRepository(db)
	tokenRepo := store.NewTokenRepository(db)

	// Initialize JWT manager
	jwtManager := utils.NewJWTManager(cfg.JWT.Secret, cfg.JWT.AccessTokenExp, cfg.JWT.RefreshTokenExp)

	// Initialize services
	authService := service.NewAuthService(userRepo, tokenRepo, jwtManager)
	userService := service.NewUserService(userRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)

	// Setup router
	router := server.NewRouter(authHandler, userHandler, jwtManager)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("HTTP server starting", zap.String("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}

func initDatabase(databaseURL string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	
	// Retry connection up to 10 times with increasing delay
	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
		if err == nil {
			break
		}
		
		logger.Info("Failed to connect to database, retrying...", 
			zap.Int("attempt", i+1), 
			zap.Error(err))
		
		// Wait before retry (exponential backoff)
		waitTime := time.Duration(i+1) * time.Second
		time.Sleep(waitTime)
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after retries: %w", err)
	}

	// Test connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	// Note: Database schema is managed by golang-migrate migrations
	// Run: make migrate-up to apply migrations
	// Migration files are in db/migrations/ directory
	logger.Info("Database connected successfully")
	logger.Info("Run 'make migrate-up' to apply database migrations")

	return db, nil
}
