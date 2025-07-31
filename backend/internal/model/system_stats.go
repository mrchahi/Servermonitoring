package model

// SystemStats represents current system resource usage
type SystemStats struct {
    CPU     CPUStats    `json:"cpu"`
    Memory  MemoryStats `json:"memory"`
    Disk    DiskStats   `json:"disk"`
    Network NetworkStats `json:"network"`
    System  SystemInfo  `json:"system"`
}

// CPUStats represents CPU usage statistics
type CPUStats struct {
    UsagePercent float64 `json:"usagePercent"`
    Temperature  float64 `json:"temperature,omitempty"`
}

// MemoryStats represents memory usage statistics
type MemoryStats struct {
    Total       uint64  `json:"total"`
    Used        uint64  `json:"used"`
    Free        uint64  `json:"free"`
    UsagePercent float64 `json:"usagePercent"`
}

// DiskStats represents disk usage statistics
type DiskStats struct {
    Total       uint64  `json:"total"`
    Used        uint64  `json:"used"`
    Free        uint64  `json:"free"`
    UsagePercent float64 `json:"usagePercent"`
}

// NetworkStats represents network usage statistics
type NetworkStats struct {
    BytesSent    uint64 `json:"bytesSent"`
    BytesReceived uint64 `json:"bytesReceived"`
}

// SystemInfo represents general system information
type SystemInfo struct {
    Hostname    string  `json:"hostname"`
    Uptime     float64 `json:"uptime"`
    LoadAverage []float64 `json:"loadAverage"`
}
