package dto

// RegisterRequest represents registration request payload
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginRequest represents login request payload
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// RefreshRequest represents refresh token request payload
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
