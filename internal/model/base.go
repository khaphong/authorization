package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseModel contains common fields for all models
type BaseModel struct {
	ID        string         `json:"id" gorm:"type:varchar(36);primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt *time.Time     `json:"deleted_at,omitempty" gorm:"index"`
	IsDeleted bool           `json:"is_deleted" gorm:"default:false;index"`
}

// BeforeCreate hook to generate UUID
func (b *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == "" {
		// Generate UUID v7 by default
		b.ID = uuid.New().String()
	}
	return
}

// SoftDelete performs soft delete by setting deleted_at and is_deleted
func (b *BaseModel) SoftDelete() {
	now := time.Now()
	b.DeletedAt = &now
	b.IsDeleted = true
}
