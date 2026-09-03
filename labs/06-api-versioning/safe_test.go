package api_versioning

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newVersionedMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/invoices", V1Handler)
	mux.HandleFunc("/api/v1/invoices/", V1Handler)

	mux.HandleFunc("/api/v2/invoices", V2Handler)
	mux.HandleFunc("/api/v2/invoices/", V2Handler)

	return mux
}

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
// Menggunakan legacy consumer model (LegacyInvoice) untuk bukti compatibility.
func TestSafeVersioning_LegacyClientSucceedsOnV1(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/invoices/", V1Handler)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/invoices/1001", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected HTTP 200, got %d", w.Code)
	}

	body := w.Body.Bytes()

	// Legacy client decode menggunakan model consumer lama
	legacyResponse, err := ParseLegacyInvoice(body)
	if err != nil {
		t.Fatalf("legacy client SHOULD succeed on V1: %v", err)
	}

	// Assert exact business contract
	if legacyResponse.ID != 1001 {
		t.Errorf("ID: expected 1001, got %d", legacyResponse.ID)
	}
	if legacyResponse.Customer != "Budi" {
		t.Errorf("Customer: expected 'Budi', got '%s'", legacyResponse.Customer)
	}
	if legacyResponse.Total != 500000 {
		t.Errorf("Total: expected 500000, got %d", legacyResponse.Total)
	}
	if legacyResponse.Status != "PAID" {
		t.Errorf("Status: expected 'PAID', got '%s'", legacyResponse.Status)
	}

	t.Log("✅ Legacy client successfully reads from V1 API")
}

// TestSafeVersioning_NewClientSucceedsOnV2 membuktikan client baru dapat
// menggunakan V2 API dengan customer sebagai object.
// Menggunakan typed response untuk bukti client-level success.
func TestSafeVersioning_NewClientSucceedsOnV2(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/invoices/", V2Handler)

	req := httptest.NewRequest(http.MethodGet, "/api/v2/invoices/1001", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected HTTP 200, got %d", w.Code)
	}

	var v2Response InvoiceV2Response
	err := json.NewDecoder(w.Body).Decode(&v2Response)
	if err != nil {
		t.Fatalf("new client SHOULD succeed on V2: %v", err)
	}

	// Assert exact business contract
	if v2Response.ID != 1001 {
		t.Errorf("ID: expected 1001, got %d", v2Response.ID)
	}
	if v2Response.Customer.ID != 15 {
		t.Errorf("Customer.ID: expected 15, got %d", v2Response.Customer.ID)
	}
	if v2Response.Customer.Name != "Budi" {
		t.Errorf("Customer.Name: expected 'Budi', got '%s'", v2Response.Customer.Name)
	}
	if v2Response.Customer.Phone != "08123" {
		t.Errorf("Customer.Phone: expected '08123', got '%s'", v2Response.Customer.Phone)
	}
	if v2Response.Total != 500000 {
		t.Errorf("Total: expected 500000, got %d", v2Response.Total)
	}
	if v2Response.Status != "PAID" {
		t.Errorf("Status: expected 'PAID', got '%s'", v2Response.Status)
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

	req1 := httptest.NewRequest(http.MethodGet, "/api/v1/invoices/1001", nil)
	w1 := httptest.NewRecorder()
	mux.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Fatalf("V1: expected HTTP 200, got %d", w1.Code)
	}

	var v1Resp InvoiceV1Response
	if err := json.NewDecoder(w1.Body).Decode(&v1Resp); err != nil {
		t.Fatalf("V1 decode failed: %v", err)
	}

	req2 := httptest.NewRequest(http.MethodGet, "/api/v2/invoices/1001", nil)
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Fatalf("V2: expected HTTP 200, got %d", w2.Code)
	}

	var v2Resp InvoiceV2Response
	if err := json.NewDecoder(w2.Body).Decode(&v2Resp); err != nil {
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

	req := httptest.NewRequest(http.MethodGet, "/api/v1/invoices/1001", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("V1: expected HTTP 200 without X-Tenant-ID, got %d", w.Code)
	}

	t.Log("✅ V1 endpoint tidak membutuhkan X-Tenant-ID header - backward compatible")
}

// TestV1Contract_RegressionGuaranteeAfterV2Registration tambahan verifikasi
// untuk memastikan tidak ada breaking change pada V1 pasca V2 registration.
func TestV1Contract_RegressionGuaranteeAfterV2Registration(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/invoices/", V1Handler)
	mux.HandleFunc("/api/v2/invoices/", V2Handler)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/invoices/1001", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("V1: expected HTTP 200, got %d", w.Code)
	}

	ct := w.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("V1: expected Content-Type 'application/json', got '%s'", ct)
	}

	body := w.Body.Bytes()
	assertV1WireContract(t, body)

	t.Log("✅ V1 contract regression guarantee: V2 registration tidak mengubah V1 output")
}

