package models

import (
	"time"
)

type Egg struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UUID           string    `gorm:"size:36;uniqueIndex;not null" json:"uuid"`
	Name           string    `gorm:"size:191;not null" json:"name"`
	Description    string    `gorm:"type:text" json:"description"`
	DockerImage    string    `gorm:"size:255;not null" json:"docker_image"`
	StartupCommand string    `gorm:"type:text;not null" json:"startup_command"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
