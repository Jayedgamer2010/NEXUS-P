package requests

type CreateServerRequest struct {
	UserID      uint   `json:"user_id" validate:"required"`
	NodeID      uint   `json:"node_id" validate:"required"`
	EggID       uint   `json:"egg_id" validate:"required"`
	Name        string `json:"name" validate:"required,min=2,max=255"`
	Description string `json:"description"`
	Memory      int    `json:"memory" validate:"required,min=128"`
	Disk        int    `json:"disk" validate:"required,min=100"`
	CPU         int    `json:"cpu"`
	Swap        int    `json:"swap"`
}

type UpdateServerRequest struct {
	Name        string `json:"name" validate:"min=1,max=255"`
	Description string `json:"description"`
	Memory      *int   `json:"memory" validate:"omitempty,min=128"`
	Disk        *int   `json:"disk" validate:"omitempty,min=100"`
	CPU         *int   `json:"cpu"`
	Swap        *int   `json:"swap"`
}

type PowerActionRequest struct {
	Action string `json:"action" validate:"required,oneof=start stop restart kill"`
}

type ServerUpdateDetailRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=255"`
	Description string `json:"description"`
}
