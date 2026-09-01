package api_versioning

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestSafeVersioning_LegacyClientSucceedsOnV1 membuktikan bahwa dengan API versioning,
// kita dapat meluncurkan perubahan tanpa merusak client yang sudah ada.
func TestSafeVersioning_LegacyClientSucceedsOnV1(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/invoices/", V1Handler)

	server := httptest.NewServer(mux)
	defer server.Close()

	// Legacy client menggunakan V1 API
	resp, err := http.Get(server.URL + "/api/v1/invoices/1001")
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected HTTP 200, got %d", resp.StatusCode)
	}

	// Legacy client berhasil decode ke struct V1
	var legacyResponse InvoiceV1Response
	err = json.NewDecoder(resp.Body).Decode(&legacyResponse)
	if err != nil {
		t.Fatalf("legacy client SHOULD succeed on V1: %v", err)
	}

	// Verify response structure
	if legacyResponse.Customer != "Budi" {
		t.Errorf("expected Customer='Budi', got '%s'", legacyResponse.Customer)
	}
	if legacyResponse.Total != 500000 {
		t.Errorf("expected Total=500000, got %d", legacyResponse.Total)
	}

	t.Log("✅ Legacy client successfully reads from V1 API")
}

// TestSafeVersioning_NewClientSucceedsOnV2 membuktikan client baru dapat
// menggunakan V2 API dengan customer sebagai object.
func TestSafeVersioning_NewClientSucceedsOnV2(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/invoices/", V2Handler)

	server := httptest.NewServer(mux)
	defer server.Close()

	// New client menggunakan V2 API
	resp, err := http.Get(server.URL + "/api/v2/invoices/1001")
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected HTTP 200, got %d", resp.StatusCode)
	}

	// New client decode ke struct V2
	var v2Response InvoiceV2Response
	err = json.NewDecoder(resp.Body).Decode(&v2Response)
	if err != nil {
		t.Fatalf("new client SHOULD succeed on V2: %v", err)
	}

	// Verify response structure
	if v2Response.Customer.Name != "Budi" {
		t.Errorf("expected Customer.Name='Budi', got '%s'", v2Response.Customer.Name)
	}
	if v2Response.Customer.Phone != "08123" {
		t.Errorf("expected Customer.Phone='08123', got '%s'", v2Response.Customer.Phone)
	}
	if v2Response.Total != 500000 {
		t.Errorf("expected Total=500000, got %d", v2Response.Total)
	}

	t.Log("✅ New client successfully reads from V2 API")
	t.Log("✅ Both V1 and V2 contracts live side-by-side safely")
}

// TestVersionedRoutes_CanRunSideBySide membuktikan bahwa V1 dan V2 dapat
// berjalan bersamaan di server yang sama dengan routing yang tepat.
func TestVersionedRoutes_CanRunSideBySide(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/invoices/", V1Handler)
	mux.HandleFunc("/api/v2/invoices/", V2Handler)

	server := httptest.NewServer(mux)
	defer server.Close()

	// Request ke V1
	resp1, err := http.Get(server.URL + "/api/v1/invoices/1001")
	if err != nil {
		t.Fatalf("failed V1 request: %v", err)
	}
	defer resp1.Body.Close()

	if resp1.StatusCode != http.StatusOK {
		t.Fatalf("V1: expected HTTP 200, got %d", resp1.StatusCode)
	}

	var v1Resp InvoiceV1Response
	if err := json.NewDecoder(resp1.Body).Decode(&v1Resp); err != nil {
		t.Fatalf("V1 decode failed: %v", err)
	}

	// Request ke V2
	resp2, err := http.Get(server.URL + "/api/v2/invoices/1001")
	if err != nil {
		t.Fatalf("failed V2 request: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("V2: expected HTTP 200, got %d", resp2.StatusCode)
	}

	var v2Resp InvoiceV2Response
	if err := json.NewDecoder(resp2.Body).Decode(&v2Resp); err != nil {
		t.Fatalf("V2 decode failed: %v", err)
	}

	// Verifikasi kontrak berbeda
	if v1Resp.Customer != "Budi" {
		t.Errorf("V1: expected Customer='Budi', got '%s'", v1Resp.Customer)
	}
	if v2Resp.Customer.Name != "Budi" {
		t.Errorf("V2: expected Customer.Name='Budi', got '%s'", v2Resp.Customer.Name)
	}

	// Verifikasi tipe customer berbeda
	t.Logf("V1 contract: customer=%v (type: %T)", v1Resp.Customer, v1Resp.Customer)
	t.Logf("V2 contract: customer=%v (type: %T)", v2Resp.Customer, v2Resp.Customer)
	t.Log("✅ Both versions co-exist safely with different contracts")
}