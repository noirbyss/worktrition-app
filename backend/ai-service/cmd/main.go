package main

import (
	"ai-service/internal/config"
	"ai-service/internal/logger"
	"ai-service/internal/provider"
	"ai-service/internal/usecase"
	"context"

	"github.com/k0kubun/pp/v3"
)

func main () {
	cfg, err := config.Load()
	if err != nil {
		panic("failed to load .env")
	}
	log,closelog,err := logger.NewLogger(cfg.LogLevel)
	if err != nil {
		panic("failed to create logger")
	}
	defer closelog()
	
	prov := provider.NewOpenRouterProvider(
			cfg.ApiKey,
			cfg.AiBaseURL,
			cfg.AiModel,
			log,
	)
	us := usecase.NewUseCase(prov,log)

	plan,err := us.StartGeneration(context.Background(),"","")

	pp.Print(plan)
}