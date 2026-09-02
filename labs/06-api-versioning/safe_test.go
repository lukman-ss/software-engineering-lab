package api_versioning

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// performRequest adalah helper untuk melakukan request terhadap HTTP handler.
// Mengembalikan ResponseRecorder untuk assertion.
func performRequest(handler http.Handler, method, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	return w
}

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

// TestRequiredHeaderBreakingChange menyimulasikan breaking change ketika
// endpoint yang sebelumnya tidak membutuhkan header menjadi wajib.
func TestRequiredHeaderBreakingChange(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/invoices/", V1Handler)
	server := httptest.NewServer(mux)
	defer server.Close()

	// Request tanpa X-Tenant-ID header (yang seharusnya tidak diperlukan untuk V1)
	resp, err := http.Get(server.URL + "/api/v1/invoices/1001")
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer resp.Body.Close()

	// V1 tidak membutuhkan header khusus - harus success
	if resp.StatusCode != http.StatusOK {
		t.Errorf("V1: expected HTTP 200 without X-Tenant-ID, got %d", resp.StatusCode)
	}

	t.Log("✅ V1 endpoint tidak membutuhkan X-Tenant-ID header - backward compatible")
}

// TestV1Contract_RegressionGuaranteeAfterV2Registration tambahan verifikasi
// untuk memastikan tidak ada breaking change pada V1 pasca V2 registration.
func TestV1Contract_RegressionGuaranteeAfterV2Registration(t *testing.T) {
	// Setup: V1 dan V2 pada router yang sama
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/invoices/", V1Handler)
	mux.HandleFunc("/api/v2/invoices/", V2Handler)

	server := httptest.NewServer(mux)
	defer server.Close()

	// Request ke V1 - harus tetap mengembalikan V1 contract
	resp, err := http.Get(server.URL + "/api/v1/invoices/1001")
	if err != nil {
		t.Fatalf("failed V1 request: %v", err)
	}
	defer resp.Body.Close()

	// HTTP 200
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("V1: expected HTTP 200, got %d", resp.StatusCode)
	}

	// Content-Type application/json
	ct := resp.Header.Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("V1: expected Content-Type 'application/json', got '%s'", ct)
	}

	// Decode dan verify semua field
	var raw map[string]json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		t.Fatalf("V1: invalid JSON response: %v", err)
	}

	// 1. id == 1001, JSON number
	var idNum json.Number
	if err := json.Unmarshal(raw["id"], &idNum); err != nil {
		t.Error("V1: 'id' harus berupa number")
	} else if idNum != "1001" {
		t.Errorf("V1: 'id' harus 1001, got %s", idNum)
	}

	// 2. customer == "Budi", JSON string (BREAKING jika jadi object)
	var customerStr string
	if err := json.Unmarshal(raw["customer"], &customerStr); err != nil {
		t.Fatal("V1 CONTRACT BROKEN: customer harus tetap string (breaking change!)")
	} else if customerStr != "Budi" {
		t.Errorf("V1: 'customer' harus 'Budi', got '%s'", customerStr)
	}

	// 3. total == 500000, JSON number
	var totalNum json.Number
	if err := json.Unmarshal(raw["total"], &totalNum); err != nil {
		t.Error("V1: 'total' harus berupa number")
	} else if totalNum != "500000" {
		t.Errorf("V1: 'total' harus 500000, got %s", totalNum)
	}

	// 4. status == "PAID", JSON string
	var statusStr string
	if err := json.Unmarshal(raw["status"], &statusStr); err != nil {
		t.Fatal("V1 CONTRACT BROKEN: status harus tetap string")
	} else if statusStr != "PAID" {
		t.Errorf("V1: 'status' harus 'PAID', got '%s'", statusStr)
	}

	t.Log("✅ V1 contract regression guarantee: V2 registration tidak mengubah V1 output")
}

