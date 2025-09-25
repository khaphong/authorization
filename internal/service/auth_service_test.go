package service

import (
	"authorization/internal/dto"
	"authorization/internal/model"
	"authorization/internal/store"
	"authorization/internal/utils"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto-migrate the schema
	if err := db.AutoMigrate(&model.User{}, &model.RefreshToken{}); err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func TestAuthService_Register_Success(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	userRepo := store.NewUserRepository(db)
	tokenRepo := store.NewTokenRepository(db)
	jwtManager := utils.NewJWTManager("test-secret", 15*time.Minute, 7*24*time.Hour)
	
	authService := NewAuthService(userRepo, tokenRepo, jwtManager)

	// Test data
	req := &dto.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	// Execute
	resp, err := authService.Register(req)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response, got nil")
	}

	if resp.User.Username != req.Username {
		t.Errorf("Expected username %s, got %s", req.Username, resp.User.Username)
	}

	if resp.User.Email != req.Email {
		t.Errorf("Expected email %s, got %s", req.Email, resp.User.Email)
	}

	if resp.User.ID == "" {
		t.Error("Expected user ID to be set")
	}

	// Verify user was created in database
	var user model.User
	if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		t.Fatalf("User was not created in database: %v", err)
	}

	if user.Email != req.Email {
		t.Errorf("Expected email %s in database, got %s", req.Email, user.Email)
	}
}

func TestAuthService_Register_DuplicateUsername(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	userRepo := store.NewUserRepository(db)
	tokenRepo := store.NewTokenRepository(db)
	jwtManager := utils.NewJWTManager("test-secret", 15*time.Minute, 7*24*time.Hour)
	
	authService := NewAuthService(userRepo, tokenRepo, jwtManager)

	// Create first user
	req1 := &dto.RegisterRequest{
		Username: "testuser",
		Email:    "test1@example.com",
		Password: "password123",
	}
	
	_, err := authService.Register(req1)
	if err != nil {
		t.Fatalf("Failed to create first user: %v", err)
	}

	// Try to create second user with same username
	req2 := &dto.RegisterRequest{
		Username: "testuser", // Same username
		Email:    "test2@example.com",
		Password: "password456",
	}

	// Execute
	_, err = authService.Register(req2)

	// Assert
	if err == nil {
		t.Fatal("Expected error for duplicate username, got nil")
	}

	if !contains(err.Error(), "already exists") {
		t.Errorf("Expected 'already exists' error, got %v", err)
	}
}

func TestAuthService_Login_Success(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	userRepo := store.NewUserRepository(db)
	tokenRepo := store.NewTokenRepository(db)
	jwtManager := utils.NewJWTManager("test-secret", 15*time.Minute, 7*24*time.Hour)
	
	authService := NewAuthService(userRepo, tokenRepo, jwtManager)

	// Create user first
	registerReq := &dto.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	
	_, err := authService.Register(registerReq)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Test login
	loginReq := &dto.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	// Execute
	resp, err := authService.Login(loginReq)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response, got nil")
	}

	if resp.AccessToken == "" {
		t.Error("Expected access token to be set")
	}

	if resp.RefreshToken == "" {
		t.Error("Expected refresh token to be set")
	}

	if resp.User.Username != loginReq.Username {
		t.Errorf("Expected username %s, got %s", loginReq.Username, resp.User.Username)
	}

	// Verify refresh token was created in database
	var tokenCount int64
	db.Model(&model.RefreshToken{}).Where("user_id = ?", resp.User.ID).Count(&tokenCount)
	if tokenCount != 1 {
		t.Errorf("Expected 1 refresh token in database, got %d", tokenCount)
	}
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	userRepo := store.NewUserRepository(db)
	tokenRepo := store.NewTokenRepository(db)
	jwtManager := utils.NewJWTManager("test-secret", 15*time.Minute, 7*24*time.Hour)
	
	authService := NewAuthService(userRepo, tokenRepo, jwtManager)

	// Create user first
	registerReq := &dto.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	
	_, err := authService.Register(registerReq)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Test login with wrong password
	loginReq := &dto.LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}

	// Execute
	_, err = authService.Login(loginReq)

	// Assert
	if err == nil {
		t.Fatal("Expected error for wrong password, got nil")
	}

	if !contains(err.Error(), "Invalid") {
		t.Errorf("Expected invalid credentials error, got %v", err)
	}
}

func TestAuthService_Login_NonexistentUser(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	userRepo := store.NewUserRepository(db)
	tokenRepo := store.NewTokenRepository(db)
	jwtManager := utils.NewJWTManager("test-secret", 15*time.Minute, 7*24*time.Hour)
	
	authService := NewAuthService(userRepo, tokenRepo, jwtManager)

	// Test login with nonexistent user
	loginReq := &dto.LoginRequest{
		Username: "nonexistent",
		Password: "password123",
	}

	// Execute
	_, err := authService.Login(loginReq)

	// Assert
	if err == nil {
		t.Fatal("Expected error for nonexistent user, got nil")
	}

	if !contains(err.Error(), "Invalid") {
		t.Errorf("Expected invalid credentials error, got %v", err)
	}
}

// Helper functions
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || s[0:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
