package api

import (
    "github.com/gin-gonic/gin"
    "github.com/mrchahi/Servermonitoring/internal/model"
    "github.com/mrchahi/Servermonitoring/internal/service"
)

type ServiceHandler struct {
    serviceManager *service.ServiceManager
}

func NewServiceHandler(sm *service.ServiceManager) *ServiceHandler {
    return &ServiceHandler{
        serviceManager: sm,
    }
}

// ListServices handles GET /api/services
func (h *ServiceHandler) ListServices(c *gin.Context) {
    services, err := h.serviceManager.ListServices()
    if err != nil {
        c.JSON(500, gin.H{
            "error": "خطا در دریافت لیست سرویس‌ها: " + err.Error(),
        })
        return
    }

    c.JSON(200, services)
}

// ControlService handles POST /api/services/:name/action
func (h *ServiceHandler) ControlService(c *gin.Context) {
    var action model.ServiceAction
    if err := c.ShouldBindJSON(&action); err != nil {
        c.JSON(400, gin.H{
            "error": "پارامترهای نامعتبر",
        })
        return
    }

    action.Name = c.Param("name")
    result := h.serviceManager.ControlService(action)
    
    if !result.Success {
        c.JSON(500, result)
        return
    }

    c.JSON(200, result)
}
