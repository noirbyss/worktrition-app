package main

import (
	"log/slog"
	"os"

	"github.com/noirbyss/worktrition-app/backend/user-service/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		slog.Error("user-service failed", "error", err)
		os.Exit(1)
	}
}
