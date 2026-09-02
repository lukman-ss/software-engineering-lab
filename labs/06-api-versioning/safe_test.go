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

// TestV2Contract_UsesNestedCustomer memastikan V2 menggunakan customer sebagai object.
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

	// HTTP 200
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("V2: expected HTTP 200, got %d", resp.StatusCode)
	}

	var v2Resp InvoiceV2Response
	if err := json.NewDecoder(resp.Body).Decode(&v2Resp); err != nil {
		t.Fatalf("V2 decode failed: %v", err)
	}

	// customer = object (bukan string)
	if v2Resp.Customer.Name == "" {
		t.Error("V2: customer.name tidak boleh kosong")
	}

	// customer.id = number
	if v2Resp.Customer.ID <= 0 {
		t.Error("V2: customer.id harus berupa number positif")
	}

	// customer.phone = string
	if v2Resp.Customer.Phone == "" {
		t.Error("V2: customer.phone tidak boleh kosong")
	}

	t.Log("✅ V2 contract verified: customer is object with id, name, phone")
}

// TestV1Contract_DetectsCustomerStringToObject memastikan contract test V1
// akan GAGAL bila `customer` berubah dari string menjadi object.
// Tujuan: regression guard untuk breaking change.
func TestV1Contract_DetectsCustomerStringToObject(t *testing.T) {
	// Response hipotetis yang melanggar V1 contract (customer jadi object)
	violatingResponse := []byte(`{
		"id": 1001,
		"customer": {"id": 15, "name": "Budi"},
		"total": 500000,
		"status": "PAID"
	}`)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(violatingResponse, &raw); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	// V1 contract: customer HARUS string
	var customerStr string
	err := json.Unmarshal(raw["customer"], &customerStr)
	if err == nil {
		t.Fatal("V1 contract test GAGAL mendeteksi: customer berubah dari string ke object (breaking change terlewat!)")
	}

	t.Logf("✅ V1 contract detected breaking change: %v", err)
}

// TestV1Contract_DetectsTotalNumberToString memastikan contract test V1
// akan mendeteksi bila `total` berubah dari number menjadi string.
// Tujuan: regression guard untuk breaking change tipe data.
func TestV1Contract_DetectsTotalNumberToString(t *testing.T) {
	// Response hipotetis yang melanggar V1 contract (total jadi string)
	violatingResponse := []byte(`{
		"id": 1001,
		"customer": "Budi",
		"total": "500000",
		"status": "PAID"
	}`)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(violatingResponse, &raw); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	// V1 contract: total HARUS number
	var totalNum json.Number
	err := json.Unmarshal(raw["total"], &totalNum)
	if err != nil {
		// Valid - breaking change terdeteksi karena total tidak bisa diparse sebagai number
		t.Logf("✅ V1 contract detected breaking change: %v", err)
		return
	}

	// Jika masih bisa di-parse (karena json.Number bisa menyerap string number),
	// gunakan strict type checking dari json.Unmarshal untuk float64
	var totalFloat float64
	if err := json.Unmarshal(raw["total"], &totalFloat); err != nil {
		// Valid - mendeteksi tipe total berubah
		t.Logf("✅ V1 contract detected breaking change: total berubah dari number ke string: %v", err)
		return
	}

	// Lebih strict: pastikan raw bukan literal string dengan tanda kutip
	rawStr := string(raw["total"])
	if len(rawStr) > 0 && rawStr[0] == '"' {
		t.Log("✅ V1 contract detected breaking change: total number → string (terdeteksi via raw JSON quotes)")
		return
	}

	t.Fatal("V1 contract test GAGAL mendeteksi: total berubah dari number ke string (breaking change terlewat!)")
}

// TestV1Contract_DetectsStatusStringToObject memastikan contract test V1
// akan mendeteksi bila `status` berubah dari string menjadi object.
// Tujuan: regression guard untuk breaking change tipe data.
func TestV1Contract_DetectsStatusStringToObject(t *testing.T) {
	// Response hipotetis yang melanggar V1 contract (status jadi object)
	violatingResponse := []byte(`{
		"id": 1001,
		"customer": "Budi",
		"total": 500000,
		"status": {"code": "PAID", "label": "Paid"}
	}`)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(violatingResponse, &raw); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	// V1 contract: status HARUS string
	var statusStr string
	err := json.Unmarshal(raw["status"], &statusStr)
	if err == nil {
		t.Fatal("V1 contract test GAGAL mendeteksi: status berubah dari string ke object (breaking change terlewat!)")
	}

	t.Logf("✅ V1 contract detected breaking change: status string → object: %v", err)
}
