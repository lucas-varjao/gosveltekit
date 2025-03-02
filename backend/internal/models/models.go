package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string `gorm:"unique;not null"`
	Email        string `gorm:"unique;not null"`
	DisplayName  string `gorm:"not null"`
	PasswordHash string `gorm:"not null"`
	Coins        int    `gorm:"default:0"`
	Role         string `gorm:"default:user"`
	LastLogin    time.Time
	Active       bool `gorm:"default:true"`
	RefreshToken string
	TokenExpiry  time.Time
}

type Match struct {
	gorm.Model
	TeamA    string    `gorm:"not null"`
	TeamB    string    `gorm:"not null"`
	Status   string    `gorm:"default:scheduled"` // scheduled, ongoing, finished
	Result   string    // ex: "2-1"
	StartsAt time.Time `gorm:"not null"`
}

type Bet struct {
	gorm.Model
	UserID     uint   `gorm:"not null"`
	MatchID    uint   `gorm:"not null"`
	Amount     int    `gorm:"not null;default:0"`
	Prediction string `gorm:"not null"`        // ex: "team_a_win", "draw"
	Status     string `gorm:"default:pending"` // pending, won, lost
}
