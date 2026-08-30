package gamificationapi

import (
	"context"
	"net/http"
	"time"

	gamificationpb "github.com/noirbyss/worktrition-app/gen/gamification-service"
)

const defaultRequestTimeout = 15 * time.Second

type Config struct {
	RequestTimeout time.Duration
}

type Handler struct {
	gamificationClient gamificationpb.GamificationServiceClient
	requestTimeout     time.Duration
}

func New(gamificationClient gamificationpb.GamificationServiceClient, cfg Config) *Handler {
	requestTimeout := cfg.RequestTimeout
	if requestTimeout <= 0 {
		requestTimeout = defaultRequestTimeout
	}

	return &Handler{
		gamificationClient: gamificationClient,
		requestTimeout:     requestTimeout,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, authMiddleware func(http.Handler) http.Handler) {
	mux.Handle("/gamification/character", authMiddleware(http.HandlerFunc(h.HandleGetCharacter)))
	mux.Handle("/gamification/rewards/workout", authMiddleware(http.HandlerFunc(h.HandleApplyWorkoutReward)))
	mux.Handle("/gamification/rewards/meal", authMiddleware(http.HandlerFunc(h.HandleApplyMealReward)))
	mux.Handle("/gamification/rewards/water", authMiddleware(http.HandlerFunc(h.HandleApplyWaterReward)))
}

func (h *Handler) grpcContext(r *http.Request) (context.Context, context.CancelFunc) {
	return context.WithTimeout(r.Context(), h.requestTimeout)
}
