package api_versioning

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// performRequest adalah helper untuk melakukan request terhadap HTTP handler.
// Mengembalikan ResponseRecorder untuk assertion.

// assertV1WireContract adalah helper untuk verifikasi kontrak wire V1.
// Menyimpan satu sumber kebenaran untuk V1 contract assertions.
func assertV1WireContract(t *testing.T, body []byte) {
	t.Helper()

	// Decode ke raw message untuk semantic contract assertion
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		t.Fatalf("V1: invalid JSON response: %v", err)
	}

	// Field existence checks
	assertV1FieldExists(t, raw, "id")
	assertV1FieldExists(t, raw, "customer")
	assertV1FieldExists(t, raw, "total")
	assertV1FieldExists(t, raw, "status")

	// 1. id = JSON number = 1001
	var idNum json.Number
	if err := json.Unmarshal(raw["id"], &idNum); err != nil {
		t.Error("V1: 'id' harus berupa JSON number")
	} else if idNum != "1001" {
		t.Errorf("V1: 'id' harus 1001, got %s", idNum)
	}

	// 2. customer = JSON string = "Budi" (BREAKING if accidentally becomes object)
	var customerStr string
	if err := json.Unmarshal(raw["customer"], &customerStr); err != nil {
		t.Error("V1: 'customer' HARUS tetap string (breaking change detected!)")
		t.Error("Jika customer menjadi object, legacy client akan gagal parse")
	} else if customerStr != "Budi" {
		t.Errorf("V1: 'customer' harus 'Budi', got '%s'", customerStr)
	}

	// 3. total = JSON number = 500000
	var totalNum json.Number
	if err := json.Unmarshal(raw["total"], &totalNum); err != nil {
		t.Error("V1: 'total' harus berupa JSON number")
	} else if totalNum != "500000" {
		t.Errorf("V1: 'total' harus 500000, got %s", totalNum)
	}

	// 4. status = JSON string = "PAID"
	var statusStr string
	if err := json.Unmarshal(raw["status"], &statusStr); err != nil {
		t.Error("V1: 'status' harus berupa JSON string")
	} else if statusStr != "PAID" {
		t.Errorf("V1: 'status' harus 'PAID', got '%s'", statusStr)
	}
}

// assertV1FieldExists memastikan field ada di response V1
func assertV1FieldExists(t *testing.T, raw map[string]json.RawMessage, field string) {
	t.Helper()
	if _, ok := raw[field]; !ok {
		t.Fatalf("V1 contract broken: missing field '%s'", field)
	}
}

// assertV2WireContract adalah helper untuk verifikasi kontrak wire V2.
// Menyimpan satu sumber kebenaran untuk V2 contract assertions.
func assertV2WireContract(t *testing.T, body []byte) {
	t.Helper()

	// Decode ke raw message untuk semantic contract assertion
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		t.Fatalf("V2: invalid JSON response: %v", err)
	}

	// Field existence checks
	assertV2FieldExists(t, raw, "id")
	assertV2FieldExists(t, raw, "customer")
	assertV2FieldExists(t, raw, "total")
	assertV2FieldExists(t, raw, "status")

	// Check customer object exists
	customerRaw, ok := raw["customer"]
	if !ok {
		t.Fatalf("V2 contract broken: missing customer")
	}
	var customerObj map[string]json.RawMessage
	if err := json.Unmarshal(customerRaw, &customerObj); err != nil {
		t.Fatal("V2: 'customer' harus berupa JSON object")
	}

	// Nested customer field existence checks
	assertV2NestedFieldExists(t, customerObj, "id")
	assertV2NestedFieldExists(t, customerObj, "name")
	assertV2NestedFieldExists(t, customerObj, "phone")

	// 1. id = JSON number = 1001
	var idNum json.Number
	if err := json.Unmarshal(raw["id"], &idNum); err != nil {
		t.Error("V2: 'id' harus berupa JSON number")
	} else if idNum != "1001" {
		t.Errorf("V2: 'id' harus 1001, got %s", idNum)
	}

	// 2. customer = JSON object
	// 3. customer.id = JSON number = 15
	var cidNum json.Number
	if err := json.Unmarshal(customerObj["id"], &cidNum); err != nil {
		t.Error("V2: 'customer.id' harus berupa JSON number")
	} else if cidNum != "15" {
		t.Errorf("V2: 'customer.id' harus 15, got %s", cidNum)
	}

	// 4. customer.name = JSON string = "Budi"
	var cname string
	if err := json.Unmarshal(customerObj["name"], &cname); err != nil {
		t.Error("V2: 'customer.name' harus berupa JSON string")
	} else if cname != "Budi" {
		t.Errorf("V2: 'customer.name' harus 'Budi', got '%s'", cname)
	}

	// 5. customer.phone = JSON string = "08123"
	var cphone string
	if err := json.Unmarshal(customerObj["phone"], &cphone); err != nil {
		t.Error("V2: 'customer.phone' harus berupa JSON string")
	} else if cphone != "08123" {
		t.Errorf("V2: 'customer.phone' harus '08123', got '%s'", cphone)
	}

	// 6. total = JSON number = 500000
	var totalNum json.Number
	if err := json.Unmarshal(raw["total"], &totalNum); err != nil {
		t.Error("V2: 'total' harus berupa JSON number")
	} else if totalNum != "500000" {
		t.Errorf("V2: 'total' harus 500000, got %s", totalNum)
	}

	// 7. status = JSON string = "PAID"
	var statusStr string
	if err := json.Unmarshal(raw["status"], &statusStr); err != nil {
		t.Error("V2: 'status' harus berupa JSON string")
	} else if statusStr != "PAID" {
		t.Errorf("V2: 'status' harus 'PAID', got '%s'", statusStr)
	}
}

