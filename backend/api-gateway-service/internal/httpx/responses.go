package httpx

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type errorResponse struct {
	Error string `json:"error"`
}

func WriteGRPCError(w http.ResponseWriter, err error) {
	code := status.Code(err)
	message := status.Convert(err).Message()
	if message == "" {
		message = "internal error"
	}

	switch code {
	case codes.InvalidArgument:
		WriteError(w, http.StatusBadRequest, message)
	case codes.Unauthenticated:
		WriteError(w, http.StatusUnauthorized, message)
	case codes.NotFound:
		WriteError(w, http.StatusNotFound, message)
	case codes.AlreadyExists:
		WriteError(w, http.StatusConflict, message)
	case codes.DeadlineExceeded:
		WriteError(w, http.StatusGatewayTimeout, "upstream service request timed out")
	case codes.Unavailable:
		WriteError(w, http.StatusServiceUnavailable, "upstream service is unavailable")
	case codes.Internal, codes.Unknown:
		WriteError(w, http.StatusInternalServerError, message)
	default:
		WriteError(w, http.StatusInternalServerError, fmt.Sprintf("upstream service error: %s", message))
	}
}

func WriteError(w http.ResponseWriter, statusCode int, message string) {
	if message == "" {
		message = http.StatusText(statusCode)
	}

	WriteJSON(w, statusCode, errorResponse{Error: message})
}

func WriteJSON(w http.ResponseWriter, statusCode int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(value); err != nil && !errors.Is(err, http.ErrHandlerTimeout) {
		return
	}
}
