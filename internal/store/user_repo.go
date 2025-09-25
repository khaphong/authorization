package store

import (
	"authorization/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) GetByID(id string) (*model.User, error) {
	var user model.User
	err := r.db.Where("id = ? AND is_deleted = false", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ? AND is_deleted = false", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ? AND is_deleted = false", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) SoftDelete(id string) error {
	var user model.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return err
	}
	
	user.SoftDelete()
	return r.db.Save(&user).Error
}

func (r *UserRepository) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("username = ? AND is_deleted = false", username).Count(&count).Error
	return count > 0, err
}

func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("email = ? AND is_deleted = false", email).Count(&count).Error
	return count > 0, err
}
