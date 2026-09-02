package api_versioning

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// AdditiveHandler mensimulasikan penambahan field baru pada backend
// (misal: "currency") tanpa membuat contract major version baru.
// Endpoint: GET /api/invoices/:id
func AdditiveHandler(w http.ResponseWriter, r *http.Request) {
	prefix := "/api/invoices/"
	if len(r.URL.Path) <= len(prefix) || r.URL.Path[:len(prefix)] != prefix {
		http.Error(w, `{"error": "invalid path"}`, http.StatusBadRequest)
		return
	}
	idStr := r.URL.Path[len(prefix):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error": "invalid id"}`, http.StatusBadRequest)
		return
	}

	// Backend menambah field `currency` secara sepihak
	invoice := map[string]interface{}{
		"id":       id,
		"customer": "Budi", // Tipe tetap string (tidak merusak legacy)
		"total":    500000,
		"status":   "PAID",
		"currency": "IDR", // ADDITIVE CHANGE: Field baru
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(invoice)
}
