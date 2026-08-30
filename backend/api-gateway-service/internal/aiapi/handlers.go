package aiapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	authn "github.com/noirbyss/worktrition-app/backend/api-gateway-service/internal/auth"
	"github.com/noirbyss/worktrition-app/backend/api-gateway-service/internal/httpx"
	aipb "github.com/noirbyss/worktrition-app/gen/ai-service"
)

const generationPathPrefix = "/ai/generations/"

type planTypeValue struct {
	value aipb.PlanType
}

type startGenerationRequest struct {
	PlanType planTypeValue `json:"plan_type"`
}

type generationResponse struct {
	GenerationID string `json:"generation_id"`
	Status       string `json:"status"`
	ErrorMessage string `json:"error_message,omitempty"`
}

func (h *Handler) HandleGenerations(w http.ResponseWriter, r *http.Request) {
	if !httpx.RequireMethod(w, r, http.MethodPost) {
		return
	}

	var req startGenerationRequest
	if !httpx.DecodeJSONBody(w, r, &req) {
		return
	}

	ctx, cancel := h.grpcContext(r)
	defer cancel()

	resp, err := h.aiClient.StartGeneration(ctx, &aipb.StartGenerationRequest{
		UserId:   authn.UserIDFromContext(r.Context()),
		PlanType: req.PlanType.value,
	})
	if err != nil {
		httpx.WriteGRPCError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusAccepted, generationResponse{
		GenerationID: resp.GetGenerationId(),
		Status:       normalizeGenerationStatus(resp.GetStatus()),
	})
}

func (h *Handler) HandleGenerationStatus(w http.ResponseWriter, r *http.Request) {
	if !httpx.RequireMethod(w, r, http.MethodGet) {
		return
	}

	generationID := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, generationPathPrefix))
	if generationID == "" || strings.Contains(generationID, "/") {
		httpx.WriteError(w, http.StatusNotFound, "generation not found")
		return
	}

	ctx, cancel := h.grpcContext(r)
	defer cancel()

	resp, err := h.aiClient.GetGenerationStatus(ctx, &aipb.GetGenerationStatusRequest{
		GenerationId: generationID,
	})
	if err != nil {
		httpx.WriteGRPCError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, generationResponse{
		GenerationID: resp.GetGenerationId(),
		Status:       normalizeGenerationStatus(resp.GetStatus()),
		ErrorMessage: resp.GetErrorMessage(),
	})
}

func (p *planTypeValue) UnmarshalJSON(data []byte) error {
	var textValue string
	if err := json.Unmarshal(data, &textValue); err == nil {
		planType, parseErr := parsePlanType(textValue)
		if parseErr != nil {
			return parseErr
		}

		p.value = planType
		return nil
	}

	var numericValue int32
	if err := json.Unmarshal(data, &numericValue); err == nil {
		planType, parseErr := parsePlanType(strconv.FormatInt(int64(numericValue), 10))
		if parseErr != nil {
			return parseErr
		}

		p.value = planType
		return nil
	}

	return fmt.Errorf("plan_type must be a string or integer")
}

func parsePlanType(raw string) (aipb.PlanType, error) {
	switch normalized := strings.ToLower(strings.TrimSpace(raw)); normalized {
	case "", "0", "plan_type_unspecified", "unspecified":
		return aipb.PlanType_PLAN_TYPE_UNSPECIFIED, nil
	case "1", "all", "plan_type_all":
		return aipb.PlanType_PLAN_TYPE_ALL, nil
	case "2", "workout", "training", "plan_type_workout":
		return aipb.PlanType_PLAN_TYPE_WORKOUT, nil
	case "3", "nutrition", "meal", "plan_type_nutrition":
		return aipb.PlanType_PLAN_TYPE_NUTRITION, nil
	default:
		return aipb.PlanType_PLAN_TYPE_UNSPECIFIED, fmt.Errorf("unsupported plan_type %q", raw)
	}
}

func normalizeGenerationStatus(status aipb.GenerationStatus) string {
	switch status {
	case aipb.GenerationStatus_GENERATION_STATUS_PENDING:
		return "pending"
	case aipb.GenerationStatus_GENERATION_STATUS_DONE:
		return "done"
	case aipb.GenerationStatus_GENERATION_STATUS_FAILED:
		return "failed"
	default:
		return "unspecified"
	}
}
