package wings

import (
	"encoding/json"
)

// Server resources response
type ServerResources struct {
	CPU        float64 `json:"cpu"`
	Memory     int64   `json:"memory"`     // bytes
	MemoryMax  int64   `json:"memory_max"` // bytes
	Disk       int64   `json:"disk"`
	DiskMax    int64   `json:"disk_max"`
	Uptime     int64   `json:"uptime"`
}

// Server details response from Wings
type ServerDetails struct {
	UUID        string `json:"uuid"`
	UUIDShort   string `json:"uuid_short"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Suspended   bool   `json:"suspended"`
	Node        struct {
		ID    uint   `json:"id"`
		UUID  string `json:"uuid"`
		Name  string `json:"name"`
		FQDN  string `json:"fqdn"`
	} `json:"node"`
	Egg struct {
		ID    uint   `json:"id"`
		UUID  string `json:"uuid"`
		Name  string `json:"name"`
	} `json:"egg"`
	Allocation struct {
		IP   string `json:"ip"`
		Port int    `json:"port"`
	} `json:"allocation"`
	Memory int64 `json:"memory"`
	Disk   int64 `json:"disk"`
	CPU    int   `json:"cpu"`
}

// CreateServerPayload - payload for creating a server
type CreateServerPayload struct {
	UUID              string `json:"uuid"`
	StartOnCompletion bool   `json:"start_on_completion"`
	// Additional fields would be nested in actual implementation
	// but for Phase 1 we keep it minimal
}

// PowerActionPayload - payload for power actions
type PowerActionPayload struct {
	Action string `json:"action"` // start, stop, restart, kill
}

// Error response from Wings
type ErrorResponse struct {
	ErrMsg string `json:"error"`
	Trace  string `json:"trace,omitempty"`
}

func (e *ErrorResponse) Error() string {
	return e.ErrMsg
}

// Success response wrapper
type SuccessResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// UnmarshalJSON helper for mixed responses
func (s *SuccessResponse) UnmarshalJSON(data []byte) error {
	type Alias SuccessResponse
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	return json.Unmarshal(data, &aux)
}
