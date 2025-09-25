package handler

import (
	"authorization/internal/constants"
	"authorization/internal/pkg/logger"
	"authorization/internal/pkg/response"
	"authorization/internal/service"
	"net/http"

	"go.uber.org/zap"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		logger.Error("User ID not found in context")
		response.InternalError(w, constants.MsgInternalError)
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		logger.Error("Failed to get user", zap.Error(err), zap.String("user_id", userID))
		
		if err.Error() == constants.MsgUserNotFound {
			response.NotFound(w, constants.MsgUserNotFound)
			return
		}
		response.InternalError(w, constants.MsgInternalError)
		return
	}

	logger.Info("User info retrieved successfully", zap.String("user_id", userID))
	response.Success(w, user)
}
