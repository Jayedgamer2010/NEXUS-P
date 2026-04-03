package models

import (
	"time"
)

type Allocation struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	NodeID    uint      `gorm:"not null;index" json:"node_id"`
	IP        string    `gorm:"size:45;not null" json:"ip"`
	Port      int       `gorm:"not null" json:"port"`
	Assigned  bool      `gorm:"default:false;index" json:"assigned"`
	ServerID  *uint     `gorm:"index" json:"server_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Node    Node    `gorm:"foreignKey:NodeID" json:"node"`
	Server  *Server `gorm:"foreignKey:ServerID" json:"server"`
}
