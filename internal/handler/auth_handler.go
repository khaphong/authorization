package handler

import (
	"authorization/internal/constants"
	"authorization/internal/dto"
	"authorization/internal/pkg/logger"
	"authorization/internal/pkg/response"
	"authorization/internal/service"
	"encoding/json"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Failed to decode register request", zap.Error(err))
		response.BadRequest(w, "Invalid request body")
		return
	}

	// Basic validation
	if req.Username == "" || req.Email == "" || req.Password == "" {
		response.BadRequest(w, "Username, email, and password are required")
		return
	}

	if len(req.Password) < 6 {
		response.BadRequest(w, "Password must be at least 6 characters long")
		return
	}

	resp, err := h.authService.Register(&req)
	if err != nil {
		logger.Error("Registration failed", zap.Error(err), zap.String("username", req.Username))
		
		if strings.Contains(err.Error(), "already exists") {
			response.Conflict(w, err.Error())
			return
		}
		response.InternalError(w, constants.MsgInternalError)
		return
	}

	logger.Info("User registered successfully", zap.String("user_id", resp.User.ID), zap.String("username", resp.User.Username))
	response.Created(w, resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Failed to decode login request", zap.Error(err))
		response.BadRequest(w, "Invalid request body")
		return
	}

	// Basic validation
	if req.Username == "" || req.Password == "" {
		response.BadRequest(w, "Username and password are required")
		return
	}

	resp, err := h.authService.Login(&req)
	if err != nil {
		logger.Error("Login failed", zap.Error(err), zap.String("username", req.Username))
		
		if strings.Contains(err.Error(), constants.MsgInvalidCredentials) {
			response.Unauthorized(w, constants.MsgInvalidCredentials)
			return
		}
		response.InternalError(w, constants.MsgInternalError)
		return
	}

	logger.Info("User logged in successfully", zap.String("user_id", resp.User.ID), zap.String("username", resp.User.Username))
	response.Success(w, resp)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Failed to decode refresh request", zap.Error(err))
		response.BadRequest(w, "Invalid request body")
		return
	}

	if req.RefreshToken == "" {
		response.BadRequest(w, "Refresh token is required")
		return
	}

	resp, err := h.authService.RefreshToken(&req)
	if err != nil {
		logger.Error("Token refresh failed", zap.Error(err))
		
		if strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "expired") {
			response.Unauthorized(w, err.Error())
			return
		}
		response.InternalError(w, constants.MsgInternalError)
		return
	}

	logger.Info("Token refreshed successfully", zap.String("user_id", resp.User.ID))
	response.Success(w, resp)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Failed to decode logout request", zap.Error(err))
		response.BadRequest(w, "Invalid request body")
		return
	}

	if req.RefreshToken == "" {
		response.BadRequest(w, "Refresh token is required")
		return
	}

	err := h.authService.Logout(req.RefreshToken)
	if err != nil {
		logger.Error("Logout failed", zap.Error(err))
		response.InternalError(w, constants.MsgInternalError)
		return
	}

	logger.Info("User logged out successfully")
	response.Success(w, dto.LogoutResponse{
		Message: constants.MsgLogoutSuccess,
	})
}
