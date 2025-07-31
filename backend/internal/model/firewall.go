package model

// Port represents a network port configuration
type Port struct {
    Number      int      `json:"number"`
    Protocol    string   `json:"protocol"` // tcp/udp
    Status      string   `json:"status"`   // open/closed
    Service     string   `json:"service"`  // service name using this port
    Description string   `json:"description"`
    AllowedIPs  []string `json:"allowedIPs,omitempty"`
}

// FirewallRule represents a firewall rule
type FirewallRule struct {
    ID          int      `json:"id"`
    Action      string   `json:"action"`      // allow/deny
    Protocol    string   `json:"protocol"`    // tcp/udp/any
    Port        int      `json:"port"`
    Source      string   `json:"source"`      // IP/CIDR
    Description string   `json:"description"` // توضیحات فارسی
    Enabled     bool     `json:"enabled"`
}

// FirewallRuleRequest represents a request to add/modify a firewall rule
type FirewallRuleRequest struct {
    Action      string   `json:"action"`      // allow/deny
    Protocol    string   `json:"protocol"`    // tcp/udp/any
    Port        int      `json:"port"`
    Source      string   `json:"source"`      // IP/CIDR
    Description string   `json:"description"`
}

// OperationResult represents the result of a firewall operation
type OperationResult struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
    Error   string `json:"error,omitempty"`
}
