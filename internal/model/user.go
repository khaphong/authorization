package model

// User model
type User struct {
	BaseModel
	Username     string `json:"username" gorm:"uniqueIndex:idx_username_active,where:is_deleted=false;not null"`
	Email        string `json:"email" gorm:"uniqueIndex:idx_email_active,where:is_deleted=false;not null"`
	PasswordHash string `json:"-" gorm:"not null"`
}

// TableName returns the table name for GORM
func (User) TableName() string {
	return "users"
}
