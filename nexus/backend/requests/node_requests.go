package requests

type CreateNodeRequest struct {
	Name            string `json:"name" validate:"required,min=1,max=100"`
	FQDN            string `json:"fqdn" validate:"required"`
	Scheme          string `json:"scheme" validate:"required,oneof=http https"`
	Memory          int64  `json:"memory" validate:"required,min=1024"`
	Disk            int64  `json:"disk" validate:"required,min=5120"`
	DaemonToken     string `json:"daemon_token" validate:"required"`
	DaemonTokenID   string `json:"daemon_token_id" validate:"required"`
	DaemonListen    int    `json:"daemon_listen" validate:"required"`
	DaemonSFTP      int    `json:"daemon_sftp"`
	DaemonBase      string `json:"daemon_base"`
	Public          bool   `json:"public"`
	Description     string `json:"description"`
	BehindProxy     bool   `json:"behind_proxy"`
	LocationID      int    `json:"location_id"`
}

type UpdateNodeRequest struct {
	Name            *string `json:"name"`
	FQDN            *string `json:"fqdn"`
	Scheme          *string `json:"scheme"`
	Memory          *int64  `json:"memory"`
	Disk            *int64  `json:"disk"`
	DaemonToken     *string `json:"daemon_token"`
	DaemonTokenID   *string `json:"daemon_token_id"`
	DaemonListen    *int    `json:"daemon_listen"`
	DaemonSFTP      *int    `json:"daemon_sftp"`
	DaemonBase      *string `json:"daemon_base"`
	Public          *bool   `json:"public"`
	Description     *string `json:"description"`
	MaintenanceMode *bool   `json:"maintenance_mode"`
}

type CreateAllocationRequest struct {
	IP      string `json:"ip" validate:"required"`
	Port    int    `json:"port" validate:"required,min=1,max=65535"`
	IPAlias string `json:"ip_alias"`
	Notes   string `json:"notes"`
}
