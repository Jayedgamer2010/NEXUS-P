package transformers

import "nexus/backend/models"

type NodeItem struct {
	ID          uint   `json:"id"`
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	FQDN        string `json:"fqdn"`
	Scheme      string `json:"scheme"`
	Memory      int    `json:"memory"`
	Disk        int    `json:"disk"`
	Maintenance bool   `json:"maintenance_mode"`
	ServerCount int    `json:"server_count"`
	AllocCount  int    `json:"allocation_count"`
}

type NodeDetail struct {
	ID                 uint        `json:"id"`
	UUID               string      `json:"uuid"`
	Name               string      `json:"name"`
	Description        string      `json:"description"`
	FQDN               string      `json:"fqdn"`
	Scheme             string      `json:"scheme"`
	BehindProxy        bool        `json:"behind_proxy"`
	MaintenanceMode    bool        `json:"maintenance_mode"`
	Memory             int         `json:"memory"`
	MemoryOverallocate int         `json:"memory_overallocate"`
	Disk               int         `json:"disk"`
	DiskOverallocate   int         `json:"disk_overallocate"`
	UploadSize         int         `json:"upload_size"`
	DaemonTokenID      string      `json:"daemon_token_id"`
	DaemonListen       int         `json:"daemon_listen"`
	DaemonSFTP         int         `json:"daemon_sftp"`
	DaemonBase         string      `json:"daemon_base"`
	CreatedAt          string      `json:"created_at"`
	UpdatedAt          string      `json:"updated_at"`
	Allocations        interface{} `json:"allocations,omitempty"`
}

func TransformNode(node models.Node) NodeItem {
	item := NodeItem{
		ID:          node.ID,
		UUID:        node.UUID,
		Name:        node.Name,
		Description: node.Description,
		FQDN:        node.FQDN,
		Scheme:      node.Scheme,
		Memory:      node.Memory,
		Disk:        node.Disk,
		Maintenance: node.MaintenanceMode,
		AllocCount:  len(node.Allocations),
	}
	if len(node.Servers) > 0 {
		item.ServerCount = len(node.Servers)
	}
	return item
}

func TransformNodeDetail(node models.Node) NodeDetail {
	detail := NodeDetail{
		ID:                 node.ID,
		UUID:               node.UUID,
		Name:               node.Name,
		Description:        node.Description,
		FQDN:               node.FQDN,
		Scheme:             node.Scheme,
		BehindProxy:        node.BehindProxy,
		MaintenanceMode:    node.MaintenanceMode,
		Memory:             node.Memory,
		MemoryOverallocate: node.MemoryOverallocate,
		Disk:               node.Disk,
		DiskOverallocate:   node.DiskOverallocate,
		UploadSize:         node.UploadSize,
		DaemonTokenID:      node.DaemonTokenID,
		DaemonListen:       node.DaemonListen,
		DaemonSFTP:         node.DaemonSFTP,
		DaemonBase:         node.DaemonBase,
		CreatedAt:          node.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:          node.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if len(node.Allocations) > 0 {
		allocs := make([]map[string]interface{}, len(node.Allocations))
		for i, a := range node.Allocations {
			allocs[i] = map[string]interface{}{
				"id":        a.ID,
				"node_id":   a.NodeID,
				"ip":        a.IP,
				"ip_alias":  a.IPAlias,
				"port":      a.Port,
				"notes":     a.Notes,
				"assigned":  a.IsAssigned(),
				"server_id": a.ServerID,
			}
		}
		detail.Allocations = allocs
	}

	return detail
}
