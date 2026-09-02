package api_versioning

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

// InvoiceV1Response adalah contract API untuk V1.
// Customer direpresentasikan sebagai string untuk backward compatibility.
type InvoiceV1Response struct {
	ID       int    `json:"id"`
	Customer string `json:"customer"`
	Total    int64  `json:"total"`
	Status   string `json:"status"`
}

// CustomerV2Response adalah representasi customer untuk API V2.
type CustomerV2Response struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

// InvoiceV2Response adalah contract API untuk V2.
// Customer sekarang object dengan informasi lebih lengkap.
type InvoiceV2Response struct {
	ID       int                `json:"id"`
	Customer CustomerV2Response `json:"customer"`
	Total    int64              `json:"total"`
	Status   string             `json:"status"`
}

// mapToV1 mengonversi domain Invoice ke V1 response contract.
// WHY: memisahkan transport layer (API contract) dari domain model.
func mapToV1(domain Invoice) InvoiceV1Response {
	return InvoiceV1Response{
		ID:       domain.ID,
		Customer: domain.Customer.Name,
		Total:    domain.Total,
		Status:   domain.Status,
	}
}

// mapToV2 mengonversi domain Invoice ke V2 response contract.
func mapToV2(domain Invoice) InvoiceV2Response {
	return InvoiceV2Response{
		ID: domain.ID,
		Customer: CustomerV2Response{
			ID:    domain.Customer.ID,
			Name:  domain.Customer.Name,
			Phone: domain.Customer.Phone,
		},
		Total:  domain.Total,
		Status: domain.Status,
	}
}

// ErrInvoiceNotFound adalah error sentinel untuk invoice yang tidak ditemukan.
// WHY: memungkinkan handler menerjemahkan ke HTTP 404.
var ErrInvoiceNotFound = errors.New("invoice not found")

// InvoiceRepository mendefinisikan interface untuk mengakses data invoice.
// Menggunakan context.Context untuk cancellation dan timeout.
type InvoiceRepository interface {
	GetInvoiceByID(ctx context.Context, id int) (Invoice, error)
}

// mockInvoiceRepository simulasi repository untuk testing.
// WHY: mensimulasikan database access tanpa dependency eksternal.
type mockInvoiceRepository struct{}

func (r *mockInvoiceRepository) GetInvoiceByID(ctx context.Context, id int) (Invoice, error) {
	// Simulasi: hanya ID 1001 yang ada di database
	if id != 1001 {
		return Invoice{}, ErrInvoiceNotFound
	}
	return Invoice{
		ID: 1001,
		Customer: Customer{
			ID:    15,
			Name:  "Budi",
			Phone: "08123",
		},
		Total:  500000,
		Status: "PAID",
	}, nil
}

// repo instance untuk penggunaan handler
var repo InvoiceRepository = &mockInvoiceRepository{}

// writeJSON menulis JSON response dengan Content-Type yang benar.
// WHY: memastikan semua response API menggunakan application/json, bukan text/plain.
func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// writeJSONError menulis error response dalam format JSON konsisten.
func writeJSONError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// parseInvoiceID mengekstrak dan memvalidasi ID dari path.
// Menangani: missing ID, non-numeric ID, extra path segment.
// Returns: (id, isMissing, isInvalid, errorMessage)
// WHY: menghindari slicing path mentah yang rapuh dan menyediakan validasi terpusat.
func parseInvoiceID(path, prefix string) (int, bool, bool, string) {
	if !strings.HasPrefix(path, prefix) {
		return 0, false, false, "invalid path"
	}

	idStr := strings.TrimPrefix(path, prefix)
	if idStr == "" {
		return 0, true, false, "missing invoice ID"
	}

	// Cek extra path segment (e.g., /api/v1/invoices/1001/extra)
	if strings.Contains(idStr, "/") {
		return 0, false, true, "invalid ID format"
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, false, true, "invalid ID: must be numeric"
	}

	return id, false, false, ""
}

// V1Handler mengembalikan response dalam format V1 (customer = string).
// Endpoint: GET /api/v1/invoices/:id
func V1Handler(w http.ResponseWriter, r *http.Request) {
	// Method validation
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	id, isMissing, isInvalid, errMsg := parseInvoiceID(r.URL.Path, "/api/v1/invoices/")
	if isMissing {
		writeJSONError(w, http.StatusBadRequest, errMsg)
		return
	}
	if isInvalid {
		writeJSONError(w, http.StatusBadRequest, errMsg)
		return
	}

	domain, err := repo.GetInvoiceByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrInvoiceNotFound) {
			writeJSONError(w, http.StatusNotFound, "invoice not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, mapToV1(domain))
}

// V2Handler mengembalikan response dalam format V2 (customer = object).
// Endpoint: GET /api/v2/invoices/:id
func V2Handler(w http.ResponseWriter, r *http.Request) {
	// Method validation
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	id, isMissing, isInvalid, errMsg := parseInvoiceID(r.URL.Path, "/api/v2/invoices/")
	if isMissing {
		writeJSONError(w, http.StatusBadRequest, errMsg)
		return
	}
	if isInvalid {
		writeJSONError(w, http.StatusBadRequest, errMsg)
		return
	}

	domain, err := repo.GetInvoiceByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrInvoiceNotFound) {
			writeJSONError(w, http.StatusNotFound, "invoice not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, mapToV2(domain))
}
