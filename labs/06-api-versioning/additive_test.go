package api_versioning

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// responseBody membaca body dari ResponseRecorder untuk assertion.
func responseBody(t *testing.T, recorder *httptest.ResponseRecorder) []byte {
	t.Helper()
	return recorder.Body.Bytes()
}

// TestAdditiveField_LegacyClientStillWorks membuktikan bahwa penambahan field
// baru pada response biasanya TIDAK menjadi breaking change, ASAL client
// toleran terhadap unknown field dan tidak menggunakan strict schema validation.
//
// Penting: Penambahan field SELALU aman? Tidak. Namun JIKA consumer
// mengabaikan unknown field secara default (seperti Go json.Unmarshal),
// maka perubahan ini backward compatible.
func TestAdditiveField_LegacyClientStillWorks(t *testing.T) {
	handler := http.HandlerFunc(AdditiveHandler)

	// Test valid request
	w := performRequest(handler, http.MethodGet, "/api/invoices/1001")
	if w.Code != http.StatusOK {
		t.Errorf("expected HTTP 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got '%s'", ct)
	}

	body := responseBody(t, w)

	// 1. Verify wire response memiliki additive field
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}

	// Assert currency field exists and is string "IDR"
	currencyRaw, ok := raw["currency"]
	if !ok {
		t.Error("ADDITIVE FIELD MISSING: 'currency' not found in response")
	} else {
		var currency string
		if err := json.Unmarshal(currencyRaw, &currency); err != nil {
			t.Errorf("ADDITIVE FIELD INVALID: 'currency' harus string, got error: %v", err)
		} else if currency != "IDR" {
			t.Errorf("ADDITIVE FIELD WRONG VALUE: expected 'IDR', got '%s'", currency)
		} else {
			t.Log("✅ Wire response contains additive field: currency = \"IDR\"")
		}
	}

	// 2. Legacy client decode ke struct LAMA (tanpa field currency)
	var legacyResponse LegacyInvoice
	err := json.Unmarshal(body, &legacyResponse)
	if err != nil {
		t.Fatalf("legacy client SHOULD succeed even with unknown field: %v", err)
	}

	// Verify data yang diketahui tetap ada
	if legacyResponse.Customer != "Budi" {
		t.Errorf("expected Customer='Budi', got '%s'", legacyResponse.Customer)
	}
	if legacyResponse.Total != 500000 {
		t.Errorf("expected Total=500000, got %d", legacyResponse.Total)
	}

	t.Log("✅ Legacy client successfully decodes response with unknown `currency` field")
	t.Log("Field `currency` diabaikan oleh decoder karena tidak ada di struct")
	t.Log("INI adalah contoh additive change yang tidak menjadi breaking change")
}

// TestAdditiveHandler_MethodNotAllowed memastikan additive endpoint
// mengembalikan HTTP 405 untuk method POST dengan Allow header.
func TestAdditiveHandler_MethodNotAllowed(t *testing.T) {
	handler := http.HandlerFunc(AdditiveHandler)
	w := performRequest(handler, http.MethodPost, "/api/invoices/1001")

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected HTTP 405, got %d", w.Code)
	}

	// Verify Allow header
	allow := w.Header().Get("Allow")
	if allow != http.MethodGet {
		t.Errorf("expected Allow header 'GET', got '%s'", allow)
	}

	// Verify Content-Type
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got '%s'", ct)
	}

	t.Log("✅ AdditiveHandler returns HTTP 405 with Allow: GET for POST method")
}

// TestAdditiveHandler_WrongPrefixNotValid memastikan path dengan prefix salah
// tidak dianggap sebagai ID 0 yang valid.
func TestAdditiveHandler_WrongPrefixNotValid(t *testing.T) {
	handler := http.HandlerFunc(AdditiveHandler)
	w := performRequest(handler, http.MethodGet, "/wrong/path/1001")

	if w.Code == http.StatusOK {
		t.Error("BUG: wrong prefix path should NOT return HTTP 200")
	}

	t.Logf("✅ AdditiveHandler returns HTTP %d for wrong prefix path (not 200)", w.Code)
}
