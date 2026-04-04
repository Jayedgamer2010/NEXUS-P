package requests

type CreateServerRequest struct {
	Name        string            `json:"name" validate:"required,min=1,max=48"`
	UserID      uint              `json:"user_id" validate:"required"`
	NodeID      uint              `json:"node_id" validate:"required"`
	EggID       uint              `json:"egg_id" validate:"required"`
	Memory      int               `json:"memory" validate:"required,min=128"`
	Disk        int               `json:"disk" validate:"required,min=256"`
	CPU         int               `json:"cpu" validate:"required,min=0,max=10000"`
	StartupCmd  string            `json:"startup" validate:"required"`
	DockerImage string            `json:"image" validate:"required"`
	Environment map[string]string `json:"environment"`
	Description string            `json:"description"`
}

type UpdateServerRequest struct {
	Name        *string `json:"name"`
	Memory      *int    `json:"memory"`
	Disk        *int    `json:"disk"`
	CPU         *int    `json:"cpu"`
	Status      *string `json:"status"`
	Suspended   *bool   `json:"suspended"`
	Description *string `json:"description"`
}

type PowerActionRequest struct {
	Action string `json:"action" validate:"required,oneof=start stop restart kill"`
}
