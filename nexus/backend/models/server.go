package models

import (
	"time"
)

type Server struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UUID           string    `gorm:"size:36;uniqueIndex;not null" json:"uuid"`
	Name           string    `gorm:"size:191;not null" json:"name"`
	UserID         uint      `gorm:"not null;index" json:"user_id"`
	NodeID         uint      `gorm:"not null;index" json:"node_id"`
	EggID          uint      `gorm:"not null;index" json:"egg_id"`
	AllocationID   uint      `gorm:"not null;uniqueIndex" json:"allocation_id"`
	Memory         int64     `gorm:"not null" json:"memory"`     // MB
	Disk           int64     `gorm:"not null" json:"disk"`       // MB
	CPU            int       `gorm:"not null" json:"cpu"`        // percentage
	Status         string    `gorm:"size:50;default:'installing';index" json:"status"` // installing/running/stopped/error
	Suspended      bool      `gorm:"default:false;index" json:"suspended"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	// Relations
	User       User       `gorm:"foreignKey:UserID" json:"user"`
	Node       Node       `gorm:"foreignKey:NodeID" json:"node"`
	Egg        Egg        `gorm:"foreignKey:EggID" json:"egg"`
	Allocation Allocation `gorm:"foreignKey:AllocationID" json:"allocation"`
}
