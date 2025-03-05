// backend/internal/models/models.go

package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	// Identity information
	Username     string `gorm:"unique;not null;index"`
	Email        string `gorm:"unique;not null;index"`
	DisplayName  string `gorm:"not null"`
	PasswordHash string `gorm:"not null"`

	// Profile information
	FirstName string
	LastName  string

	// Account status
	Active        bool `gorm:"default:true"`
	EmailVerified bool `gorm:"default:false"`
	LastLogin     time.Time
	LastActive    time.Time

	// Access control
	Role        string `gorm:"default:user"`
	Permissions string `gorm:"type:text"` // JSON string of permissions

	// Authentication
	RefreshToken     string
	TokenExpiry      time.Time
	ResetToken       string
	ResetTokenExpiry time.Time
}
