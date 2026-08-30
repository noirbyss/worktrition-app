package userapi

import (
	"net/http"

	authn "github.com/noirbyss/worktrition-app/backend/api-gateway-service/internal/auth"
	"github.com/noirbyss/worktrition-app/backend/api-gateway-service/internal/httpx"
	userpb "github.com/noirbyss/worktrition-app/gen/user-service"
)

func (h *Handler) HandleGetCurrentUser(w http.ResponseWriter, r *http.Request) {
	if !httpx.RequireMethod(w, r, http.MethodGet) {
		return
	}

	ctx, cancel := h.grpcContext(r)
	defer cancel()

	resp, err := h.userClient.GetUser(ctx, &userpb.GetUserRequest{
		UserId: authn.UserIDFromContext(r.Context()),
	})
	if err != nil {
		httpx.WriteGRPCError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) HandleProfile(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGetProfile(w, r)
	case http.MethodPut:
		h.handleSaveProfile(w, r)
	default:
		httpx.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) handleGetProfile(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := h.grpcContext(r)
	defer cancel()

	resp, err := h.userClient.GetProfile(ctx, &userpb.GetProfileRequest{
		UserId: authn.UserIDFromContext(r.Context()),
	})
	if err != nil {
		httpx.WriteGRPCError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) handleSaveProfile(w http.ResponseWriter, r *http.Request) {
	var req userpb.SaveProfileRequest
	if !httpx.DecodeJSONBody(w, r, &req) {
		return
	}
	req.UserId = authn.UserIDFromContext(r.Context())

	ctx, cancel := h.grpcContext(r)
	defer cancel()

	resp, err := h.userClient.SaveProfile(ctx, &req)
	if err != nil {
		httpx.WriteGRPCError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, resp)
}
