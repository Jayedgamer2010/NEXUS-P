package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Egg struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UUID          string    `gorm:"uniqueIndex;size:36" json:"uuid"`
	Author        string    `gorm:"size:255" json:"author"`
	Name          string    `gorm:"size:255" json:"name"`
	Description   string    `gorm:"type:text" json:"description"`
	DockerImage   string    `gorm:"size:255" json:"docker_image"`
	Startup       string    `gorm:"type:text" json:"startup"`
	ConfigStop    string    `gorm:"size:255" json:"config_stop"`
	ScriptInstall string    `gorm:"type:text" json:"script_install"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Servers       []Server  `gorm:"foreignKey:EggID" json:"servers,omitempty"`
}

func (e *Egg) BeforeCreate(tx *gorm.DB) error {
	if e.UUID == "" {
		e.UUID = uuid.New().String()
	}
	return nil
}
