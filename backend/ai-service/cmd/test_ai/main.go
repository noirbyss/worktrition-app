package main

import (
    "ai-service/internal/config"
    "ai-service/internal/logger"
    "ai-service/internal/provider"
    "context"
    "encoding/json"
    "fmt"
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

    prov := provider.NewOpenRouterProvider(cfg.ApiKey, cfg.AiBaseURL, cfg.AiModel, zapLog)

    fmt.Println("Запрос к ИИ")

    resp, err := prov.GeneratePlan(context.Background(), provider.SystemPrompt, provider.UserPrompt)
    if err != nil {
        log.Fatalf("AI error: %v", err)
    }

    fmt.Println("Ответ получен")

    var result map[string]any
    if err := json.Unmarshal([]byte(resp), &result); err != nil {
        fmt.Println(resp) 
        return
    }

    out, _ := json.MarshalIndent(result, "", "  ")
    fmt.Println(string(out))
}