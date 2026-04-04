package models

import (
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type Node struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	UUID             string    `gorm:"size:36;uniqueIndex;not null" json:"uuid"`
	Public           bool      `gorm:"default:false" json:"public"`
	Name             string    `gorm:"size:100;not null" json:"name"`
	Description      string    `gorm:"type:text" json:"description"`
	LocationID       int       `gorm:"default:0" json:"location_id"`
	FQDN             string    `gorm:"size:255;not null" json:"fqdn"`
	Scheme           string    `gorm:"size:10;default:'https'" json:"scheme"`
	BehindProxy      bool      `gorm:"default:false" json:"behind_proxy"`
	MaintenanceMode  bool      `gorm:"default:false" json:"maintenance_mode"`
	Memory           int64     `gorm:"not null" json:"memory"`
	MemoryOveralloc  int       `gorm:"default:0" json:"memory_overallocate"`
	Disk             int64     `gorm:"not null" json:"disk"`
	DiskOveralloc    int       `gorm:"default:0" json:"disk_overallocate"`
	UploadSize       int64     `gorm:"default:100" json:"upload_size"`
	DaemonTokenID    string    `gorm:"size:64;not null" json:"-"`
	DaemonToken      string    `gorm:"size:255;not null" json:"-"`
	DaemonListen     int       `gorm:"default:8080" json:"daemon_listen"`
	DaemonSFTP       int       `gorm:"default:2022" json:"daemon_sftp"`
	DaemonBase       string    `gorm:"size:255;default:'/var/lib/pterodactyl'" json:"daemon_base"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (n *Node) BeforeCreate(tx *gorm.DB) error {
	if n.Scheme == "" {
		n.Scheme = "https"
	}
	if n.DaemonListen == 0 {
		n.DaemonListen = 8080
	}
	if n.DaemonSFTP == 0 {
		n.DaemonSFTP = 2022
	}
	if n.DaemonBase == "" {
		n.DaemonBase = "/var/lib/pterodactyl"
	}
	return nil
}

// GetConnectionAddress returns fqdn:port for connecting to Wings
func (n *Node) GetConnectionAddress() string {
	return fmt.Sprintf("%s:%d", n.FQDN, n.DaemonListen)
}

// IsOnline pings the Wings health endpoint
func (n *Node) IsOnline() bool {
	url := fmt.Sprintf("%s://%s:%d/api/health", n.Scheme, n.FQDN, n.DaemonListen)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// GetUsedMemory returns total memory allocated to servers on this node
func (n *Node) GetUsedMemory(db *gorm.DB) int64 {
	var total int64
	db.Model(&Server{}).Where("node_id = ?", n.ID).Select("COALESCE(SUM(memory), 0)").Scan(&total)
	return total
}

// GetUsedDisk returns total disk allocated to servers on this node
func (n *Node) GetUsedDisk(db *gorm.DB) int64 {
	var total int64
	db.Model(&Server{}).Where("node_id = ?", n.ID).Select("COALESCE(SUM(disk), 0)").Scan(&total)
	return total
}
