package api_versioning

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestBreakingChange_LegacyClientFails membuktikan bahwa bukan semua perubahan
// backend yang berhasil compile otomatis menjadi backward compatible.
//
// Mental model: HTTP 200 ≠ backward compatible
func TestBreakingChange_LegacyClientFails(t *testing.T) {
	// Setup: server yang mengembalikan response V2 (customer sebagai object)
	server := httptest.NewServer(http.HandlerFunc(UnsafeHandler))
	defer server.Close()

	// Legacy client mencoba decode ke struct yang mengharapkan customer = string
	url := server.URL + "/api/invoices/1001"

	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("failed to connect to server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected HTTP 200, got %d", resp.StatusCode)
	}

	// Decode JSON body seperti legacy client menggunakan helper
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	_, err = ParseLegacyInvoice(body)

	// ASSERTION: Decode HARUS gagal karena `customer` sudah bukan string
	if err == nil {
		t.Error("expected legacy client decode to FAIL, but it succeeded")
		t.Error("This indicates a BREAKING CHANGE went undetected!")
	}

	t.Logf("✅ Breaking change confirmed: legacy client fails to decode: %v", err)
	t.Log("Server returns HTTP 200 with valid JSON, but client cannot parse it")
	t.Log("This demonstrates: HTTP 200 != backward compatible")
}

// TestUnsafeHandler_GETValidInvoice memastikan unsafe endpoint
// mengembalikan HTTP 200 untuk invoice yang valid.
func TestUnsafeHandler_GETValidInvoice(t *testing.T) {
	handler := http.HandlerFunc(UnsafeHandler)
	w := performRequest(handler, http.MethodGet, "/api/invoices/1001")

	if w.Code != http.StatusOK {
		t.Errorf("expected HTTP 200, got %d", w.Code)
	}

	// Verify Content-Type
	ct := w.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got '%s'", ct)
	}

	t.Logf("✅ UnsafeHandler returns HTTP 200 for GET /api/invoices/1001")
}

// TestUnsafeHandler_GETInvalidID memastikan unsafe endpoint
// mengembalikan HTTP 400 untuk ID non-numeric.
func TestUnsafeHandler_GETInvalidID(t *testing.T) {
	handler := http.HandlerFunc(UnsafeHandler)
	w := performRequest(handler, http.MethodGet, "/api/invoices/abc")

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected HTTP 400, got %d", w.Code)
	}

	t.Logf("✅ UnsafeHandler returns HTTP 400 for non-numeric ID")
}

// TestUnsafeHandler_GETMissingID memastikan unsafe endpoint
// mengembalikan HTTP 400 untuk path tanpa ID.
func TestUnsafeHandler_GETMissingID(t *testing.T) {
	handler := http.HandlerFunc(UnsafeHandler)
	w := performRequest(handler, http.MethodGet, "/api/invoices/")

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected HTTP 400, got %d", w.Code)
	}

	t.Logf("✅ UnsafeHandler returns HTTP 400 for missing ID")
}

// TestUnsafeHandler_GETExtraPath memastikan unsafe endpoint
// mengembalikan HTTP 400 untuk path dengan extra segment.
func TestUnsafeHandler_GETExtraPath(t *testing.T) {
	handler := http.HandlerFunc(UnsafeHandler)
	w := performRequest(handler, http.MethodGet, "/api/invoices/1001/extra")

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected HTTP 400, got %d", w.Code)
	}

	t.Logf("✅ UnsafeHandler returns HTTP 400 for extra path segment")
}

// TestUnsafeHandler_WrongPrefixNeverFallsThrough memastikan path dengan prefix salah
// tidak pernah dianggap sebagai ID 0 yang valid.
func TestUnsafeHandler_WrongPrefixNeverFallsThrough(t *testing.T) {
	handler := http.HandlerFunc(UnsafeHandler)

	// GET /wrong/path/1001 - prefix bukan "/api/invoices/"
	w := performRequest(handler, http.MethodGet, "/wrong/path/1001")
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected HTTP 400 for wrong prefix, got %d", w.Code)
	}
	t.Logf("✅ UnsafeHandler returns HTTP 400 for wrong prefix path")

	// Pastikan tidak pernah mengembalikan 200 untuk path yang salah
	if w.Code == http.StatusOK {
		t.Error("BUG: wrong prefix path should NOT return HTTP 200")
	}
}
