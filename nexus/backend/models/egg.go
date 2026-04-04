package models

import (
	"time"
)

type Egg struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UUID            string    `gorm:"size:36;uniqueIndex;not null" json:"uuid"`
	Author          string    `gorm:"size:255;not null" json:"author"`
	Name            string    `gorm:"size:191;not null" json:"name"`
	Description     string    `gorm:"type:text" json:"description"`
	DockerImage     string    `gorm:"size:255;not null" json:"docker_image"`
	DockerImages    string    `gorm:"type:json" json:"docker_images"`
	Startup         string    `gorm:"type:text" json:"startup"`
	ConfigFiles     string    `gorm:"type:json" json:"config_files"`
	ConfigStartup   string    `gorm:"type:json" json:"config_startup"`
	ConfigStop      string    `gorm:"type:text" json:"config_stop"`
	ConfigLogs      string    `gorm:"type:json" json:"config_logs"`
	ScriptInstall   string    `gorm:"type:text" json:"script_install"`
	ScriptEntry     string    `gorm:"type:text" json:"script_entry"`
	ScriptContainer string    `gorm:"type:text" json:"script_container"`
	CopyScriptFrom  string    `gorm:"size:255" json:"copy_script_from"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
