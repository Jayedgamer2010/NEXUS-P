package transformers

import "nexus/backend/models"

type EggItem struct {
	ID          uint   `json:"id"`
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	DockerImage string `json:"docker_image"`
	ServerCount int    `json:"server_count"`
}

type EggDetail struct {
	ID            uint   `json:"id"`
	UUID          string `json:"uuid"`
	Author        string `json:"author"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	DockerImage   string `json:"docker_image"`
	Startup       string `json:"startup"`
	ConfigStop    string `json:"config_stop"`
	ScriptInstall string `json:"script_install"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

func TransformEgg(egg models.Egg) EggItem {
	item := EggItem{
		ID:          egg.ID,
		UUID:        egg.UUID,
		Name:        egg.Name,
		Description: egg.Description,
		DockerImage: egg.DockerImage,
	}
	if len(egg.Servers) > 0 {
		item.ServerCount = len(egg.Servers)
	}
	return item
}

func TransformEggDetail(egg models.Egg) EggDetail {
	return EggDetail{
		ID:            egg.ID,
		UUID:          egg.UUID,
		Author:        egg.Author,
		Name:          egg.Name,
		Description:   egg.Description,
		DockerImage:   egg.DockerImage,
		Startup:       egg.Startup,
		ConfigStop:    egg.ConfigStop,
		ScriptInstall: egg.ScriptInstall,
		CreatedAt:     egg.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     egg.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func TransformEggs(eggs []models.Egg) []EggItem {
	items := make([]EggItem, len(eggs))
	for i, e := range eggs {
		items[i] = TransformEgg(e)
	}
	return items
}
