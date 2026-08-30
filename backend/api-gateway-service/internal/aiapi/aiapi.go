package aiapi

import (
	"context"
	"net/http"
	"time"

	aipb "github.com/noirbyss/worktrition-app/gen/ai-service"
)

const defaultRequestTimeout = 15 * time.Second

type Config struct {
	RequestTimeout time.Duration
}

type Handler struct {
	aiClient       aipb.AiServiceClient
	requestTimeout time.Duration
}

func New(aiClient aipb.AiServiceClient, cfg Config) *Handler {
	requestTimeout := cfg.RequestTimeout
	if requestTimeout <= 0 {
		requestTimeout = defaultRequestTimeout
	}

	return &Handler{
		aiClient:       aiClient,
		requestTimeout: requestTimeout,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, authMiddleware func(http.Handler) http.Handler) {
	mux.Handle("/ai/generations", authMiddleware(http.HandlerFunc(h.HandleGenerations)))
	mux.Handle("/ai/generations/", authMiddleware(http.HandlerFunc(h.HandleGenerationStatus)))
}

func (h *Handler) grpcContext(r *http.Request) (context.Context, context.CancelFunc) {
	return context.WithTimeout(r.Context(), h.requestTimeout)
}
