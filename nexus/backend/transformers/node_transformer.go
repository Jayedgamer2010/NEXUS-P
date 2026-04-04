package transformers

import (
	"nexus/backend/models"
)

type NodeTransformed struct {
	ID              uint   `json:"id"`
	UUID            string `json:"uuid"`
	Public          bool   `json:"public"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	LocationID      int    `json:"location_id"`
	FQDN            string `json:"fqdn"`
	Scheme          string `json:"scheme"`
	BehindProxy     bool   `json:"behind_proxy"`
	MaintenanceMode bool   `json:"maintenance_mode"`
	Memory          int64  `json:"memory"`
	Disk            int64  `json:"disk"`
	DaemonListen    int    `json:"daemon_listen"`
	DaemonSFTP      int    `json:"daemon_sftp"`
	UsedMemory      int64  `json:"used_memory"`
	UsedDisk        int64  `json:"used_disk"`
	ServerCount     int64  `json:"server_count"`
	CreatedAt       string `json:"created_at"`
}

func TransformNode(node models.Node, usedMem, usedDisk, srvCount int64) NodeTransformed {
	return NodeTransformed{
		ID:              node.ID,
		UUID:            node.UUID,
		Public:          node.Public,
		Name:            node.Name,
		Description:     node.Description,
		LocationID:      node.LocationID,
		FQDN:            node.FQDN,
		Scheme:          node.Scheme,
		BehindProxy:     node.BehindProxy,
		MaintenanceMode: node.MaintenanceMode,
		Memory:          node.Memory,
		Disk:            node.Disk,
		DaemonListen:    node.DaemonListen,
		DaemonSFTP:      node.DaemonSFTP,
		UsedMemory:      usedMem,
		UsedDisk:        usedDisk,
		ServerCount:     srvCount,
		CreatedAt:       node.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

type EggTransformed struct {
	ID           uint   `json:"id"`
	UUID         string `json:"uuid"`
	Author       string `json:"author"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	DockerImage  string `json:"docker_image"`
	DockerImages string `json:"docker_images"`
	Startup      string `json:"startup"`
	CreatedAt    string `json:"created_at"`
}

func TransformEgg(egg models.Egg) EggTransformed {
	return EggTransformed{
		ID:           egg.ID,
		UUID:         egg.UUID,
		Author:       egg.Author,
		Name:         egg.Name,
		Description:  egg.Description,
		DockerImage:  egg.DockerImage,
		DockerImages: egg.DockerImages,
		Startup:      egg.Startup,
		CreatedAt:    egg.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func TransformEggs(eggs []models.Egg) []EggTransformed {
	result := make([]EggTransformed, len(eggs))
	for i, e := range eggs {
		result[i] = TransformEgg(e)
	}
	return result
}

type AllocationTransformed struct {
	ID         uint             `json:"id"`
	NodeID     uint             `json:"node_id"`
	IP         string           `json:"ip"`
	IPAlias    string           `json:"ip_alias"`
	Port       int              `json:"port"`
	ServerID   *uint            `json:"server_id"`
	Assigned   bool             `json:"assigned"`
	ServerName string           `json:"server_name,omitempty"`
	Notes      string           `json:"notes"`
	CreatedAt  string           `json:"created_at"`
}

func TransformAllocation(a models.Allocation) AllocationTransformed {
	t := AllocationTransformed{
		ID:         a.ID,
		NodeID:     a.NodeID,
		IP:         a.IP,
		IPAlias:    a.IPAlias,
		Port:       a.Port,
		ServerID:   a.ServerID,
		Assigned:   a.IsAssigned(),
		Notes:      a.Notes,
		CreatedAt:  a.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
	if a.Server != nil {
		t.ServerName = a.Server.Name
	}
	return t
}

func TransformAllocations(allocs []models.Allocation) []AllocationTransformed {
	result := make([]AllocationTransformed, len(allocs))
	for i, a := range allocs {
		result[i] = TransformAllocation(a)
	}
	return result
}
