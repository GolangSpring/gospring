package repository

import (
	"fmt"
	"github.com/lib/pq"
	"time"
)

type User struct {
	Name       string         `gorm:"type:varchar(100);not null" json:"name"`
	Email      string         `gorm:"type:varchar(100);unique;not null" json:"email"`
	Password   string         `gorm:"type:varchar(255);not null" json:"-"` // Excluded from JSON responses
	IsVerified bool           `gorm:"default:false" json:"is_verified"`
	Roles      pq.StringArray `gorm:"type:text[]" json:"roles"`

	ID        uint       `gorm:"primaryKey" json:"id"` // Auto-increment primary key
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (user *User) Validate() error {
	if user.Name == "" {
		return fmt.Errorf("name is required")
	}
	if user.Email == "" {
		return fmt.Errorf("email is required")
	}
	if user.Password == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}