// assertV2FieldExists memastikan field ada di response V2
func assertV2FieldExists(t *testing.T, raw map[string]json.RawMessage, field string) {
	t.Helper()
	if _, ok := raw[field]; !ok {
		t.Fatalf("V2 contract broken: missing field '%s'", field)
	}
}

// assertV2NestedFieldExists memastikan nested field ada di object customer
func assertV2NestedFieldExists(t *testing.T, obj map[string]json.RawMessage, field string) {
	t.Helper()
	if _, ok := obj[field]; !ok {
		t.Fatalf("V2 contract broken: missing customer.%s", field)
	}
}
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

// TestV1Contract_DoesNotIntroduceRequiredTenantHeader membuktikan bahwa
// V1 tidak secara tidak sengaja memperkenalkan X-Tenant-ID sebagai required header.
func TestV1Contract_DoesNotIntroduceRequiredTenantHeader(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/invoices/", V1Handler)
	server := httptest.NewServer(mux)
	defer server.Close()

	// Request tanpa X-Tenant-ID header - legacy client tidak mengirim header ini
	resp, err := http.Get(server.URL + "/api/v1/invoices/1001")
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer resp.Body.Close()

	// V1 tidak membutuhkan header khusus - harus success (backward compatible)
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

	// Baca body sekali untuk semua asserts
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("V1: failed to read body: %v", err)
	}

	// Use helper untuk semua V1 contract assertions including field existence
	assertV1WireContract(t, body)

	t.Log("✅ V1 contract regression guarantee: V2 registration tidak mengubah V1 output")
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

	// HTTP 200 requirement
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("V1: expected HTTP 200, got %d", resp.StatusCode)
	}

	// Content-Type application/json
	ct := resp.Header.Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("V1: expected Content-Type 'application/json', got '%s'", ct)
	}

	// Baca body sekali untuk semua asserts
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("V1: failed to read body: %v", err)
	}

	// Use helper untuk semua V1 contract assertions
	assertV1WireContract(t, body)

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
		t.Fatalf("V2: failed to connect: %v", err)
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

	// Baca body sekali untuk semua assertions
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("V2: failed to read body: %v", err)
	}

	// Use helper untuk semua V2 contract assertions including field existence
	assertV2WireContract(t, body)

	// Additional: verify decode ke struct V2 tetap berhasil (new client compatibility)
	var v2Resp InvoiceV2Response
	if err := json.Unmarshal(body, &v2Resp); err != nil {
		t.Fatalf("V2 struct decode failed: %v", err)
	}

	t.Log("✅ V2 contract verified: customer is object with id, name, phone")
}

