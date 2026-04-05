package models

import "time"

type Ticket struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	Subject   string    `gorm:"size:255" json:"subject"`
	Status    string    `gorm:"size:20;default:open" json:"status"`
	Priority  string    `gorm:"size:20;default:low" json:"priority"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
