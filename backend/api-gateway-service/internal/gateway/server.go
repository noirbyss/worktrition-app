package gateway

import (
	"net/http"
	"time"

	authn "github.com/noirbyss/worktrition-app/backend/api-gateway-service/internal/auth"
	"github.com/noirbyss/worktrition-app/backend/api-gateway-service/internal/config"
	"github.com/noirbyss/worktrition-app/backend/api-gateway-service/internal/httpx"
)

type Server struct {
	mux           *http.ServeMux
	allowedOrigin string
}

func New(cfg *config.Config, registrars ...RouteRegistrar) *Server {
	server := &Server{
		mux:           http.NewServeMux(),
		allowedOrigin: cfg.AllowedOrigin,
	}

	authMiddleware := authn.Middleware([]byte(cfg.UserJWTSecret), time.Now)

	server.routes(authMiddleware, registrars...)

	return server
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	httpx.CORS(s.allowedOrigin, s.mux).ServeHTTP(w, r)
}
