package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ApiKey               string
	AiBaseURL            string
	AiModel              string
	GRPCPort             string
	LogLevel             string
	DatabaseURL          string
	UserServiceAddr      string
	NutritionServiceAddr string
	WorkoutServiceAddr   string
}

func Load() (*Config, error) {
	_ = godotenv.Load(".env")

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, errors.New("OPENAI_API_KEY environment variable is required")
	}

	return &Config{
		ApiKey:               apiKey,
		AiBaseURL:            getEnv("OPENAI_BASE_URL", "https://openrouter.ai/api/v1"),
		AiModel:              getEnv("OPENAI_MODEL", "openai/gpt-4o-mini"),
		GRPCPort:             getEnvAny([]string{"AI_GRPC_PORT", "GRPC_PORT"}, "50056"),
		LogLevel:             getEnvAny([]string{"AI_LOG_LEVEL", "LOG_LEVEL"}, "WARN"),
		DatabaseURL:          os.Getenv("DATABASE_URL"),
		UserServiceAddr:      getEnv("USER_SERVICE_ADDR", "localhost:50051"),
		NutritionServiceAddr: getEnv("NUTRITION_SERVICE_ADDR", "localhost:50052"),
		WorkoutServiceAddr:   getEnv("WORKOUT_SERVICE_ADDR", "localhost:50054"),
	}, nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvAny(keys []string, fallback string) string {
	for _, key := range keys {
		if val := os.Getenv(key); val != "" {
			return val
		}
	}

	return fallback
}
