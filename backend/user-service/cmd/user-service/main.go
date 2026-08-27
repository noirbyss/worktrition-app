package main

import (
	"log"

	"github.com/noirbyss/worktrition-app/backend/user-service/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	log.Printf(
		"user-service started: environment=%s grpc_address=%s log_level=%s",
		cfg.Environment,
		cfg.GRPC.Address(),
		cfg.Log.Level,
	)
}
