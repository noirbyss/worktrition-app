package workoutapi

import (
	"context"
	"net/http"
	"time"

	workoutpb "github.com/noirbyss/worktrition-app/gen/workout-service"
)

const grpcRequestTimeout = 5 * time.Second

type Handler struct {
	workoutClient workoutpb.WorkoutServiceClient
}

func New(workoutClient workoutpb.WorkoutServiceClient) *Handler {
	return &Handler{
		workoutClient: workoutClient,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, authMiddleware func(http.Handler) http.Handler) {
	mux.Handle("/workout/plan", authMiddleware(http.HandlerFunc(h.HandlePlan)))
	mux.Handle("/workout/training/complete", authMiddleware(http.HandlerFunc(h.HandleCompleteTraining)))
	mux.Handle("/workout/stats", authMiddleware(http.HandlerFunc(h.HandleGetStats)))
}

func (h *Handler) grpcContext(r *http.Request) (context.Context, context.CancelFunc) {
	return context.WithTimeout(r.Context(), grpcRequestTimeout)
}
