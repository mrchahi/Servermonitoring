package api

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/mrchahi/Servermonitoring/internal/config"
    "github.com/mrchahi/Servermonitoring/internal/service"
)

type Server struct {
    router            *gin.Engine
    config           *config.Config
    monitoringService *service.MonitoringService
    serviceManager    *service.ServiceManager
    firewallManager   *service.FirewallManager
    logManager        *service.LogManager
    wsHandler        *WebSocketHandler
    serviceHandler    *ServiceHandler
    firewallHandler   *FirewallHandler
    logHandler        *LogHandler
}

func NewServer(cfg *config.Config) *Server {
    router := gin.Default()
    monitoringService := service.NewMonitoringService(time.Duration(cfg.Server.UpdateInterval) * time.Second)
    serviceManager := service.NewServiceManager()
    firewallManager := service.NewFirewallManager()
    logManager := service.NewLogManager()
    wsHandler := NewWebSocketHandler(monitoringService)
    serviceHandler := NewServiceHandler(serviceManager)
    firewallHandler := NewFirewallHandler(firewallManager)
    logHandler := NewLogHandler(logManager)

    server := &Server{
        router:            router,
        config:           cfg,
        monitoringService: monitoringService,
        serviceManager:    serviceManager,
        firewallManager:   firewallManager,
        logManager:        logManager,
        wsHandler:        wsHandler,
        serviceHandler:    serviceHandler,
        firewallHandler:   firewallHandler,
        logHandler:        logHandler,
    }

    server.setupRoutes()
    return server
}

func (s *Server) setupRoutes() {
    // Middleware for CORS and security
    s.router.Use(s.corsMiddleware())
    s.router.Use(s.securityMiddleware())

    // WebSocket endpoint for real-time monitoring
    s.router.GET("/ws/stats", s.wsHandler.HandleWebSocket)

    // API endpoints
    api := s.router.Group("/api")
    {
        api.GET("/health", s.handleHealth)
        
        // Service management endpoints
        api.GET("/services", s.serviceHandler.ListServices)
        api.POST("/services/:name/action", s.serviceHandler.ControlService)

        // Port and firewall management endpoints
        api.GET("/ports", s.firewallHandler.ListPorts)
        api.GET("/firewall/rules", s.firewallHandler.ListRules)
        api.POST("/firewall/rules", s.firewallHandler.AddRule)
        api.DELETE("/firewall/rules/:id", s.firewallHandler.DeleteRule)
        api.POST("/firewall/enable", s.firewallHandler.EnableFirewall)
        api.POST("/firewall/disable", s.firewallHandler.DisableFirewall)

        // Log management endpoints
        api.GET("/logs", s.logHandler.GetLogs)
        api.GET("/logs/download", s.logHandler.DownloadLogs)
        api.GET("/logs/stats", s.logHandler.GetLogStats)
        api.POST("/logs/filter", s.logHandler.FilterLogs)
    }
}

func (s *Server) Start() error {
    // Start the monitoring service
    s.monitoringService.Start()

    // Start the HTTP server
    addr := fmt.Sprintf(":%d", s.config.Server.Port)
    return s.router.Run(addr)
}

func (s *Server) corsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    }
}

func (s *Server) securityMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // IP whitelist check
        if len(s.config.Security.AllowedIPs) > 0 {
            clientIP := c.ClientIP()
            allowed := false
            for _, ip := range s.config.Security.AllowedIPs {
                if ip == clientIP {
                    allowed = true
                    break
                }
            }
            if !allowed {
                c.AbortWithStatus(403)
                return
            }
        }
        
        c.Next()
    }
}

func (s *Server) handleHealth(c *gin.Context) {
    c.JSON(200, gin.H{
        "status": "ok",
    })
}
