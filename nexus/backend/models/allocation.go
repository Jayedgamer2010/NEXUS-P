package models

import (
	"fmt"
	"time"
)

type Allocation struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	NodeID    uint      `gorm:"not null;index" json:"node_id"`
	IP        string    `gorm:"size:45;not null" json:"ip"`
	IPAlias   string    `gorm:"size:45" json:"ip_alias"`
	Port      int       `gorm:"not null" json:"port"`
	ServerID  *uint     `gorm:"index" json:"server_id"`
	Notes     string    `gorm:"type:text" json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Node   Node    `gorm:"foreignKey:NodeID" json:"node,omitempty"`
	Server *Server `gorm:"foreignKey:ServerID" json:"server,omitempty"`
}

func (a *Allocation) IsAssigned() bool {
	return a.ServerID != nil
}

func (a *Allocation) GetDisplayName() string {
	return fmt.Sprintf("%s:%d", a.IP, a.Port)
}

func (a *Allocation) Assign(serverID uint) {
	a.ServerID = &serverID
}

func (a *Allocation) Unassign() {
	a.ServerID = nil
}
