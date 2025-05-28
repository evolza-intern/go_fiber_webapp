package models

// SystemStats represents system statistics
type SystemStats struct {
	TotalUsers    int     `json:"total_users"`
	ActiveUsers   int     `json:"active_users"`
	InactiveUsers int     `json:"inactive_users"`
	MemoryUsage   string  `json:"memory_usage"`
	CPUUsage      float64 `json:"cpu_usage"`
	Uptime        string  `json:"uptime"`
}
