package api

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
    "github.com/mrchahi/Servermonitoring/internal/service"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        // TODO: Add proper origin checking in production
        return true
    },
}

type WebSocketHandler struct {
    monitoringService *service.MonitoringService
}

func NewWebSocketHandler(ms *service.MonitoringService) *WebSocketHandler {
    return &WebSocketHandler{
        monitoringService: ms,
    }
}

func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        return
    }
    defer conn.Close()

    // Subscribe to system stats updates
    statsChan := h.monitoringService.Subscribe()
    defer h.monitoringService.Unsubscribe(statsChan)

    // Send updates to the client
    for stats := range statsChan {
        if err := conn.WriteJSON(stats); err != nil {
            break
        }
    }
}
