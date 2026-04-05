package models

import "fmt"

type Allocation struct {
	ID       uint    `gorm:"primaryKey" json:"id"`
	NodeID   uint    `gorm:"index;not null" json:"node_id"`
	IP       string  `gorm:"size:45;not null" json:"ip"`
	IPAlias  string  `gorm:"size:45" json:"ip_alias"`
	Port     int     `gorm:"not null" json:"port"`
	Notes    string  `gorm:"size:255" json:"notes"`
	ServerID *uint   `gorm:"index" json:"server_id"`
	Node     *Node   `gorm:"foreignKey:NodeID" json:"node,omitempty"`
	Server   *Server `gorm:"foreignKey:ServerID" json:"server,omitempty"`
}

func (a *Allocation) IsAssigned() bool {
	return a.ServerID != nil
}

func (a *Allocation) GetDisplayName() string {
	if a.IPAlias != "" {
		return fmt.Sprintf("%s:%d", a.IPAlias, a.Port)
	}
	return fmt.Sprintf("%s:%d", a.IP, a.Port)
}
