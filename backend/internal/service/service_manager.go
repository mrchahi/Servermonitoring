package service

import (
    "fmt"
    "os/exec"
    "strings"
    "github.com/mrchahi/Servermonitoring/internal/model"
)

// ServiceManager handles system service operations
type ServiceManager struct{}

// NewServiceManager creates a new service manager
func NewServiceManager() *ServiceManager {
    return &ServiceManager{}
}

// ListServices returns all system services
func (sm *ServiceManager) ListServices() ([]model.Service, error) {
    cmd := exec.Command("systemctl", "list-units", "--type=service", "--all", "--no-pager", "--plain", "--no-legend")
    output, err := cmd.Output()
    if err != nil {
        return nil, fmt.Errorf("error listing services: %v", err)
    }

    services := []model.Service{}
    lines := strings.Split(string(output), "\n")

    for _, line := range lines {
        if line == "" {
            continue
        }

        fields := strings.Fields(line)
        if len(fields) < 4 {
            continue
        }

        name := strings.TrimSuffix(fields[0], ".service")
        status := sm.parseStatus(fields[3])
        
        // Get additional service info
        description := sm.getServiceDescription(name)
        autoStart := sm.isServiceEnabled(name)
        port := sm.getServicePort(name)

        service := model.Service{
            Name:        name,
            DisplayName: sm.getDisplayName(name),
            Status:      status,
            Port:        port,
            Description: description,
            AutoStart:   autoStart,
        }

        services = append(services, service)
    }

    return services, nil
}

// ControlService handles service control actions (start/stop/restart/enable/disable)
func (sm *ServiceManager) ControlService(action model.ServiceAction) model.ServiceStatus {
    var cmd *exec.Cmd
    switch action.Action {
    case "start", "stop", "restart":
        cmd = exec.Command("systemctl", action.Action, action.Name+".service")
    case "enable", "disable":
        cmd = exec.Command("systemctl", action.Action, action.Name+".service")
    default:
        return model.ServiceStatus{
            Success: false,
            Error:   "عملیات نامعتبر",
        }
    }

    if err := cmd.Run(); err != nil {
        return model.ServiceStatus{
            Success: false,
            Error:   fmt.Sprintf("خطا در اجرای دستور: %v", err),
        }
    }

    message := fmt.Sprintf("عملیات %s برای سرویس %s با موفقیت انجام شد", 
        sm.translateAction(action.Action), action.Name)
    
    return model.ServiceStatus{
        Success: true,
        Message: message,
    }
}

// Helper functions

func (sm *ServiceManager) parseStatus(status string) string {
    switch strings.ToLower(status) {
    case "running":
        return "active"
    case "dead":
        return "inactive"
    case "failed":
        return "failed"
    default:
        return "unknown"
    }
}

func (sm *ServiceManager) getServiceDescription(name string) string {
    cmd := exec.Command("systemctl", "show", name+".service", "--property=Description", "--no-pager")
    output, err := cmd.Output()
    if err != nil {
        return ""
    }
    return strings.TrimPrefix(string(output), "Description=")
}

func (sm *ServiceManager) isServiceEnabled(name string) bool {
    cmd := exec.Command("systemctl", "is-enabled", name+".service")
    output, _ := cmd.Output()
    return strings.TrimSpace(string(output)) == "enabled"
}

func (sm *ServiceManager) getServicePort(name string) int {
    // This is a simplified version. In real implementation,
    // you might want to use 'ss' or 'netstat' to get the actual port
    return 0
}

func (sm *ServiceManager) getDisplayName(name string) string {
    // Map common service names to Persian display names
    serviceNames := map[string]string{
        "nginx":     "وب سرور Nginx",
        "apache2":   "وب سرور Apache",
        "mysql":     "پایگاه داده MySQL",
        "postgresql": "پایگاه داده PostgreSQL",
        "mongodb":   "پایگاه داده MongoDB",
        "redis":     "Redis",
        "ssh":       "SSH",
        "ufw":       "فایروال UFW",
        "fail2ban":  "Fail2Ban",
    }

    if displayName, ok := serviceNames[name]; ok {
        return displayName
    }
    return name
}

func (sm *ServiceManager) translateAction(action string) string {
    actions := map[string]string{
        "start":   "شروع",
        "stop":    "توقف",
        "restart": "راه‌اندازی مجدد",
        "enable":  "فعال‌سازی",
        "disable": "غیرفعال‌سازی",
    }

    if translated, ok := actions[action]; ok {
        return translated
    }
    return action
}
