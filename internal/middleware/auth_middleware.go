package middleware

import (
	"authorization/internal/constants"
	"authorization/internal/pkg/response"
	"authorization/internal/utils"
	"context"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	jwtManager *utils.JWTManager
}

func NewAuthMiddleware(jwtManager *utils.JWTManager) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
	}
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.Unauthorized(w, constants.MsgUnauthorized)
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			response.Unauthorized(w, "Invalid authorization header format")
			return
		}

		claims, err := m.jwtManager.ValidateAccessToken(bearerToken[1])
		if err != nil {
			response.Unauthorized(w, constants.MsgInvalidToken)
			return
		}

		// Add user info to context
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "username", claims.Username)
		
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
