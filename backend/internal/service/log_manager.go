package service

import (
    "bufio"
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "time"
    "regexp"
    "sync"
    "github.com/mrchahi/Servermonitoring/internal/model"
)

// LogManager handles system log operations
type LogManager struct {
    logPaths map[string]string
    mutex    sync.RWMutex
    cache    map[string][]model.LogEntry
    summary  *model.LogSummary
}

// NewLogManager creates a new log manager
func NewLogManager() *LogManager {
    // Common log file paths
    logPaths := map[string]string{
        "auth":     "/var/log/auth.log",
        "syslog":   "/var/log/syslog",
        "kern":     "/var/log/kern.log",
        "fail2ban": "/var/log/fail2ban.log",
    }

    return &LogManager{
        logPaths: logPaths,
        cache:    make(map[string][]model.LogEntry),
        summary: &model.LogSummary{
            SourceCounts: make(map[string]int),
            RecentErrors: make([]model.LogEntry, 0),
        },
    }
}

// GetLogs returns filtered log entries
func (lm *LogManager) GetLogs(filter model.LogFilter) ([]model.LogEntry, error) {
    lm.mutex.RLock()
    defer lm.mutex.RUnlock()

    var entries []model.LogEntry
    logPath, exists := lm.logPaths[filter.Source]
    if !exists {
        return nil, fmt.Errorf("منبع لاگ نامعتبر")
    }

    // Try to get from cache first
    if cachedEntries, ok := lm.cache[filter.Source]; ok {
        entries = lm.filterEntries(cachedEntries, filter)
    } else {
        // Read from file if not in cache
        var err error
        entries, err = lm.readLogFile(logPath, filter)
        if err != nil {
            return nil, err
        }
        // Cache the entries
        lm.cache[filter.Source] = entries
    }

    return entries, nil
}

// GetSummary returns log summary statistics
func (lm *LogManager) GetSummary() (*model.LogSummary, error) {
    lm.mutex.RLock()
    defer lm.mutex.RUnlock()

    // Return cached summary if it's recent enough (less than 5 minutes old)
    if lm.summary != nil && time.Since(lm.summary.LastUpdateTime) < 5*time.Minute {
        return lm.summary, nil
    }

    // Update summary
    summary := &model.LogSummary{
        SourceCounts: make(map[string]int),
        RecentErrors: make([]model.LogEntry, 0),
    }

    for source := range lm.logPaths {
        entries, err := lm.GetLogs(model.LogFilter{
            Source: source,
            Limit:  1000, // Last 1000 entries
        })
        if err != nil {
            continue
        }

        summary.TotalEntries += len(entries)
        summary.SourceCounts[source] = len(entries)

        for _, entry := range entries {
            switch strings.ToLower(entry.Level) {
            case "error", "err", "critical", "alert", "emergency":
                summary.ErrorCount++
                if len(summary.RecentErrors) < 10 {
                    summary.RecentErrors = append(summary.RecentErrors, entry)
                }
            case "warning", "warn":
                summary.WarningCount++
            }
        }
    }

    summary.LastUpdateTime = time.Now()
    lm.summary = summary
    return summary, nil
}

// DownloadLog prepares a log file for download
func (lm *LogManager) DownloadLog(source string) (string, error) {
    logPath, exists := lm.logPaths[source]
    if !exists {
        return "", fmt.Errorf("منبع لاگ نامعتبر")
    }

    // Create a temporary copy of the log file
    tmpDir := os.TempDir()
    tmpFile := filepath.Join(tmpDir, fmt.Sprintf("%s_%s.log", source, time.Now().Format("20060102_150405")))
    
    input, err := os.ReadFile(logPath)
    if err != nil {
        return "", err
    }

    err = os.WriteFile(tmpFile, input, 0644)
    if err != nil {
        return "", err
    }

    return tmpFile, nil
}

// Helper functions

func (lm *LogManager) readLogFile(path string, filter model.LogFilter) ([]model.LogEntry, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var entries []model.LogEntry
    scanner := bufio.NewScanner(file)
    
    // Common log timestamp formats
    timeFormats := []string{
        "Jan _2 15:04:05",
        "Jan _2 15:04:05 2006",
        "2006-01-02 15:04:05",
    }

    for scanner.Scan() {
        line := scanner.Text()
        if entry := lm.parseLine(line, timeFormats); entry != nil {
            entries = append(entries, *entry)
        }
    }

    if err := scanner.Err(); err != nil {
        return nil, err
    }

    return entries, nil
}

func (lm *LogManager) parseLine(line string, timeFormats []string) *model.LogEntry {
    // Basic log line pattern
    pattern := `^(\w{3}\s+\d{1,2}\s+\d{2}:\d{2}:\d{2})\s+(\S+)\s+([^:]+):\s+(.+)$`
    re := regexp.MustCompile(pattern)
    matches := re.FindStringSubmatch(line)

    if len(matches) < 5 {
        return nil
    }

    // Try parsing timestamp
    var timestamp time.Time
    var err error
    for _, format := range timeFormats {
        timestamp, err = time.Parse(format, matches[1])
        if err == nil {
            break
        }
    }

    if err != nil {
        return nil
    }

    // Determine log level
    level := "info"
    messageLower := strings.ToLower(matches[4])
    if strings.Contains(messageLower, "error") || strings.Contains(messageLower, "failed") {
        level = "error"
    } else if strings.Contains(messageLower, "warning") || strings.Contains(messageLower, "warn") {
        level = "warning"
    }

    // Extract process info
    processName := matches[2]
    processID := 0
    if pidMatch := regexp.MustCompile(`\[(\d+)\]`).FindStringSubmatch(matches[3]); len(pidMatch) > 1 {
        processID, _ = strconv.Atoi(pidMatch[1])
    }

    // Extract IP address if present
    var ip string
    if ipMatch := regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}\b`).FindString(matches[4]); ipMatch != "" {
        ip = ipMatch
    }

    return &model.LogEntry{
        Timestamp:   timestamp,
        Source:      matches[2],
        Level:       level,
        Message:     matches[4],
        ProcessName: processName,
        ProcessID:   processID,
        IP:         ip,
    }
}

func (lm *LogManager) filterEntries(entries []model.LogEntry, filter model.LogFilter) []model.LogEntry {
    filtered := make([]model.LogEntry, 0)

    for _, entry := range entries {
        // Apply time filter
        if !filter.StartTime.IsZero() && entry.Timestamp.Before(filter.StartTime) {
            continue
        }
        if !filter.EndTime.IsZero() && entry.Timestamp.After(filter.EndTime) {
            continue
        }

        // Apply level filter
        if filter.Level != "" && !strings.EqualFold(entry.Level, filter.Level) {
            continue
        }

        // Apply search filter
        if filter.Search != "" {
            searchLower := strings.ToLower(filter.Search)
            messageLower := strings.ToLower(entry.Message)
            if !strings.Contains(messageLower, searchLower) {
                continue
            }
        }

        filtered = append(filtered, entry)
    }

    // Apply limit
    if filter.Limit > 0 && len(filtered) > filter.Limit {
        filtered = filtered[len(filtered)-filter.Limit:]
    }

    return filtered
}
