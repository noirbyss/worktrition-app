package userapi

import (
	"net/http"

	"github.com/noirbyss/worktrition-app/backend/api-gateway-service/internal/httpx"
	userpb "github.com/noirbyss/worktrition-app/gen/user-service"
)

type authResponse struct {
	UserID               string `json:"user_id"`
	ProfileCompleted     bool   `json:"profile_completed"`
	AccessToken          string `json:"access_token"`
	AccessTokenExpiresAt int64  `json:"access_token_expires_at"`
}

func (h *Handler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	if !httpx.RequireMethod(w, r, http.MethodPost) {
		return
	}

	var req userpb.RegisterRequest
	if !httpx.DecodeJSONBody(w, r, &req) {
		return
	}

	resp, err := h.userClient.Register(r.Context(), &req)
	if err != nil {
		httpx.WriteGRPCError(w, err)
		return
	}

	h.setRefreshTokenCookie(w, resp.GetRefreshToken(), resp.GetRefreshTokenExpiresAt())
	httpx.WriteJSON(w, http.StatusOK, authResponse{
		UserID:               resp.GetUserId(),
		ProfileCompleted:     resp.GetProfileCompleted(),
		AccessToken:          resp.GetAccessToken(),
		AccessTokenExpiresAt: resp.GetAccessTokenExpiresAt(),
	})
}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if !httpx.RequireMethod(w, r, http.MethodPost) {
		return
	}

	var req userpb.LoginRequest
	if !httpx.DecodeJSONBody(w, r, &req) {
		return
	}

	resp, err := h.userClient.Login(r.Context(), &req)
	if err != nil {
		httpx.WriteGRPCError(w, err)
		return
	}

	h.setRefreshTokenCookie(w, resp.GetRefreshToken(), resp.GetRefreshTokenExpiresAt())
	httpx.WriteJSON(w, http.StatusOK, authResponse{
		UserID:               resp.GetUserId(),
		ProfileCompleted:     resp.GetProfileCompleted(),
		AccessToken:          resp.GetAccessToken(),
		AccessTokenExpiresAt: resp.GetAccessTokenExpiresAt(),
	})
}

func (h *Handler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	if !httpx.RequireMethod(w, r, http.MethodPost) {
		return
	}

	refreshToken, ok := h.refreshTokenFromCookie(w, r)
	if !ok {
		return
	}

	resp, err := h.userClient.RefreshToken(r.Context(), &userpb.RefreshTokenRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		httpx.WriteGRPCError(w, err)
		return
	}

	h.setRefreshTokenCookie(w, resp.GetRefreshToken(), resp.GetRefreshTokenExpiresAt())
	httpx.WriteJSON(w, http.StatusOK, authResponse{
		UserID:               resp.GetUserId(),
		ProfileCompleted:     resp.GetProfileCompleted(),
		AccessToken:          resp.GetAccessToken(),
		AccessTokenExpiresAt: resp.GetAccessTokenExpiresAt(),
	})
}

func (h *Handler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	if !httpx.RequireMethod(w, r, http.MethodPost) {
		return
	}

	refreshToken, ok := h.refreshTokenFromCookie(w, r)
	if !ok {
		return
	}

	if _, err := h.userClient.Logout(r.Context(), &userpb.LogoutRequest{
		RefreshToken: refreshToken,
	}); err != nil {
		httpx.WriteGRPCError(w, err)
		return
	}

	h.clearRefreshTokenCookie(w)
	httpx.WriteJSON(w, http.StatusOK, map[string]bool{"success": true})
}
