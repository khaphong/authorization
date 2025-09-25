package model

import "time"

// RefreshToken model
type RefreshToken struct {
	ID        string    `json:"id" gorm:"type:varchar(36);primaryKey"`
	UserID    string    `json:"user_id" gorm:"not null;index"`
	TokenHash string    `json:"-" gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null;index"`
	CreatedAt time.Time `json:"created_at"`
	Revoked   bool      `json:"revoked" gorm:"default:false;index"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
}

// TableName returns the table name for GORM
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