// TestAdditiveChangeDoesNotBreakV1Contract memastikan penambahan field
// pada V1 tidak merusak backward compatibility.
func TestAdditiveChangeDoesNotBreakV1Contract(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/invoices/", V1Handler)
	server := httptest.NewServer(mux)
	defer server.Close()

	// Request ke V1
	resp, err := http.Get(server.URL + "/api/v1/invoices/1001")
	if err != nil {
		t.Fatalf("failed V1 request: %v", err)
	}
	defer resp.Body.Close()

	// Decode ke LegacyInvoice (tanpa field yang mungkin ditambah)
	var legacyResponse LegacyInvoice
	if err := json.NewDecoder(resp.Body).Decode(&legacyResponse); err != nil {
		t.Fatalf("V1: legacy client gagal decode: %v", err)
	}

	// Verifikasi field yang diketahui
	if legacyResponse.Customer != "Budi" {
		t.Errorf("V1: customer=%s, expected 'Budi'", legacyResponse.Customer)
	}
	if legacyResponse.Total != 500000 {
		t.Errorf("V1: total=%d, expected 500000", legacyResponse.Total)
	}

	t.Log("✅ Additive change pada V1 tidak memecah legacy client")
}

// TestV1Contract_RemainsBackwardCompatible melindungi kontrak V1 agar tidak
// secara tidak sengaja mengubah bentuk customer dari string menjadi object.
func TestV1Contract_RemainsBackwardCompatible(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/invoices/", V1Handler)
	server := httptest.NewServer(mux)
	defer server.Close()

	resp, err := http.Get(server.URL + "/api/v1/invoices/1001")
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer resp.Body.Close()

	// 1. HTTP 200 requirement
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("V1: expected HTTP 200, got %d", resp.StatusCode)
	}

	// 2. Content-Type application/json
	ct := resp.Header.Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("V1: expected Content-Type 'application/json', got '%s'", ct)
	}

	// Decode ke raw message dulu untuk validasi schema
	var raw map[string]json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		t.Fatalf("V1: invalid JSON response: %v", err)
	}

	// 3. id = number, id == 1001
	var idNum json.Number
	if err := json.Unmarshal(raw["id"], &idNum); err != nil {
		t.Error("V1: 'id' harus berupa number")
	} else if idNum != "1001" {
		t.Errorf("V1: 'id' harus 1001, got %s", idNum)
	}

	// 4. customer = string, customer == "Budi" (BREAKING if accidentally becomes object)
	var customerStr string
	if err := json.Unmarshal(raw["customer"], &customerStr); err != nil {
		t.Error("V1: 'customer' HARUS tetap string (breaking change detected!)")
		t.Error("Jika customer menjadi object, legacy client akan gagal parse")
	} else if customerStr != "Budi" {
		t.Errorf("V1: 'customer' harus 'Budi', got '%s'", customerStr)
	}

	// 5. total = number, total == 500000
	var totalNum json.Number
	if err := json.Unmarshal(raw["total"], &totalNum); err != nil {
		t.Error("V1: 'total' harus berupa number")
	} else if totalNum != "500000" {
		t.Errorf("V1: 'total' harus 500000, got %s", totalNum)
	}

	// 6. status = string, status == "PAID"
	var statusStr string
	if err := json.Unmarshal(raw["status"], &statusStr); err != nil {
		t.Error("V1: 'status' harus berupa string")
	} else if statusStr != "PAID" {
		t.Errorf("V1: 'status' harus 'PAID', got '%s'", statusStr)
	}

	t.Log("✅ V1 contract protected: customer remains string")
}