// TestV1Handler_WrongPrefixReturnsBadRequest memastikan V1Handler
// tidak menerima request dengan prefix yang salah.
func TestV1Handler_WrongPrefixReturnsBadRequest(t *testing.T) {
	handler := http.HandlerFunc(V1Handler)

	// GET /api/v2/invoices/1001 ketika dipanggil langsung ke V1Handler
	w := performRequest(handler, http.MethodGet, "/api/v2/invoices/1001")
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected HTTP 400 for wrong prefix to V1Handler, got %d", w.Code)
	}
	t.Logf("✅ V1Handler returns HTTP 400 for wrong prefix path")
}

// TestV2Handler_WrongPrefixReturnsBadRequest memastikan V2Handler
// tidak menerima request dengan prefix yang salah.
func TestV2Handler_WrongPrefixReturnsBadRequest(t *testing.T) {
	handler := http.HandlerFunc(V2Handler)

	// GET /api/v1/invoices/1001 ketika dipanggil langsung ke V2Handler
	w := performRequest(handler, http.MethodGet, "/api/v1/invoices/1001")
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected HTTP 400 for wrong prefix to V2Handler, got %d", w.Code)
	}
	t.Logf("✅ V2Handler returns HTTP 400 for wrong prefix path")
}

// TestTrailingSlash_NoSlashReturnsBadRequest memastikan route tanpa trailing slash
// tidak mengembalikan data invoice yang valid (harus 400/404, bukan 200).
func TestTrailingSlash_NoSlashReturnsBadRequest(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/invoices/", V1Handler)
	mux.HandleFunc("/api/v2/invoices/", V2Handler)

	server := httptest.NewServer(mux)
	defer server.Close()

	// GET /api/v1/invoices (tanpa trailing slash)
	// http.ServeMux dengan pattern "/api/v1/invoices/" akan redirect atau
	// menjalankan handler tergantung konfigurasi.
	// Namun penting: tidak boleh return 200 dengan id=0.
	resp, err := http.Get(server.URL + "/api/v1/invoices")
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer resp.Body.Close()

	// Pastikan tidak pernah return 200 untuk path tanpa ID
	if resp.StatusCode == http.StatusOK {
		t.Errorf("BUG: path without trailing slash should NOT return HTTP 200")
	}
	t.Logf("✅ Route /api/v1/invoices (no slash) returns HTTP %d, not 200", resp.StatusCode)

	// GET /api/v2/invoices tanpa trailing slash
	resp2, err := http.Get(server.URL + "/api/v2/invoices")
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode == http.StatusOK {
		t.Errorf("BUG: V2 path without trailing slash should NOT return HTTP 200")
	}
	t.Logf("✅ Route /api/v2/invoices (no slash) returns HTTP %d, not 200", resp2.StatusCode)
}

// TestIDValidation_ZeroAndNegativeIDs memastikan ID nol dan negatif
// dianggap invalid (400 Bad Request), bukan 404 Not Found.
func TestIDValidation_ZeroAndNegativeIDs(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/invoices/", V1Handler)
	server := httptest.NewServer(mux)
	defer server.Close()

	// ID 0 → 400 (invalid ID: must be positive)
	resp, err := http.Get(server.URL + "/api/v1/invoices/0")
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected HTTP 400 for ID=0, got %d", resp.StatusCode)
	}
	t.Logf("✅ V1 returns HTTP 400 for ID=0 (invalid: must be positive)")

	// ID negatif → 400 (invalid ID: must be positive)
	resp2, err := http.Get(server.URL + "/api/v1/invoices/-1")
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusBadRequest {
		t.Errorf("expected HTTP 400 for ID=-1, got %d", resp2.StatusCode)
	}
	t.Logf("✅ V1 returns HTTP 400 for ID=-1 (invalid: must be positive)")
}

// TestIDValidation_ValidButNotFound memastikan ID valid positif
// yang tidak ditemukan mengembalikan 404, bukan 400.
func TestIDValidation_ValidButNotFound(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/invoices/", V1Handler)
	server := httptest.NewServer(mux)
	defer server.Close()

	// ID 9999 (valid positive, but not in mock)
	resp, err := http.Get(server.URL + "/api/v1/invoices/9999")
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected HTTP 404 for valid-but-not-found ID, got %d", resp.StatusCode)
	}
	t.Logf("✅ V1 returns HTTP 404 for valid ID=9999 (not found)")
}
