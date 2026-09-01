package api_versioning

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// AdditiveHandler mensimulasikan penambahan field baru pada backend
// (misal: "currency") tanpa membuat contract major version baru.
func AdditiveHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
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