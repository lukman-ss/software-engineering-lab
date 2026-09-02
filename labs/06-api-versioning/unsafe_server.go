package api_versioning

import (
	"net/http"
)

// UnsafeHandler merepresentasikan anti-pattern ketika backend mengubah kontrak API
// secara merusak tanpa versioning. Response V2 dengan `customer` sebagai object
// akan membuat legacy client (yang mengharapkan string) gagal decode.
//
// Mental model: API adalah kontrak. Backend berhasil compile ≠ backward compatible.
// HTTP 200 ≠ backward compatible.
//
// Endpoint: GET /api/invoices/:id
// Hanya menerima GET method. Method lain mengembalikan 405.
func UnsafeHandler(w http.ResponseWriter, r *http.Request) {
	// HTTP method validation - hanya GET yang diizinkan
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Gunakan helper parseInvoiceID untuk konsistensi dengan V1/V2
	id, isMissing, isInvalid, errMsg := parseInvoiceID(r.URL.Path, "/api/invoices/")

	if isMissing {
		writeJSONError(w, http.StatusBadRequest, errMsg)
		return
	}
	if isInvalid {
		writeJSONError(w, http.StatusBadRequest, errMsg)
		return
	}

	// Simulasi: backend pernah punya customer="Budi", sekarang jadi object
	// Ini adalah contoh BREAKING CHANGE yang tidak terdeteksi unit test
	invoice := map[string]interface{}{
		"id": id,
		"customer": map[string]interface{}{ // BREAKING: sebelumnya string "Budi"
			"id":    15,
			"name":  "Budi",
			"phone": "08123",
		},
		"total":  500000,
		"status": "PAID",
	}

	writeJSON(w, http.StatusOK, invoice)
}
