package service

import (
	"authorization/internal/constants"
	"authorization/internal/dto"
	"authorization/internal/store"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type UserService struct {
	userRepo *store.UserRepository
}

func NewUserService(userRepo *store.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetUserByID(userID string) (*dto.UserInfo, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf(constants.MsgUserNotFound)
		}
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return &dto.UserInfo{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}
