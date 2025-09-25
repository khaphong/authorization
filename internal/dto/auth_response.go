package dto

import "time"

// LoginResponse represents login response payload
type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         UserInfo  `json:"user"`
}

// UserInfo represents user information in responses
type UserInfo struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// RegisterResponse represents registration response payload
type RegisterResponse struct {
	Message string   `json:"message"`
	User    UserInfo `json:"user"`
}

// LogoutResponse represents logout response payload
type LogoutResponse struct {
	Message string `json:"message"`
}
