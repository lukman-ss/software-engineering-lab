package api_versioning

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// InvoiceV1Response adalah contract API untuk V1.
// Customer direpresentasikan sebagai string untuk backward compatibility
// dengan client yang sudah ada.
type InvoiceV1Response struct {
	ID       int    `json:"id"`
	Customer string `json:"customer"` // String: nama customer
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
	Customer CustomerV2Response `json:"customer"` // Object
	Total    int64              `json:"total"`
	Status   string             `json:"status"`
}

// mapToV1 mengonversi domain Invoice ke V1 response contract.
// Ini memisahkan transport layer (API contract) dari domain model.
func mapToV1(domain Invoice) InvoiceV1Response {
	return InvoiceV1Response{
		ID:       domain.ID,
		Customer: domain.Customer.Name, // Hanya pakai nama untuk backward compat
		Total:    int64(domain.Total),
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

// mockInvoiceRepository simulasi repository untuk testing.
type mockInvoiceRepository struct{}

func (r *mockInvoiceRepository) GetInvoiceByID(ctx any, id int) (Invoice, error) {
	return Invoice{
		ID: id,
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

// V1Handler mengembalikan response dalam format V1 (customer = string).
// Endpoint: GET /api/v1/invoices/:id
func V1Handler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/v1/invoices/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error": "invalid id"}`, http.StatusBadRequest)
		return
	}

	domain, err := repo.GetInvoiceByID(r.Context(), id)
	if err != nil {
		http.Error(w, `{"error": "not found"}`, http.StatusNotFound)
		return
	}

	resp := mapToV1(domain)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// V2Handler mengembalikan response dalam format V2 (customer = object).
// Endpoint: GET /api/v2/invoices/:id
func V2Handler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/v2/invoices/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error": "invalid id"}`, http.StatusBadRequest)
		return
	}

	domain, err := repo.GetInvoiceByID(r.Context(), id)
	if err != nil {
		http.Error(w, `{"error": "not found"}`, http.StatusNotFound)
		return
	}

	resp := mapToV2(domain)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
