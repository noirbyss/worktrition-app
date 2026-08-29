package userapi

import (
	"net/http"

	userpb "github.com/noirbyss/worktrition-app/gen/user-service"
)

type Config struct {
	RefreshTokenCookieName string
	SecureCookies          bool
}

type Handler struct {
	userClient             userpb.UserServiceClient
	refreshTokenCookieName string
	secureCookies          bool
}

func New(userClient userpb.UserServiceClient, cfg Config) *Handler {
	return &Handler{
		userClient:             userClient,
		refreshTokenCookieName: cfg.RefreshTokenCookieName,
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
}
