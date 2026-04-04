package transformers

import (
	"nexus/backend/models"
)

type ServerTransformed struct {
	ID           uint                 `json:"id"`
	UUID         string               `json:"uuid"`
	UUIDShort    string               `json:"uuid_short"`
	Name         string               `json:"name"`
	Description  string               `json:"description"`
	Status       string               `json:"status"`
	Suspended    bool                 `json:"suspended"`
	Memory       int                  `json:"memory"`
	Disk         int                  `json:"disk"`
	CPU          int                  `json:"cpu"`
	Node         *NodeBrief         `json:"node,omitempty"`
	Egg          *EggBrief          `json:"egg,omitempty"`
	User         *UserBrief         `json:"user,omitempty"`
	Allocation   *AllocationBrief   `json:"allocation,omitempty"`
	CreatedAt    string               `json:"created_at"`
}

type NodeBrief struct {
	ID   uint   `json:"id"`
	UUID string `json:"uuid"`
	Name string `json:"name"`
	FQDN string `json:"fqdn"`
}

type EggBrief struct {
	ID   uint   `json:"id"`
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type UserBrief struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type AllocationBrief struct {
	ID   uint   `json:"id"`
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

func TransformServer(server models.Server) ServerTransformed {
	s := ServerTransformed{
		ID:        server.ID,
		UUID:      server.UUID,
		UUIDShort: server.UUIDShort,
		Name:      server.Name,
		Status:    server.Status,
		Suspended: server.Suspended,
		Memory:    server.Memory,
		Disk:      server.Disk,
		CPU:       server.CPU,
		CreatedAt: server.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if server.Node.ID > 0 {
		s.Node = &NodeBrief{
			ID:   server.Node.ID,
			UUID: server.Node.UUID,
			Name: server.Node.Name,
			FQDN: server.Node.FQDN,
		}
	}
	if server.Egg.ID > 0 {
		s.Egg = &EggBrief{
			ID:   server.Egg.ID,
			UUID: server.Egg.UUID,
			Name: server.Egg.Name,
		}
	}
	if server.User.ID > 0 {
		s.User = &UserBrief{
			ID:       server.User.ID,
			Username: server.User.Username,
			Email:    server.User.Email,
		}
	}
	if server.Allocation.ID > 0 {
		s.Allocation = &AllocationBrief{
			ID:   server.Allocation.ID,
			IP:   server.Allocation.IP,
			Port: server.Allocation.Port,
		}
	}

	return s
}

func TransformServers(servers []models.Server) []ServerTransformed {
	result := make([]ServerTransformed, len(servers))
	for i, s := range servers {
		result[i] = TransformServer(s)
	}
	return result
}
