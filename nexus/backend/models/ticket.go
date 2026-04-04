package models

import (
	"time"
)

const (
	TicketOpen     = "open"
	TicketAnswered = "answered"
	TicketClosed   = "closed"
)

type Ticket struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UUID      string    `gorm:"size:36;uniqueIndex;not null" json:"uuid"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Subject   string    `gorm:"size:255;not null" json:"subject"`
	Message   string    `gorm:"type:text" json:"message"`
	Status    string    `gorm:"size:20;default:'open';index" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (t *Ticket) IsOpen() bool {
	return t.Status != TicketClosed
}
