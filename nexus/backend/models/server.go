package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	StatusInstalling    = "installing"
	StatusInstallFailed = "install_failed"
	StatusSuspended     = "suspended"
	StatusRunning       = "running"
	StatusOffline       = "offline"
)

type Server struct {
	ID           uint        `gorm:"primaryKey" json:"id"`
	UUID         string      `gorm:"uniqueIndex;size:36" json:"uuid"`
	UUIDShort    string      `gorm:"uniqueIndex;size:8" json:"uuid_short"`
	Name         string      `gorm:"size:255" json:"name"`
	Description  string      `gorm:"type:text" json:"description"`
	UserID       uint        `gorm:"index;not null" json:"user_id"`
	NodeID       uint        `gorm:"index;not null" json:"node_id"`
	EggID        uint        `gorm:"index;not null" json:"egg_id"`
	AllocationID uint        `gorm:"index;not null" json:"allocation_id"`
	Memory       int         `gorm:"not null" json:"memory"`
	Disk         int         `gorm:"not null" json:"disk"`
	CPU          int         `gorm:"default:0" json:"cpu"`
	Swap         int         `gorm:"default:0" json:"swap"`
	IO           int         `gorm:"default:500" json:"io"`
	Image        string      `gorm:"size:255" json:"image"`
	Startup      string      `gorm:"type:text" json:"startup"`
	Environment  string      `gorm:"type:text" json:"environment"`
	Status       string      `gorm:"size:50;default:installing" json:"status"`
	Installed    bool        `gorm:"default:false" json:"installed"`
	Suspended    bool        `gorm:"default:false" json:"suspended"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	User         *User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Node         *Node       `gorm:"foreignKey:NodeID" json:"node,omitempty"`
	Egg          *Egg        `gorm:"foreignKey:EggID" json:"egg,omitempty"`
	Allocation   *Allocation `gorm:"foreignKey:AllocationID" json:"allocation,omitempty"`
}

func (s *Server) BeforeCreate(tx *gorm.DB) error {
	if s.UUID == "" {
		s.UUID = uuid.New().String()
		s.UUIDShort = s.UUID[:8]
	}
	return nil
}

func (s *Server) IsRunning() bool {
	return s.Status == StatusRunning
}

func (s *Server) IsSuspended() bool {
	return s.Suspended
}