// TestV1Contract_RemainsBackwardCompatible melindungi kontrak V1 agar tidak
// secara tidak sengaja mengubah bentuk customer dari string menjadi object.
func TestV1Contract_RemainsBackwardCompatible(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/invoices/", V1Handler)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/invoices/1001", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("V1: expected HTTP 200, got %d", w.Code)
	}

	ct := w.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("V1: expected Content-Type 'application/json', got '%s'", ct)
	}

	body := w.Body.Bytes()
	assertV1WireContract(t, body)

	t.Log("✅ V1 contract protected: customer remains string")
}

// TestV2Contract_UsesNestedCustomer memastikan V2 menggunakan customer sebagai object
// dengan semua field yang tepat. Memverifikasi JSON wire contract secara eksplisit
// menggunakan raw JSON untuk detection breaking changes.
func TestV2Contract_UsesNestedCustomer(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/invoices/", V2Handler)

	req := httptest.NewRequest(http.MethodGet, "/api/v2/invoices/1001", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("V2: expected HTTP 200, got %d", w.Code)
	}

	ct := w.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("V2: expected Content-Type 'application/json', got '%s'", ct)
	}

	body := w.Body.Bytes()
	assertV2WireContract(t, body)

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
// adalah invalid collection request pada lab ini dan harus menghasilkan 400 Bad Request.
func TestTrailingSlash_NoSlashReturnsBadRequest(t *testing.T) {
	mux := newVersionedMux()

	req1 := httptest.NewRequest(http.MethodGet, "/api/v1/invoices", nil)
	w1 := httptest.NewRecorder()
	mux.ServeHTTP(w1, req1)

	if w1.Code != http.StatusBadRequest {
		t.Fatalf("expected HTTP 400, got %d", w1.Code)
	}

	req2 := httptest.NewRequest(http.MethodGet, "/api/v2/invoices", nil)
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, req2)

	if w2.Code != http.StatusBadRequest {
		t.Fatalf("expected HTTP 400, got %d", w2.Code)
	}
}

// TestIDValidation_ZeroAndNegativeIDs memastikan ID nol dan negatif
// dianggap invalid (400 Bad Request), bukan 404 Not Found.
func TestIDValidation_ZeroAndNegativeIDs(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/invoices/", V1Handler)

	req1 := httptest.NewRequest(http.MethodGet, "/api/v1/invoices/0", nil)
	w1 := httptest.NewRecorder()
	mux.ServeHTTP(w1, req1)

	if w1.Code != http.StatusBadRequest {
		t.Errorf("expected HTTP 400 for ID=0, got %d", w1.Code)
	}
	t.Logf("✅ V1 returns HTTP 400 for ID=0 (invalid: must be positive)")

	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/invoices/-1", nil)
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, req2)

	if w2.Code != http.StatusBadRequest {
		t.Errorf("expected HTTP 400 for ID=-1, got %d", w2.Code)
	}
	t.Logf("✅ V1 returns HTTP 400 for ID=-1 (invalid: must be positive)")
}

// TestIDValidation_ValidButNotFound memastikan ID valid positif
// yang tidak ditemukan mengembalikan 404, bukan 400.
func TestIDValidation_ValidButNotFound(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/invoices/", V1Handler)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/invoices/9999", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected HTTP 404 for valid-but-not-found ID, got %d", w.Code)
	}
	t.Logf("✅ V1 returns HTTP 404 for valid ID=9999 (not found)")
}

// TestHandlers_MethodNotAllowed memverifikasi semua handler mengembalikan
// 405 Method Not Allowed dengan Allow: GET dan Content-Type: application/json
// untuk method selain GET.
func TestHandlers_MethodNotAllowed(t *testing.T) {
	tests := []struct {
		name    string
		handler http.Handler
		path    string
	}{
		{"V1Handler", http.HandlerFunc(V1Handler), "/api/v1/invoices/1001"},
		{"V2Handler", http.HandlerFunc(V2Handler), "/api/v2/invoices/1001"},
		{"UnsafeHandler", http.HandlerFunc(UnsafeHandler), "/api/invoices/1001"},
		{"AdditiveHandler", http.HandlerFunc(AdditiveHandler), "/api/invoices/1001"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := performRequest(tt.handler, http.MethodPost, tt.path)

			if w.Code != http.StatusMethodNotAllowed {
				t.Errorf("expected HTTP 405, got %d", w.Code)
			}

			allow := w.Header().Get("Allow")
			if allow != http.MethodGet {
				t.Errorf("expected Allow header 'GET', got '%s'", allow)
			}

			ct := w.Header().Get("Content-Type")
			if ct != "application/json" {
				t.Errorf("expected Content-Type 'application/json', got '%s'", ct)
			}

			// Verify error JSON field exists
			body := w.Body.Bytes()
			var resp map[string]interface{}
			if err := json.Unmarshal(body, &resp); err != nil {
				t.Errorf("expected JSON error response, got: %s", string(body))
			} else if _, ok := resp["error"]; !ok {
				t.Errorf("expected JSON error field 'error', got: %v", resp)
			}

			t.Logf("✅ %s returns 405 + Allow: GET + Content-Type: application/json", tt.name)
		})
	}
}
