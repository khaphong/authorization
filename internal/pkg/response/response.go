package response

import (
	"authorization/internal/dto"
	"encoding/json"
	"net/http"
)

// JSON sends a JSON response with the given status code
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// Success sends a success response
func Success(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusOK, data)
}

// Created sends a created response
func Created(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusCreated, data)
}

// Error sends an error response
func Error(w http.ResponseWriter, statusCode int, code, message string) {
	JSON(w, statusCode, dto.ErrorResponse{
		Error: message,
		Code:  code,
	})
}

// ValidationError sends a validation error response
func ValidationError(w http.ResponseWriter, errors []dto.ValidationError) {
	JSON(w, http.StatusBadRequest, dto.ErrorResponse{
		Error:   "Validation failed",
		Code:    "VALIDATION_ERROR",
		Details: errors,
	})
}

// BadRequest sends a bad request response
func BadRequest(w http.ResponseWriter, message string) {
	Error(w, http.StatusBadRequest, "BAD_REQUEST", message)
}

// Unauthorized sends an unauthorized response
func Unauthorized(w http.ResponseWriter, message string) {
	Error(w, http.StatusUnauthorized, "UNAUTHORIZED", message)
}

// Forbidden sends a forbidden response
func Forbidden(w http.ResponseWriter, message string) {
	Error(w, http.StatusForbidden, "FORBIDDEN", message)
}

// NotFound sends a not found response
func NotFound(w http.ResponseWriter, message string) {
	Error(w, http.StatusNotFound, "NOT_FOUND", message)
}

// Conflict sends a conflict response
func Conflict(w http.ResponseWriter, message string) {
	Error(w, http.StatusConflict, "CONFLICT", message)
}

// InternalError sends an internal server error response
func InternalError(w http.ResponseWriter, message string) {
	Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", message)
}
