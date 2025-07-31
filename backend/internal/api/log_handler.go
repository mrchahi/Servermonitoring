package api

import (
    "net/http"
    "path/filepath"
    "github.com/gin-gonic/gin"
    "github.com/mrchahi/Servermonitoring/internal/model"
    "github.com/mrchahi/Servermonitoring/internal/service"
)

type LogHandler struct {
    logManager *service.LogManager
}

func NewLogHandler(lm *service.LogManager) *LogHandler {
    return &LogHandler{
        logManager: lm,
    }
}

// GetLogs handles GET /api/logs
func (h *LogHandler) GetLogs(c *gin.Context) {
    var filter model.LogFilter
    if err := c.ShouldBindQuery(&filter); err != nil {
        c.JSON(400, gin.H{
            "error": "پارامترهای نامعتبر",
        })
        return
    }

    logs, err := h.logManager.GetLogs(filter)
    if err != nil {
        c.JSON(500, gin.H{
            "error": "خطا در دریافت لاگ‌ها: " + err.Error(),
        })
        return
    }

    c.JSON(200, logs)
}

// GetLogSummary handles GET /api/logs/summary
func (h *LogHandler) GetLogSummary(c *gin.Context) {
    summary, err := h.logManager.GetSummary()
    if err != nil {
        c.JSON(500, gin.H{
            "error": "خطا در دریافت خلاصه لاگ‌ها: " + err.Error(),
        })
        return
    }

    c.JSON(200, summary)
}

// DownloadLog handles GET /api/logs/:source/download
func (h *LogHandler) DownloadLog(c *gin.Context) {
    source := c.Param("source")
    
    filePath, err := h.logManager.DownloadLog(source)
    if err != nil {
        c.JSON(500, gin.H{
            "error": "خطا در آماده‌سازی فایل لاگ: " + err.Error(),
        })
        return
    }

    filename := filepath.Base(filePath)
    c.Header("Content-Description", "File Transfer")
    c.Header("Content-Disposition", "attachment; filename="+filename)
    c.Header("Content-Type", "text/plain")
    c.File(filePath)
}
