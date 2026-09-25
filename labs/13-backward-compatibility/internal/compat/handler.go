package compat

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type APIHandler struct {
	svc *Service
}

func NewAPIHandler(svc *Service) *APIHandler {
	return &APIHandler{svc: svc}
}

// GetUserV1 is the legacy endpoint
// Demonstrates Deprecation headers when legacy field/endpoint is hit
func (h *APIHandler) GetUserV1(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	dto, err := h.svc.GetLegacyUser(id)
	if err != nil {
		if errors.Is(err, ErrLegacyUnavailable) {
			w.WriteHeader(http.StatusGone) // 410 Gone after sunset
			w.Write([]byte(`{"error": "legacy v1 endpoint has been permanently removed"}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Inject Deprecation Headers
	w.Header().Set("Deprecation", "true")
	w.Header().Set("Sunset", "Mon, 31 Dec 2026 23:59:59 GMT") // Sunset date
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(dto)
}

// GetUserV2 is the modern endpoint
func (h *APIHandler) GetUserV2(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	dto, err := h.svc.GetModernUser(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto)
}
