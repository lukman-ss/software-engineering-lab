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
	// Method validation
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	prefix := "/api/invoices/"
	if len(r.URL.Path) <= len(prefix) || r.URL.Path[:len(prefix)] != prefix {
		writeJSONError(w, http.StatusBadRequest, "invalid path")
		return
	}
	idStr := r.URL.Path[len(prefix):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid id")
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

	writeJSON(w, http.StatusOK, invoice)
}