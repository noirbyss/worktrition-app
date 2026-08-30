package workoutapi

import (
	"context"
	"net/http"
	"time"

	workoutpb "github.com/noirbyss/worktrition-app/gen/workout-service"
)

const defaultRequestTimeout = 15 * time.Second

type Config struct {
	RequestTimeout time.Duration
}

type Handler struct {
	workoutClient  workoutpb.WorkoutServiceClient
	requestTimeout time.Duration
}

func New(workoutClient workoutpb.WorkoutServiceClient, cfg Config) *Handler {
	requestTimeout := cfg.RequestTimeout
	if requestTimeout <= 0 {
		requestTimeout = defaultRequestTimeout
	}

	return &Handler{
		workoutClient:  workoutClient,
		requestTimeout: requestTimeout,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, authMiddleware func(http.Handler) http.Handler) {
	mux.Handle("/workout/plan", authMiddleware(http.HandlerFunc(h.HandlePlan)))
	mux.Handle("/workout/training/complete", authMiddleware(http.HandlerFunc(h.HandleCompleteTraining)))
	mux.Handle("/workout/stats", authMiddleware(http.HandlerFunc(h.HandleGetStats)))
}

func (h *Handler) grpcContext(r *http.Request) (context.Context, context.CancelFunc) {
	return context.WithTimeout(r.Context(), h.requestTimeout)
}
