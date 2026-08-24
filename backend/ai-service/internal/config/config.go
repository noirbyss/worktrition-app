package config

import (
	"errors"
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	ApiKey string
	AiBaseURL string
	AiModel string
	GRPCPort string
	LogLevel string
	
}

func Load() (*Config, error) {
	_ = godotenv.Load(".env")
	
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, errors.New("OPENAI_API_KEY environment variable is required")
	}
	
	return &Config{
		ApiKey: apiKey,
		AiBaseURL: getEnv("OPENAI_BASE_URL", "https://openrouter.ai.api/v1"),
		AiModel: getEnv("OPENAI_MODEL", "openai/gpt-4o-mini"),
		GRPCPort: getEnv("GRPC_PORT", "50052"),
		LogLevel: getEnv("LOG_LEVEL", "WARN"),
	}, nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
	
}