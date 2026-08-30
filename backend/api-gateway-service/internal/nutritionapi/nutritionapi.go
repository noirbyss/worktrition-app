package nutritionapi

import (
	"context"
	"net/http"
	"time"

	nutritionpb "github.com/noirbyss/worktrition-app/gen/nutrition-service"
)

const defaultRequestTimeout = 15 * time.Second

type Config struct {
	RequestTimeout time.Duration
}

type Handler struct {
	nutritionClient nutritionpb.NutritionServiceClient
	requestTimeout  time.Duration
}

func New(nutritionClient nutritionpb.NutritionServiceClient, cfg Config) *Handler {
	requestTimeout := cfg.RequestTimeout
	if requestTimeout <= 0 {
		requestTimeout = defaultRequestTimeout
	}

	return &Handler{
		nutritionClient: nutritionClient,
		requestTimeout:  requestTimeout,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, authMiddleware func(http.Handler) http.Handler) {
	mux.Handle("/nutrition/plan", authMiddleware(http.HandlerFunc(h.HandlePlan)))
	mux.Handle("/nutrition/meals/complete", authMiddleware(http.HandlerFunc(h.HandleCompleteMeal)))
	mux.Handle("/nutrition/water/complete", authMiddleware(http.HandlerFunc(h.HandleCompleteWater)))
	mux.Handle("/nutrition/stats", authMiddleware(http.HandlerFunc(h.HandleGetStats)))
}

func (h *Handler) grpcContext(r *http.Request) (context.Context, context.CancelFunc) {
	return context.WithTimeout(r.Context(), h.requestTimeout)
}
