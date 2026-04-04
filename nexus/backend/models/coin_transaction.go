package models

import (
	"time"
)

type CoinTransaction struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UUID      string    `gorm:"size:36;uniqueIndex;not null" json:"uuid"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Amount    int       `gorm:"not null" json:"amount"`
	Balance   int       `gorm:"not null" json:"balance"`
	Reason    string    `gorm:"size:255;not null" json:"reason"`
	CreatedAt time.Time `json:"created_at"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
