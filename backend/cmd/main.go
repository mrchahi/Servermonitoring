package main

import (
    "log"
    "github.com/mrchahi/Servermonitoring/internal/api"
    "github.com/mrchahi/Servermonitoring/internal/config"
)

func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    // Initialize server
    server := api.NewServer(cfg)
    
    // Start server
    if err := server.Start(); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}
