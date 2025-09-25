package constants

// Error codes
const (
	ErrUserNotFound      = "USER_NOT_FOUND"
	ErrUserAlreadyExists = "USER_ALREADY_EXISTS"
	ErrInvalidCredentials = "INVALID_CREDENTIALS"
	ErrInvalidToken      = "INVALID_TOKEN"
	ErrTokenExpired      = "TOKEN_EXPIRED"
	ErrTokenRevoked      = "TOKEN_REVOKED"
	ErrValidation        = "VALIDATION_ERROR"
	ErrInternal          = "INTERNAL_ERROR"
	ErrUnauthorized      = "UNAUTHORIZED"
	ErrForbidden         = "FORBIDDEN"
)

// Error messages
const (
	MsgUserNotFound       = "User not found"
	MsgUserAlreadyExists  = "User already exists"
	MsgInvalidCredentials = "Invalid username or password"
	MsgInvalidToken       = "Invalid token"
	MsgTokenExpired       = "Token has expired"
	MsgTokenRevoked       = "Token has been revoked"
	MsgUnauthorized       = "Authentication required"
	MsgForbidden          = "Access denied"
	MsgInternalError      = "Internal server error"
)

// Success messages
const (
	MsgRegisterSuccess = "User registered successfully"
	MsgLoginSuccess    = "Login successful"
	MsgLogoutSuccess   = "Logged out successfully"
	MsgTokenRefreshed  = "Token refreshed successfully"
)
