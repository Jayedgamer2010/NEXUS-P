package requests

type CreateNodeRequest struct {
	Name             string `json:"name" validate:"required,min=2,max=255"`
	Description      string `json:"description"`
	FQDN             string `json:"fqdn" validate:"required,min=3,max=255"`
	Scheme           string `json:"scheme" validate:"required,oneof=http https"`
	Memory           int    `json:"memory" validate:"required,min=1"`
	MemoryOveralloc  int    `json:"memory_overallocate"`
	Disk             int    `json:"disk" validate:"required,min=1"`
	DiskOveralloc    int    `json:"disk_overallocate"`
	DaemonListen     int    `json:"daemon_listen"`
	DaemonSFTP       int    `json:"daemon_sftp"`
	DaemonTokenID    string `json:"daemon_token_id"`
	DaemonToken      string `json:"daemon_token"`
	DaemonBase       string `json:"daemon_base"`
	BehindProxy      bool   `json:"behind_proxy"`
	MaintenanceMode  bool   `json:"maintenance_mode"`
}

type UpdateNodeRequest struct {
	Name             string `json:"name" validate:"omitempty,min=2,max=255"`
	Description      string `json:"description"`
	FQDN             string `json:"fqdn" validate:"omitempty,min=3,max=255"`
	Scheme           string `json:"scheme" validate:"omitempty,oneof=http https"`
	Memory           int    `json:"memory" validate:"omitempty,min=1"`
	MemoryOveralloc  int    `json:"memory_overallocate"`
	Disk             int    `json:"disk" validate:"omitempty,min=1"`
	DiskOveralloc    int    `json:"disk_overallocate"`
	DaemonListen     int    `json:"daemon_listen"`
	DaemonSFTP       int    `json:"daemon_sftp"`
	DaemonTokenID    string `json:"daemon_token_id"`
	DaemonToken      string `json:"daemon_token"`
	DaemonBase       string `json:"daemon_base"`
	BehindProxy      bool   `json:"behind_proxy"`
	MaintenanceMode  bool   `json:"maintenance_mode"`
}

type CreateAllocationRequest struct {
	IP      string `json:"ip" validate:"required"`
	Port    int    `json:"port" validate:"required,min=1,max=65535"`
	IPAlias string `json:"ip_alias"`
	Notes   string `json:"notes"`
}

type CreateAllocationsRequest struct {
	IP        string `json:"ip" validate:"required"`
	IPAlias   string `json:"ip_alias"`
	PortStart int    `json:"port_start" validate:"required,min=1,max=65535"`
	PortEnd   int    `json:"port_end" validate:"required,min=1,max=65535"`
	Notes     string `json:"notes"`
}
