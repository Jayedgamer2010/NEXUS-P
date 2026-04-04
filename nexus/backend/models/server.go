package models

import (
	"time"

	"gorm.io/gorm"
)

// Server status constants
const (
	StatusInstalling    = "installing"
	StatusInstallFailed = "install_failed"
	StatusSuspended     = "suspended"
	StatusRunning       = "running"
	StatusOffline       = "offline"
	StatusStarting      = "starting"
	StatusStopping      = "stopping"
)

// Power action constants
const (
	PowerStart  = "start"
	PowerStop   = "stop"
	PowerRestart = "restart"
	PowerKill   = "kill"
)

type Server struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	UUID             string    `gorm:"size:36;uniqueIndex;not null" json:"uuid"`
	UUIDShort        string    `gorm:"size:8;uniqueIndex;not null" json:"uuid_short"`
	Name             string    `gorm:"size:191;not null" json:"name"`
	Description      string    `gorm:"type:text" json:"description"`
	UserID           uint      `gorm:"not null;index" json:"user_id"`
	NodeID           uint      `gorm:"not null;index" json:"node_id"`
	EggID            uint      `gorm:"not null;index" json:"egg_id"`
	AllocationID     uint      `gorm:"not null;index" json:"allocation_id"`
	Memory           int       `gorm:"not null" json:"memory"`
	MemoryOveralloc  int       `gorm:"default:0" json:"memory_overallocate"`
	Disk             int       `gorm:"not null" json:"disk"`
	DiskOveralloc    int       `gorm:"default:0" json:"disk_overallocate"`
	CPU              int       `gorm:"not null" json:"cpu"`
	Threads          string    `gorm:"size:20" json:"threads"`
	IO               int       `gorm:"default:500" json:"io"`
	Image            string    `gorm:"size:255" json:"image"`
	Startup          string    `gorm:"type:text" json:"startup"`
	EnvVariables     string    `gorm:"type:json" json:"env_variables"`
	Status           string    `gorm:"size:50;default:'installing';index" json:"status"`
	Installed        bool      `gorm:"default:false" json:"installed"`
	Suspended        bool      `gorm:"default:false;index" json:"suspended"`
	SkipScripts      bool      `gorm:"default:false" json:"skip_scripts"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	User       User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Node       Node       `gorm:"foreignKey:NodeID" json:"node,omitempty"`
	Egg        Egg        `gorm:"foreignKey:EggID" json:"egg,omitempty"`
	Allocation Allocation `gorm:"foreignKey:AllocationID" json:"allocation,omitempty"`
}

func (s *Server) BeforeCreate(tx *gorm.DB) error {
	if s.Status == "" {
		s.Status = StatusInstalling
	}
	if s.IO == 0 {
		s.IO = 500
	}
	return nil
}

func (s *Server) IsRunning() bool {
	return s.Status == StatusRunning || s.Status == StatusStarting
}

func (s *Server) IsSuspended() bool {
	return s.Suspended
}

func (s *Server) GetShortUUID() string {
	return s.UUIDShort
}
