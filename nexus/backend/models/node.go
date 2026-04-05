package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Node struct {
	ID                  uint         `gorm:"primaryKey" json:"id"`
	UUID                string       `gorm:"uniqueIndex;size:36" json:"uuid"`
	Name                string       `gorm:"size:255" json:"name"`
	Description         string       `gorm:"size:500" json:"description"`
	FQDN                string       `gorm:"size:255" json:"fqdn"`
	Scheme              string       `gorm:"size:5;default:https" json:"scheme"`
	BehindProxy         bool         `gorm:"default:false" json:"behind_proxy"`
	MaintenanceMode     bool         `gorm:"default:false" json:"maintenance_mode"`
	Memory              int          `gorm:"not null" json:"memory"`
	MemoryOverallocate  int          `gorm:"default:0" json:"memory_overallocate"`
	Disk                int          `gorm:"not null" json:"disk"`
	DiskOverallocate    int          `gorm:"default:0" json:"disk_overallocate"`
	UploadSize          int          `gorm:"default:100" json:"upload_size"`
	DaemonTokenID       string       `gorm:"size:255" json:"daemon_token_id"`
	DaemonToken         string       `gorm:"size:255" json:"-"`
	DaemonListen        int          `gorm:"default:8080" json:"daemon_listen"`
	DaemonSFTP          int          `gorm:"default:2022" json:"daemon_sftp"`
	DaemonBase          string       `gorm:"size:255;default:/var/lib/pterodactyl" json:"daemon_base"`
	CreatedAt           time.Time    `json:"created_at"`
	UpdatedAt           time.Time    `json:"updated_at"`
	Allocations         []Allocation `gorm:"foreignKey:NodeID" json:"allocations,omitempty"`
	Servers             []Server     `gorm:"foreignKey:NodeID" json:"servers,omitempty"`
}

func (n *Node) BeforeCreate(tx *gorm.DB) error {
	if n.UUID == "" {
		n.UUID = uuid.New().String()
	}
	return nil
}

func (n *Node) GetConnectionAddress() string {
	return fmt.Sprintf("%s://%s:%d", n.Scheme, n.FQDN, n.DaemonListen)
}
