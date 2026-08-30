package gamificationapi

import (
	"net/http"

	authn "github.com/noirbyss/worktrition-app/backend/api-gateway-service/internal/auth"
	"github.com/noirbyss/worktrition-app/backend/api-gateway-service/internal/httpx"
	gamificationpb "github.com/noirbyss/worktrition-app/gen/gamification-service"
)

type applyWorkoutRewardRequest struct {
	IsStrength bool `json:"is_strength"`
}

func (h *Handler) HandleGetCharacter(w http.ResponseWriter, r *http.Request) {
	if !httpx.RequireMethod(w, r, http.MethodGet) {
		return
	}

	ctx, cancel := h.grpcContext(r)
	defer cancel()

	resp, err := h.gamificationClient.GetCharacter(ctx, &gamificationpb.GetCharacterRequest{
		UserId: authn.UserIDFromContext(r.Context()),
	})
	if err != nil {
		httpx.WriteGRPCError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) HandleApplyWorkoutReward(w http.ResponseWriter, r *http.Request) {
	if !httpx.RequireMethod(w, r, http.MethodPost) {
		return
	}

	var req applyWorkoutRewardRequest
	if !httpx.DecodeJSONBody(w, r, &req) {
		return
	}

	ctx, cancel := h.grpcContext(r)
	defer cancel()

	resp, err := h.gamificationClient.ApplyWorkoutReward(ctx, &gamificationpb.ApplyWorkoutRewardRequest{
		UserId:     authn.UserIDFromContext(r.Context()),
		IsStrength: req.IsStrength,
	})
	if err != nil {
		httpx.WriteGRPCError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) HandleApplyMealReward(w http.ResponseWriter, r *http.Request) {
	if !httpx.RequireMethod(w, r, http.MethodPost) {
		return
	}

	ctx, cancel := h.grpcContext(r)
	defer cancel()

	resp, err := h.gamificationClient.ApplyMealReward(ctx, &gamificationpb.ApplyMealRewardRequest{
		UserId: authn.UserIDFromContext(r.Context()),
	})
	if err != nil {
		httpx.WriteGRPCError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) HandleApplyWaterReward(w http.ResponseWriter, r *http.Request) {
	if !httpx.RequireMethod(w, r, http.MethodPost) {
		return
	}

	ctx, cancel := h.grpcContext(r)
	defer cancel()

	resp, err := h.gamificationClient.ApplyWaterReward(ctx, &gamificationpb.ApplyWaterRewardRequest{
		UserId: authn.UserIDFromContext(r.Context()),
	})
	if err != nil {
		httpx.WriteGRPCError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, resp)
}
