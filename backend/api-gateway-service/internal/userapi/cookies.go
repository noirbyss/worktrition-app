package userapi

import (
	"net/http"
	"strings"
	"time"

	"github.com/noirbyss/worktrition-app/backend/api-gateway-service/internal/httpx"
)

func (h *Handler) refreshTokenFromCookie(w http.ResponseWriter, r *http.Request) (string, bool) {
	cookie, err := r.Cookie(h.refreshTokenCookieName)
	if err != nil || strings.TrimSpace(cookie.Value) == "" {
		httpx.WriteError(w, http.StatusUnauthorized, "refresh token cookie is required")
		return "", false
	}

	return cookie.Value, true
}

func (h *Handler) setRefreshTokenCookie(w http.ResponseWriter, value string, expiresAt int64) {
	expires := time.Unix(expiresAt, 0)
	maxAge := int(time.Until(expires).Seconds())
	if maxAge < 0 {
		maxAge = -1
	}

	http.SetCookie(w, &http.Cookie{
		Name:     h.refreshTokenCookieName,
		Value:    value,
		Path:     "/auth",
		Expires:  expires,
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   h.secureCookies,
		SameSite: http.SameSiteLaxMode,
	})
}

func (h *Handler) clearRefreshTokenCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     h.refreshTokenCookieName,
		Value:    "",
		Path:     "/auth",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   h.secureCookies,
		SameSite: http.SameSiteLaxMode,
	})
}
