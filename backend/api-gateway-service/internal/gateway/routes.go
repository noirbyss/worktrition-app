package gateway

import "net/http"

type Middleware = func(http.Handler) http.Handler

type RouteRegistrar interface {
	RegisterRoutes(mux *http.ServeMux, authMiddleware Middleware)
}

func (s *Server) routes(authMiddleware Middleware, registrars ...RouteRegistrar) {
	for _, registrar := range registrars {
		registrar.RegisterRoutes(s.mux, authMiddleware)
	}
}
