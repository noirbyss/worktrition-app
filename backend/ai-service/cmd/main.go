package main

import (
	"ai-service/internal/app"
	"ai-service/internal/config"
	"ai-service/internal/logger"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	zapLog, closeLog, err := logger.NewLogger(cfg.LogLevel)
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer closeLog()

	application := app.New(cfg, zapLog)
	if err := application.Run(); err != nil {
		log.Fatalf("failed to run application: %v", err)
	}
}