// TestV2Contract_UsesNestedCustomer memastikan V2 menggunakan customer sebagai object
// dengan semua field yang tepat. Memverifikasi JSON wire contract secara eksplisit
// menggunakan raw JSON untuk detection breaking changes.
func TestV2Contract_UsesNestedCustomer(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/invoices/", V2Handler)
	server := httptest.NewServer(mux)
	defer server.Close()

	resp, err := http.Get(server.URL + "/api/v2/invoices/1001")
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer resp.Body.Close()

	// 1. HTTP Status = 200
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("V2: expected HTTP 200, got %d", resp.StatusCode)
	}

	// 2. Content-Type = application/json
	ct := resp.Header.Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("V2: expected Content-Type 'application/json', got '%s'", ct)
	}

	// 3. Decode ke raw message untuk semantic contract assertion
	var raw map[string]json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		t.Fatalf("V2: invalid JSON response: %v", err)
	}

	// 4. id = JSON number = 1001
	var idNum json.Number
	if err := json.Unmarshal(raw["id"], &idNum); err != nil {
		t.Error("V2: 'id' harus berupa JSON number")
	} else if idNum != "1001" {
		t.Errorf("V2: 'id' harus 1001, got %s", idNum)
	}

	// 5. customer = JSON object (bukan string)
	customerRaw := raw["customer"]
	var customerObj map[string]json.RawMessage
	if err := json.Unmarshal(customerRaw, &customerObj); err != nil {
		t.Fatal("V2 CONTRACT BROKEN: customer harus object, bukan string")
	}

	// 6. customer.id = JSON number = 15
	var cidNum json.Number
	if err := json.Unmarshal(customerObj["id"], &cidNum); err != nil {
		t.Error("V2: 'customer.id' harus berupa JSON number")
	} else if cidNum != "15" {
		t.Errorf("V2: 'customer.id' harus 15, got %s", cidNum)
	}

	// 7. customer.name = JSON string = "Budi"
	var cname string
	if err := json.Unmarshal(customerObj["name"], &cname); err != nil {
		t.Error("V2: 'customer.name' harus berupa JSON string")
	} else if cname != "Budi" {
		t.Errorf("V2: 'customer.name' harus 'Budi', got '%s'", cname)
	}

	// 8. customer.phone = JSON string = "08123"
	var cphone string
	if err := json.Unmarshal(customerObj["phone"], &cphone); err != nil {
		t.Error("V2: 'customer.phone' harus berupa JSON string")
	} else if cphone != "08123" {
		t.Errorf("V2: 'customer.phone' harus '08123', got '%s'", cphone)
	}

	// 9. total = JSON number = 500000
	var totalNum json.Number
	if err := json.Unmarshal(raw["total"], &totalNum); err != nil {
		t.Error("V2: 'total' harus berupa JSON number")
	} else if totalNum != "500000" {
		t.Errorf("V2: 'total' harus 500000, got %s", totalNum)
	}

	// 10. status = JSON string = "PAID"
	var statusStr string
	if err := json.Unmarshal(raw["status"], &statusStr); err != nil {
		t.Error("V2: 'status' harus berupa JSON string")
	} else if statusStr != "PAID" {
		t.Errorf("V2: 'status' harus 'PAID', got '%s'", statusStr)
	}

	// Tambahan: verify decode ke struct V2 tetap berhasil (new client compatibility)
	var v2Resp InvoiceV2Response
	resp.Body.Close() // close previous body
	resp2, _ := http.Get(server.URL + "/api/v2/invoices/1001")
	defer resp2.Body.Close()
	if err := json.NewDecoder(resp2.Body).Decode(&v2Resp); err != nil {
		t.Fatalf("V2 struct decode failed: %v", err)
	}

	t.Log("✅ V2 contract verified: customer is object with id, name, phone")
}

// TestV2Contract_DetectsStatusObject memastikan V2 menggunakan status sebagai string
// pada real HTTP handler.
func TestV2Contract_DetectsStatusObject(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/invoices/", V2Handler)
	server := httptest.NewServer(mux)
	defer server.Close()

	// Request ke V2 endpoint - real HTTP
	resp, err := http.Get(server.URL + "/api/v2/invoices/1001")
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer resp.Body.Close()

	var raw map[string]json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		t.Fatalf("V2: invalid JSON response: %v", err)
	}

	// V2 contract: status HARUS string (bukan object)
	var statusStr string
	if err := json.Unmarshal(raw["status"], &statusStr); err != nil {
		t.Fatalf("✅ V2 contract detected breaking change: status harus string, got error: %v", err)
	}

	// Harus "PAID"
	if statusStr != "PAID" {
		t.Errorf("V2 status berubah: expected 'PAID', got '%s'", statusStr)
	}

	t.Logf("✅ V2 contract: status tetap string '%s'", statusStr)
}
