package api_versioning

import (
	"net/http"
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

	// Gunakan helper parseInvoiceID untuk konsistensi
	id, isMissing, isInvalid, errMsg := parseInvoiceID(r.URL.Path, "/api/invoices/")

	if isMissing {
		writeJSONError(w, http.StatusBadRequest, errMsg)
		return
	}
	if isInvalid {
		writeJSONError(w, http.StatusBadRequest, errMsg)
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
