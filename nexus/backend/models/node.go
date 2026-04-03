package models

import (
	"time"
)

type Node struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UUID            string    `gorm:"size:36;uniqueIndex;not null" json:"uuid"`
	Name            string    `gorm:"size:100;not null" json:"name"`
	FQDN            string    `gorm:"size:255;not null" json:"fqdn"`
	Scheme          string    `gorm:"size:10;default:'https'" json:"scheme"` // http or https
	WingsPort       int       `gorm:"default:8080" json:"wings_port"`
	Memory          int64     `gorm:"not null" json:"memory"`           // MB
	MemoryOveralloc int       `gorm:"default:0" json:"memory_overalloc"` // percentage
	Disk            int64     `gorm:"not null" json:"disk"`             // MB
	DiskOveralloc   int       `gorm:"default:0" json:"disk_overalloc"`   // percentage
	TokenID         string    `gorm:"size:64;not null" json:"token_id"`
	Token           string    `gorm:"size:255;not null" json:"-"`       // encrypted in production
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
