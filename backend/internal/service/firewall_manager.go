package service

import (
    "fmt"
    "os/exec"
    "strings"
    "strconv"
    "regexp"
    "github.com/mrchahi/Servermonitoring/internal/model"
)

// FirewallManager handles firewall operations using UFW
type FirewallManager struct{}

// NewFirewallManager creates a new firewall manager
func NewFirewallManager() *FirewallManager {
    return &FirewallManager{}
}

// ListPorts returns all open ports and their status
func (fm *FirewallManager) ListPorts() ([]model.Port, error) {
    // Use ss command to get listening ports
    cmd := exec.Command("ss", "-tuln")
    output, err := cmd.Output()
    if err != nil {
        return nil, fmt.Errorf("error getting port list: %v", err)
    }

    return fm.parsePortList(string(output))
}

// ListRules returns all firewall rules
func (fm *FirewallManager) ListRules() ([]model.FirewallRule, error) {
    cmd := exec.Command("ufw", "status", "numbered")
    output, err := cmd.Output()
    if err != nil {
        return nil, fmt.Errorf("error getting firewall rules: %v", err)
    }

    return fm.parseFirewallRules(string(output))
}

// AddRule adds a new firewall rule
func (fm *FirewallManager) AddRule(rule model.FirewallRuleRequest) model.OperationResult {
    // Validate rule
    if err := fm.validateRule(rule); err != nil {
        return model.OperationResult{
            Success: false,
            Error:   err.Error(),
        }
    }

    // Construct UFW command
    args := []string{"allow"}
    if rule.Action == "deny" {
        args = []string{"deny"}
    }

    if rule.Source != "" {
        args = append(args, "from", rule.Source)
    }

    args = append(args, "to", "any", "port", strconv.Itoa(rule.Port))
    if rule.Protocol != "any" {
        args = append(args, "proto", rule.Protocol)
    }

    // Execute UFW command
    cmd := exec.Command("ufw", args...)
    if err := cmd.Run(); err != nil {
        return model.OperationResult{
            Success: false,
            Error:   fmt.Sprintf("خطا در اجرای دستور: %v", err),
        }
    }

    return model.OperationResult{
        Success: true,
        Message: "قانون فایروال با موفقیت اضافه شد",
    }
}

// DeleteRule deletes a firewall rule by its number
func (fm *FirewallManager) DeleteRule(ruleNumber int) model.OperationResult {
    cmd := exec.Command("ufw", "delete", strconv.Itoa(ruleNumber))
    if err := cmd.Run(); err != nil {
        return model.OperationResult{
            Success: false,
            Error:   fmt.Sprintf("خطا در حذف قانون: %v", err),
        }
    }

    return model.OperationResult{
        Success: true,
        Message: "قانون فایروال با موفقیت حذف شد",
    }
}

// EnableFirewall enables the firewall
func (fm *FirewallManager) EnableFirewall() model.OperationResult {
    cmd := exec.Command("ufw", "enable")
    if err := cmd.Run(); err != nil {
        return model.OperationResult{
            Success: false,
            Error:   fmt.Sprintf("خطا در فعال‌سازی فایروال: %v", err),
        }
    }

    return model.OperationResult{
        Success: true,
        Message: "فایروال با موفقیت فعال شد",
    }
}

// DisableFirewall disables the firewall
func (fm *FirewallManager) DisableFirewall() model.OperationResult {
    cmd := exec.Command("ufw", "disable")
    if err := cmd.Run(); err != nil {
        return model.OperationResult{
            Success: false,
            Error:   fmt.Sprintf("خطا در غیرفعال‌سازی فایروال: %v", err),
        }
    }

    return model.OperationResult{
        Success: true,
        Message: "فایروال با موفقیت غیرفعال شد",
    }
}

// Helper functions

func (fm *FirewallManager) parsePortList(output string) ([]model.Port, error) {
    lines := strings.Split(output, "\n")
    ports := []model.Port{}
    
    // Skip header line
    for _, line := range lines[1:] {
        if line == "" {
            continue
        }

        fields := strings.Fields(line)
        if len(fields) < 5 {
            continue
        }

        // Parse local address field (e.g., "*:80" or "127.0.0.1:80")
        addrParts := strings.Split(fields[4], ":")
        if len(addrParts) != 2 {
            continue
        }

        portNum, err := strconv.Atoi(addrParts[1])
        if err != nil {
            continue
        }

        protocol := strings.ToLower(strings.TrimPrefix(fields[0], "LISTEN"))
        
        port := model.Port{
            Number:   portNum,
            Protocol: protocol,
            Status:   "open",
            Service:  fm.getServiceNameForPort(portNum),
        }

        ports = append(ports, port)
    }

    return ports, nil
}

func (fm *FirewallManager) parseFirewallRules(output string) ([]model.FirewallRule, error) {
    lines := strings.Split(output, "\n")
    rules := []model.FirewallRule{}
    ruleRegex := regexp.MustCompile(`^\[ *(\d+)\] (.+)$`)

    for _, line := range lines {
        matches := ruleRegex.FindStringSubmatch(line)
        if len(matches) != 3 {
            continue
        }

        id, _ := strconv.Atoi(matches[1])
        ruleStr := matches[2]

        rule := fm.parseRuleString(id, ruleStr)
        if rule != nil {
            rules = append(rules, *rule)
        }
    }

    return rules, nil
}

func (fm *FirewallManager) parseRuleString(id int, ruleStr string) *model.FirewallRule {
    fields := strings.Fields(ruleStr)
    if len(fields) < 3 {
        return nil
    }

    rule := &model.FirewallRule{
        ID:      id,
        Action:  fields[0],
        Enabled: true,
    }

    // Parse remaining fields for port, protocol, and source
    for i := 1; i < len(fields); i++ {
        switch fields[i] {
        case "from":
            if i+1 < len(fields) {
                rule.Source = fields[i+1]
                i++
            }
        case "port":
            if i+1 < len(fields) {
                port, _ := strconv.Atoi(fields[i+1])
                rule.Port = port
                i++
            }
        case "proto":
            if i+1 < len(fields) {
                rule.Protocol = fields[i+1]
                i++
            }
        }
    }

    return rule
}

func (fm *FirewallManager) validateRule(rule model.FirewallRuleRequest) error {
    if rule.Action != "allow" && rule.Action != "deny" {
        return fmt.Errorf("عملیات نامعتبر")
    }

    if rule.Port < 1 || rule.Port > 65535 {
        return fmt.Errorf("پورت نامعتبر")
    }

    if rule.Protocol != "tcp" && rule.Protocol != "udp" && rule.Protocol != "any" {
        return fmt.Errorf("پروتکل نامعتبر")
    }

    return nil
}

func (fm *FirewallManager) getServiceNameForPort(port int) string {
    // Common port to service name mapping
    commonPorts := map[int]string{
        22:    "SSH",
        80:    "HTTP",
        443:   "HTTPS",
        3306:  "MySQL",
        5432:  "PostgreSQL",
        27017: "MongoDB",
        6379:  "Redis",
    }

    if service, ok := commonPorts[port]; ok {
        return service
    }

    return "نامشخص"
}
