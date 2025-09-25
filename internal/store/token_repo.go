package store

import (
	"authorization/internal/model"
	"time"

	"gorm.io/gorm"
)

type TokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) Create(token *model.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *TokenRepository) GetByTokenHash(tokenHash string) (*model.RefreshToken, error) {
	var token model.RefreshToken
	err := r.db.Preload("User").Where("token_hash = ? AND revoked = false", tokenHash).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *TokenRepository) GetByUserID(userID string) ([]*model.RefreshToken, error) {
	var tokens []*model.RefreshToken
	err := r.db.Where("user_id = ? AND revoked = false", userID).Find(&tokens).Error
	return tokens, err
}

func (r *TokenRepository) RevokeToken(tokenHash string) error {
	return r.db.Model(&model.RefreshToken{}).Where("token_hash = ?", tokenHash).Update("revoked", true).Error
}

func (r *TokenRepository) RevokeAllUserTokens(userID string) error {
	return r.db.Model(&model.RefreshToken{}).Where("user_id = ?", userID).Update("revoked", true).Error
}

func (r *TokenRepository) CleanExpiredTokens() error {
	return r.db.Where("expires_at < ? OR revoked = true", time.Now()).Delete(&model.RefreshToken{}).Error
}

func (r *TokenRepository) IsTokenValid(tokenHash string) (bool, error) {
	var token model.RefreshToken
	err := r.db.Where("token_hash = ? AND revoked = false AND expires_at > ?", tokenHash, time.Now()).First(&token).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
