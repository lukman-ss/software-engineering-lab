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
	body, _ := io.ReadAll(resp.Body)
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
