package transformers

import "nexus/backend/models"

type ServerItem struct {
	ID           uint   `json:"id"`
	UUID         string `json:"uuid"`
	UUIDShort    string `json:"uuid_short"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	Suspended    bool   `json:"suspended"`
	Installed    bool   `json:"installed"`
	Memory       int    `json:"memory"`
	Disk         int    `json:"disk"`
	CPU          int    `json:"cpu"`
	NodeName     string `json:"node_name,omitempty"`
	UserName     string `json:"user_name,omitempty"`
	EggName      string `json:"egg_name,omitempty"`
	Allocation   string `json:"allocation,omitempty"`
	CreatedAt    string `json:"created_at"`
}

type ServerDetail struct {
	ID           uint   `json:"id"`
	UUID         string `json:"uuid"`
	UUIDShort    string `json:"uuid_short"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Status       string `json:"status"`
	Suspended    bool   `json:"suspended"`
	Installed    bool   `json:"installed"`
	Memory       int    `json:"memory"`
	Disk         int    `json:"disk"`
	CPU          int    `json:"cpu"`
	Swap         int    `json:"swap"`
	IO           int    `json:"io"`
	Image        string `json:"image"`
	Startup      string `json:"startup"`
	UserID       uint   `json:"user_id"`
	NodeID       uint   `json:"node_id"`
	EggID        uint   `json:"egg_id"`
	AllocationID uint   `json:"allocation_id"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

func TransformServer(server models.Server) ServerItem {
	item := ServerItem{
		ID:        server.ID,
		UUID:      server.UUID,
		UUIDShort: server.UUIDShort,
		Name:      server.Name,
		Status:    server.Status,
		Suspended: server.Suspended,
		Installed: server.Installed,
		Memory:    server.Memory,
		Disk:      server.Disk,
		CPU:       server.CPU,
		CreatedAt: server.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if server.User != nil {
		item.UserName = server.User.Username
	}
	if server.Node != nil {
		item.NodeName = server.Node.Name
	}
	if server.Egg != nil {
		item.EggName = server.Egg.Name
	}
	if server.Allocation != nil {
		item.Allocation = server.Allocation.GetDisplayName()
	}

	return item
}

func TransformServerDetail(server models.Server) ServerDetail {
	return ServerDetail{
		ID:           server.ID,
		UUID:         server.UUID,
		UUIDShort:    server.UUIDShort,
		Name:         server.Name,
		Description:  server.Description,
		Status:       server.Status,
		Suspended:    server.Suspended,
		Installed:    server.Installed,
		Memory:       server.Memory,
		Disk:         server.Disk,
		CPU:          server.CPU,
		Swap:         server.Swap,
		IO:           server.IO,
		Image:        server.Image,
		Startup:      server.Startup,
		UserID:       server.UserID,
		NodeID:       server.NodeID,
		EggID:        server.EggID,
		AllocationID: server.AllocationID,
		CreatedAt:    server.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:    server.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func TransformServers(servers []models.Server) []ServerItem {
	items := make([]ServerItem, len(servers))
	for i, s := range servers {
		items[i] = TransformServer(s)
	}
	return items
}
