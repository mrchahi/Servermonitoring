package api

import (
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/mrchahi/Servermonitoring/internal/model"
    "github.com/mrchahi/Servermonitoring/internal/service"
)

type FirewallHandler struct {
    firewallManager *service.FirewallManager
}

func NewFirewallHandler(fm *service.FirewallManager) *FirewallHandler {
    return &FirewallHandler{
        firewallManager: fm,
    }
}

// ListPorts handles GET /api/ports
func (h *FirewallHandler) ListPorts(c *gin.Context) {
    ports, err := h.firewallManager.ListPorts()
    if err != nil {
        c.JSON(500, gin.H{
            "error": "خطا در دریافت لیست پورت‌ها: " + err.Error(),
        })
        return
    }

    c.JSON(200, ports)
}

// ListRules handles GET /api/firewall/rules
func (h *FirewallHandler) ListRules(c *gin.Context) {
    rules, err := h.firewallManager.ListRules()
    if err != nil {
        c.JSON(500, gin.H{
            "error": "خطا در دریافت قوانین فایروال: " + err.Error(),
        })
        return
    }

    c.JSON(200, rules)
}

// AddRule handles POST /api/firewall/rules
func (h *FirewallHandler) AddRule(c *gin.Context) {
    var rule model.FirewallRuleRequest
    if err := c.ShouldBindJSON(&rule); err != nil {
        c.JSON(400, gin.H{
            "error": "پارامترهای نامعتبر",
        })
        return
    }

    result := h.firewallManager.AddRule(rule)
    if !result.Success {
        c.JSON(500, result)
        return
    }

    c.JSON(200, result)
}

// DeleteRule handles DELETE /api/firewall/rules/:id
func (h *FirewallHandler) DeleteRule(c *gin.Context) {
    ruleID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(400, gin.H{
            "error": "شناسه قانون نامعتبر است",
        })
        return
    }

    result := h.firewallManager.DeleteRule(ruleID)
    if !result.Success {
        c.JSON(500, result)
        return
    }

    c.JSON(200, result)
}

// EnableFirewall handles POST /api/firewall/enable
func (h *FirewallHandler) EnableFirewall(c *gin.Context) {
    result := h.firewallManager.EnableFirewall()
    if !result.Success {
        c.JSON(500, result)
        return
    }

    c.JSON(200, result)
}

// DisableFirewall handles POST /api/firewall/disable
func (h *FirewallHandler) DisableFirewall(c *gin.Context) {
    result := h.firewallManager.DisableFirewall()
    if !result.Success {
        c.JSON(500, result)
        return
    }

    c.JSON(200, result)
}
