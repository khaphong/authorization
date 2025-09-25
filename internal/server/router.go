package server

import (
	"authorization/internal/handler"
	"authorization/internal/middleware"
	"authorization/internal/pkg/response"
	"authorization/internal/utils"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewRouter(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	jwtManager *utils.JWTManager,
) *chi.Mux {
	r := chi.NewRouter()

	// Built-in middleware
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.Timeout(60 * time.Second))

	// Custom middleware
	r.Use(middleware.LoggingMiddleware)

	// CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		response.Success(w, map[string]interface{}{
			"status":    "healthy",
			"service":   "authorization",
			"timestamp": time.Now().Unix(),
		})
	})

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Auth routes (public)
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
			r.Post("/refresh", authHandler.RefreshToken)
			r.Post("/logout", authHandler.Logout)
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			authMiddleware := middleware.NewAuthMiddleware(jwtManager)
			r.Use(authMiddleware.RequireAuth)

			r.Get("/me", userHandler.Me)
		})
	})

	return r
}
