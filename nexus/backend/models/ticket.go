package models

import (
	"time"
)

type Ticket struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Subject   string    `gorm:"size:255;not null" json:"subject"`
	Status    string    `gorm:"size:20;default:'open';index" json:"status"` // open/closed/pending
	Priority  string    `gorm:"size:20;default:'medium';index" json:"priority"` // low/medium/high
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"user"`
}
