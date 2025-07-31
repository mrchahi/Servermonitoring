package model

import "time"

// LogEntry represents a single log entry
type LogEntry struct {
    Timestamp   time.Time `json:"timestamp"`
    Source      string    `json:"source"`      // auth.log, syslog, etc.
    Level       string    `json:"level"`       // info, warning, error, etc.
    Message     string    `json:"message"`
    User        string    `json:"user,omitempty"`
    IP          string    `json:"ip,omitempty"`
    ProcessName string    `json:"processName,omitempty"`
    ProcessID   int       `json:"processId,omitempty"`
}

// LogFilter represents log filtering options
type LogFilter struct {
    Source    string    `json:"source"`
    Level     string    `json:"level"`
    StartTime time.Time `json:"startTime"`
    EndTime   time.Time `json:"endTime"`
    Search    string    `json:"search"`
    Limit     int       `json:"limit"`
}

// LogSummary represents a summary of log statistics
type LogSummary struct {
    TotalEntries   int            `json:"totalEntries"`
    ErrorCount     int            `json:"errorCount"`
    WarningCount   int            `json:"warningCount"`
    SourceCounts   map[string]int `json:"sourceCounts"`
    RecentErrors   []LogEntry     `json:"recentErrors"`
    LastUpdateTime time.Time      `json:"lastUpdateTime"`
}
