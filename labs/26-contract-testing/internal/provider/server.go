package provider

import (
	"encoding/json"
	"net/http"
	"strings"

	"labs/26-contract-testing/internal/model"
)

// ProviderV1 implements original API server fulfilling V1 contract.
type ProviderV1 struct{}

func NewProviderV1() *ProviderV1 {
	return &ProviderV1{}
}

func (p *ProviderV1) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if strings.HasPrefix(r.URL.Path, "/v1/orders/") {
		id := strings.TrimPrefix(r.URL.Path, "/v1/orders/")
		resp := model.OrderResponseV1{
			ID:     id,
			Status: "IN_PROGRESS",
			Customer: model.CustomerResponseV1{
				ID:   "cust-100",
				Name: "Budi Santoso",
			},
			Total: 150000,
			Notes: "Internal note not used by mobile app",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	http.NotFound(w, r)
}

// ProviderBreaking implements API server with breaking schema changes.
type ProviderBreaking struct{}

func NewProviderBreaking() *ProviderBreaking {
	return &ProviderBreaking{}
}

func (p *ProviderBreaking) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if strings.HasPrefix(r.URL.Path, "/v1/orders/") {
		id := strings.TrimPrefix(r.URL.Path, "/v1/orders/")
		// Introduces 3 breaking changes: status casing, field rename customer.full_name, total string
		resp := model.OrderResponseBreaking{
			ID:     id,
			Status: "in_progress",
			Customer: model.CustomerResponseBreaking{
				ID:       "cust-100",
				FullName: "Budi Santoso",
			},
			Total: "150000",
			Notes: "Breaking change version",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	http.NotFound(w, r)
}

// ProviderDual supports V1 (contract compliant) alongside V2 (new feature schema).
type ProviderDual struct{}

func NewProviderDual() *ProviderDual {
	return &ProviderDual{}
}

func (p *ProviderDual) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if strings.HasPrefix(r.URL.Path, "/v1/orders/") {
		id := strings.TrimPrefix(r.URL.Path, "/v1/orders/")
		respV1 := model.OrderResponseV1{
			ID:     id,
			Status: "IN_PROGRESS",
			Customer: model.CustomerResponseV1{
				ID:   "cust-100",
				Name: "Budi Santoso",
			},
			Total: 150000,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(respV1)
		return
	}

	if strings.HasPrefix(r.URL.Path, "/v2/orders/") {
		id := strings.TrimPrefix(r.URL.Path, "/v2/orders/")
		respV2 := model.OrderResponseV2{
			ID:     id,
			Status: "in_progress",
			Customer: model.CustomerResponseV2{
				ID:       "cust-100",
				FullName: "Budi Santoso",
			},
			Total:    "150000",
			Currency: "IDR",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(respV2)
		return
	}

	http.NotFound(w, r)
}
