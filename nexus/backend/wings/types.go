package wings

import "encoding/json"

type ServerDetails struct {
	UUID        string `json:"uuid"`
	UUIDShort   string `json:"uuid_short"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Suspended   bool   `json:"suspended"`
	Memory      int    `json:"memory"`
	Disk        int    `json:"disk"`
	CPU         int    `json:"cpu"`
}

type ServerResources struct {
	CPUAbsolute    float64 `json:"cpu_absolute"`
	MemoryBytes    int64   `json:"memory_bytes"`
	MemoryLimit    int64   `json:"memory_limit_bytes"`
	DiskBytes      int64   `json:"disk_bytes"`
	DiskLimit      int64   `json:"disk_limit_bytes"`
	NetworkRXBytes int64   `json:"network_rx_bytes"`
	NetworkTXBytes int64   `json:"network_tx_bytes"`
	State          string  `json:"state"`
	Uptime         int64   `json:"uptime"`
}

type SystemInfo struct {
	Version   string `json:"version"`
	Uptime    int64  `json:"uptime"`
	Hostname  string `json:"hostname"`
	Architecture string `json:"architecture"`
}

type CreateServerPayload struct {
	UUID              string            `json:"uuid"`
	StartOnCompletion bool              `json:"start_on_completion"`
	Build             BuildConfig       `json:"build"`
	Container         ContainerConfig  `json:"container"`
	Allocation        AllocationConfig `json:"allocation"`
}

type BuildConfig struct {
	MemoryLimit    int     `json:"memory_limit"`
	Swap           int     `json:"swap"`
	Disk           int     `json:"disk"`
	IOWeight       int     `json:"io_weight"`
	CPU            int     `json:"cpu"`
	Threads        string  `json:"threads,omitempty"`
	OOMDisabled    bool    `json:"oom_disabled"`
}

type ContainerConfig struct {
	Image       string            `json:"image"`
	Startup     string            `json:"startup_command"`
	Environment map[string]string `json:"environment"`
}

type AllocationConfig struct {
	Default    int           `json:"default"`
	Additional []interface{} `json:"additional"`
}

type PowerActionPayload struct {
	Action string `json:"signal"`
}

type ErrorResponse struct {
	ErrMsg string `json:"error"`
	Trace  string `json:"trace,omitempty"`
}

func (e *ErrorResponse) Error() string {
	return e.ErrMsg
}

type SuccessResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func (s *SuccessResponse) UnmarshalJSON(data []byte) error {
	type Alias SuccessResponse
	aux := &struct{ *Alias }{Alias: (*Alias)(s)}
	return json.Unmarshal(data, &aux)
}
