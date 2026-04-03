package models

import (
	"time"
)

type CoinTransaction struct {
	ID      uint      `gorm:"primaryKey" json:"id"`
	UserID  uint      `gorm:"not null;index" json:"user_id"`
	Amount  int       `gorm:"not null" json:"amount"` // can be negative
	Reason  string    `gorm:"size:255;not null" json:"reason"`
	CreatedAt time.Time `json:"created_at"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"user"`
}
