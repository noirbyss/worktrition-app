package userapi

import (
	"fmt"
	"net/http"
	"strconv"

	authn "github.com/noirbyss/worktrition-app/backend/api-gateway-service/internal/auth"
	"github.com/noirbyss/worktrition-app/backend/api-gateway-service/internal/httpx"
	userpb "github.com/noirbyss/worktrition-app/gen/user-service"
)

type saveWeightMeasurementRequest struct {
	WeightKG float64 `json:"weight_kg"`
}

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

func (h *Handler) HandleWeightMeasurement(w http.ResponseWriter, r *http.Request) {
	if !httpx.RequireMethod(w, r, http.MethodPost) {
		return
	}

	var req saveWeightMeasurementRequest
	if !httpx.DecodeJSONBody(w, r, &req) {
		return
	}

	ctx, cancel := h.grpcContext(r)
	defer cancel()

	resp, err := h.userClient.SaveWeightMeasurement(ctx, &userpb.SaveWeightMeasurementRequest{
		UserId:   authn.UserIDFromContext(r.Context()),
		WeightKg: req.WeightKG,
	})
	if err != nil {
		httpx.WriteGRPCError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) HandleWeightHistory(w http.ResponseWriter, r *http.Request) {
	if !httpx.RequireMethod(w, r, http.MethodGet) {
		return
	}

	limit, err := parseWeightHistoryLimit(r.URL.Query().Get("limit"))
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := h.grpcContext(r)
	defer cancel()

	resp, err := h.userClient.GetWeightHistory(ctx, &userpb.GetWeightHistoryRequest{
		UserId: authn.UserIDFromContext(r.Context()),
		Limit:  int32(limit),
	})
	if err != nil {
		httpx.WriteGRPCError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, resp)
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

func parseWeightHistoryLimit(raw string) (int, error) {
	if raw == "" {
		return 14, nil
	}

	limit, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("limit must be a positive integer")
	}

	return limit, nil
}
