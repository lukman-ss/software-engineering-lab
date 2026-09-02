package api_versioning

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// UnsafeHandler merepresentasikan anti-pattern ketika backend mengubah kontrak API
// secara merusak tanpa versioning. Response V2 dengan `customer` sebagai object
// akan membuat legacy client (yang mengharapkan string) gagal decode.
//
// Mental model: API adalah kontrak. Backend berhasil compile ≠ backward compatible.
// HTTP 200 ≠ backward compatible.
func UnsafeHandler(w http.ResponseWriter, r *http.Request) {
	// Endpoint: GET /api/invoices/:id
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(invoice)
}
