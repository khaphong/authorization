package service

import (
	"authorization/internal/constants"
	"authorization/internal/dto"
	"authorization/internal/model"
	"authorization/internal/store"
	"authorization/internal/utils"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type AuthService struct {
	userRepo   *store.UserRepository
	tokenRepo  *store.TokenRepository
	jwtManager *utils.JWTManager
}

func NewAuthService(userRepo *store.UserRepository, tokenRepo *store.TokenRepository, jwtManager *utils.JWTManager) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		jwtManager: jwtManager,
	}
}

func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	// Check if username already exists
	exists, err := s.userRepo.ExistsByUsername(req.Username)
	if err != nil {
		return nil, fmt.Errorf("error checking username: %w", err)
	}
	if exists {
		return nil, fmt.Errorf(constants.MsgUserAlreadyExists)
	}

	// Check if email already exists
	exists, err = s.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("error checking email: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("email already exists")
	}

	// Hash password
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	// Create user
	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
	}

	// Generate UUID v7 for user ID
	user.ID = utils.GenerateUUIDv7()

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	return &dto.RegisterResponse{
		Message: constants.MsgRegisterSuccess,
		User: dto.UserInfo{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

func (s *AuthService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// Get user by username
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf(constants.MsgInvalidCredentials)
		}
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	// Verify password
	valid, err := utils.VerifyPassword(req.Password, user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("error verifying password: %w", err)
	}

	if !valid {
		return nil, fmt.Errorf(constants.MsgInvalidCredentials)
	}

	// Generate access token
	accessToken, expiresAt, err := s.jwtManager.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("error generating access token: %w", err)
	}

	// Generate refresh token
	refreshTokenStr, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("error generating refresh token: %w", err)
	}

	// Hash refresh token for storage
	refreshTokenHash := utils.HashRefreshToken(refreshTokenStr)

	// Save refresh token
	refreshToken := &model.RefreshToken{
		ID:        utils.GenerateUUIDv7(),
		UserID:    user.ID,
		TokenHash: refreshTokenHash,
		ExpiresAt: time.Now().Add(s.jwtManager.GetRefreshTokenExpiration()),
	}

	if err := s.tokenRepo.Create(refreshToken); err != nil {
		return nil, fmt.Errorf("error saving refresh token: %w", err)
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		ExpiresAt:    expiresAt,
		User: dto.UserInfo{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

func (s *AuthService) RefreshToken(req *dto.RefreshRequest) (*dto.LoginResponse, error) {
	// Hash the provided refresh token
	tokenHash := utils.HashRefreshToken(req.RefreshToken)

	// Get refresh token from database
	refreshToken, err := s.tokenRepo.GetByTokenHash(tokenHash)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf(constants.MsgInvalidToken)
		}
		return nil, fmt.Errorf("error getting refresh token: %w", err)
	}

	// Check if token is expired
	if time.Now().After(refreshToken.ExpiresAt) {
		// Revoke expired token
		s.tokenRepo.RevokeToken(tokenHash)
		return nil, fmt.Errorf(constants.MsgTokenExpired)
	}

	// Generate new access token
	accessToken, expiresAt, err := s.jwtManager.GenerateAccessToken(refreshToken.User.ID, refreshToken.User.Username)
	if err != nil {
		return nil, fmt.Errorf("error generating access token: %w", err)
	}

	// Generate new refresh token
	newRefreshTokenStr, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("error generating refresh token: %w", err)
	}

	newRefreshTokenHash := utils.HashRefreshToken(newRefreshTokenStr)

	// Revoke old refresh token
	if err := s.tokenRepo.RevokeToken(tokenHash); err != nil {
		return nil, fmt.Errorf("error revoking old refresh token: %w", err)
	}

	// Save new refresh token
	newRefreshToken := &model.RefreshToken{
		ID:        utils.GenerateUUIDv7(),
		UserID:    refreshToken.User.ID,
		TokenHash: newRefreshTokenHash,
		ExpiresAt: time.Now().Add(s.jwtManager.GetRefreshTokenExpiration()),
	}

	if err := s.tokenRepo.Create(newRefreshToken); err != nil {
		return nil, fmt.Errorf("error saving refresh token: %w", err)
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshTokenStr,
		ExpiresAt:    expiresAt,
		User: dto.UserInfo{
			ID:        refreshToken.User.ID,
			Username:  refreshToken.User.Username,
			Email:     refreshToken.User.Email,
			CreatedAt: refreshToken.User.CreatedAt,
		},
	}, nil
}

func (s *AuthService) Logout(refreshToken string) error {
	tokenHash := utils.HashRefreshToken(refreshToken)
	return s.tokenRepo.RevokeToken(tokenHash)
}
