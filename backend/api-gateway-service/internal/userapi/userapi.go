package userapi

import (
	"context"
	"net/http"
	"time"

	userpb "github.com/noirbyss/worktrition-app/gen/user-service"
)

const defaultRequestTimeout = 15 * time.Second

type Config struct {
	RefreshTokenCookieName string
	RequestTimeout         time.Duration
	SecureCookies          bool
}

type Handler struct {
	userClient             userpb.UserServiceClient
	refreshTokenCookieName string
	requestTimeout         time.Duration
	secureCookies          bool
}

func New(userClient userpb.UserServiceClient, cfg Config) *Handler {
	requestTimeout := cfg.RequestTimeout
	if requestTimeout <= 0 {
		requestTimeout = defaultRequestTimeout
	}

	return &Handler{
		userClient:             userClient,
		refreshTokenCookieName: cfg.RefreshTokenCookieName,
		requestTimeout:         requestTimeout,
		secureCookies:          cfg.SecureCookies,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, authMiddleware func(http.Handler) http.Handler) {
	mux.HandleFunc("/auth/register", h.HandleRegister)
	mux.HandleFunc("/auth/login", h.HandleLogin)
	mux.HandleFunc("/auth/refresh", h.HandleRefresh)
	mux.HandleFunc("/auth/logout", h.HandleLogout)
	mux.Handle("/users/me", authMiddleware(http.HandlerFunc(h.HandleGetCurrentUser)))
	mux.Handle("/profile", authMiddleware(http.HandlerFunc(h.HandleProfile)))
	mux.Handle("/profile/weight", authMiddleware(http.HandlerFunc(h.HandleWeightMeasurement)))
	mux.Handle("/profile/weight-history", authMiddleware(http.HandlerFunc(h.HandleWeightHistory)))
}

func (h *Handler) grpcContext(r *http.Request) (context.Context, context.CancelFunc) {
	return context.WithTimeout(r.Context(), h.requestTimeout)
}
