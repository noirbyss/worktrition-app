package main

import (
	"log/slog"
	"os"

	"github.com/noirbyss/worktrition-app/backend/api-gateway-service/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		slog.Error("api-gateway-service stopped with error", "error", err)
		os.Exit(1)
	}
}
