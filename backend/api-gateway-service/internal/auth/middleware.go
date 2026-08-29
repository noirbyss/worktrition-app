package authn

import (
	"net/http"
	"strings"
	"time"

	"github.com/noirbyss/worktrition-app/backend/api-gateway-service/internal/httpx"
)

func Middleware(secret []byte, now func() time.Time) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
			scheme, token, ok := strings.Cut(authHeader, " ")
			if !ok || !strings.EqualFold(scheme, "Bearer") || strings.TrimSpace(token) == "" {
				httpx.WriteError(w, http.StatusUnauthorized, "authorization bearer token is required")
				return
			}

			userID, err := verifyAccessToken(strings.TrimSpace(token), secret, now())
			if err != nil {
				httpx.WriteError(w, http.StatusUnauthorized, err.Error())
				return
			}

			next.ServeHTTP(w, r.WithContext(contextWithUserID(r.Context(), userID)))
		})
	}
}
