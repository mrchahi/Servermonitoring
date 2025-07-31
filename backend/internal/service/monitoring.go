package service

import (
    "time"
    "github.com/mrchahi/Servermonitoring/internal/model"
    "github.com/shirou/gopsutil/v3/cpu"
    "github.com/shirou/gopsutil/v3/mem"
    "github.com/shirou/gopsutil/v3/disk"
    "github.com/shirou/gopsutil/v3/host"
    "github.com/shirou/gopsutil/v3/net"
)

// MonitoringService handles system monitoring operations
type MonitoringService struct {
    updateInterval time.Duration
    subscribers    map[chan *model.SystemStats]bool
}

// NewMonitoringService creates a new monitoring service
func NewMonitoringService(updateInterval time.Duration) *MonitoringService {
    return &MonitoringService{
        updateInterval: updateInterval,
        subscribers:    make(map[chan *model.SystemStats]bool),
    }
}

// Start begins monitoring system resources
func (s *MonitoringService) Start() {
    go func() {
        ticker := time.NewTicker(s.updateInterval)
        defer ticker.Stop()

        for range ticker.C {
            stats, err := s.collectStats()
            if err != nil {
                continue
            }

            s.broadcast(stats)
        }
    }()
}

// Subscribe returns a channel that receives system stats updates
func (s *MonitoringService) Subscribe() chan *model.SystemStats {
    ch := make(chan *model.SystemStats)
    s.subscribers[ch] = true
    return ch
}

// Unsubscribe removes a subscriber
func (s *MonitoringService) Unsubscribe(ch chan *model.SystemStats) {
    delete(s.subscribers, ch)
    close(ch)
}

// collectStats gathers current system statistics
func (s *MonitoringService) collectStats() (*model.SystemStats, error) {
    stats := &model.SystemStats{}

    // CPU usage
    cpuPercent, err := cpu.Percent(time.Second, false)
    if err == nil && len(cpuPercent) > 0 {
        stats.CPU.UsagePercent = cpuPercent[0]
    }

    // Memory usage
    if vmstat, err := mem.VirtualMemory(); err == nil {
        stats.Memory = model.MemoryStats{
            Total:        vmstat.Total,
            Used:         vmstat.Used,
            Free:         vmstat.Free,
            UsagePercent: vmstat.UsedPercent,
        }
    }

    // Disk usage
    if diskStat, err := disk.Usage("/"); err == nil {
        stats.Disk = model.DiskStats{
            Total:        diskStat.Total,
            Used:         diskStat.Used,
            Free:         diskStat.Free,
            UsagePercent: diskStat.UsedPercent,
        }
    }

    // Network stats
    if netStats, err := net.IOCounters(false); err == nil && len(netStats) > 0 {
        stats.Network = model.NetworkStats{
            BytesSent:    netStats[0].BytesSent,
            BytesReceived: netStats[0].BytesRecv,
        }
    }

    // System info
    if hostInfo, err := host.Info(); err == nil {
        loadAvg, _ := load.Avg()
        stats.System = model.SystemInfo{
            Hostname:    hostInfo.Hostname,
            Uptime:      float64(hostInfo.Uptime),
            LoadAverage: []float64{loadAvg.Load1, loadAvg.Load5, loadAvg.Load15},
        }
    }

    return stats, nil
}

// broadcast sends updates to all subscribers
func (s *MonitoringService) broadcast(stats *model.SystemStats) {
    for subscriber := range s.subscribers {
        select {
        case subscriber <- stats:
        default:
            // Skip if subscriber is not ready
        }
    }
}
