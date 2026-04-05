package wings

import "time"

type ServerResources struct {
	CPUAbsolute     float64 `json:"cpu_absolute"`
	MemoryBytes     int64   `json:"memory_bytes"`
	DiskBytes       int64   `json:"disk_bytes"`
	NetworkRxBytes  int64   `json:"network_rx_bytes"`
	NetworkTxBytes  int64   `json:"network_tx_bytes"`
	State           string  `json:"state"`
	Uptime          int64   `json:"uptime"`
}

type CreateServerPayload struct {
	UUID              string            `json:"uuid"`
	StartOnCompletion bool              `json:"start_on_completion"`
	Image             string            `json:"image"`
	Startup           StartupConfig     `json:"startup"`
	Environment       map[string]string `json:"environment"`
	Limits            ServerLimits      `json:"limits"`
	FeatureLimits     FeatureLimits     `json:"feature_limits"`
	Allocations       AllocationConfig  `json:"allocations"`
}

type ServerLimits struct {
	Memory  int    `json:"memory"`
	Swap    int    `json:"swap"`
	Disk    int    `json:"disk"`
	IO      int    `json:"io"`
	CPU     int    `json:"cpu"`
	Threads string `json:"threads"`
}

type StartupConfig struct {
	Done             string   `json:"done"`
	UserInteraction  []string `json:"user_interaction"`
	StripAnsi        []string `json:"strip_ansi"`
}

type AllocationConfig struct {
	Default    int   `json:"default"`
	Additional []int `json:"additional"`
}

type FeatureLimits struct {
	Databases   int `json:"databases"`
	Allocations int `json:"allocations"`
	Backups     int `json:"backups"`
}

type CreateTokenResponse struct {
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
}

type ConsoleWSResponse struct {
	Data struct {
		Token     string `json:"token"`
		Socket    string `json:"socket"`
	} `json:"data"`
}

type PowerPayload struct {
	Action string `json:"action"`
}

type SystemInfo struct {
	Version   string `json:"version"`
	Docker    bool   `json:"docker"`
	System    string `json:"system"`
}

type ServerInfo struct {
	State         string `json:"state"`
	Suspended     bool   `json:"suspended"`
	IsInstalling  bool   `json:"is_installing"`
	Uptime        int64  `json:"uptime"`
	MemoryBytes   int64  `json:"memory_bytes"`
	DiskBytes     int64  `json:"disk_bytes"`
	CPUAbsolute   float64 `json:"cpu_absolute"`
	NetworkRx     int64  `json:"network_rx_bytes"`
	NetworkTx     int64  `json:"network_tx_bytes"`
}

type ServerConsoleToken struct {
	Data struct {
		Token  string `json:"token"`
		Socket string `json:"socket"`
	} `json:"data"`
}

type _time struct{}

func (t _time) Now() time.Time {
	return time.Now()
}
