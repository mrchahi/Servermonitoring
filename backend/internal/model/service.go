package model

// Service represents a system service
type Service struct {
    Name        string `json:"name"`
    DisplayName string `json:"displayName"` // نام فارسی و قابل فهم
    Status      string `json:"status"`      // active, inactive, failed
    Port        int    `json:"port,omitempty"`
    Description string `json:"description"` // توضیحات فارسی
    AutoStart   bool   `json:"autoStart"`   // enabled/disabled on boot
}

// ServiceAction represents an action that can be performed on a service
type ServiceAction struct {
    Action string `json:"action"` // start, stop, restart, enable, disable
    Name   string `json:"name"`
}

// ServiceStatus represents service operation result
type ServiceStatus struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
    Error   string `json:"error,omitempty"`
}
