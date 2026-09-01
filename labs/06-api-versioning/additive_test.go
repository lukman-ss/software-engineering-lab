package api_versioning

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestAdditiveField_LegacyClientStillWorks membuktikan bahwa penambahan field
// baru pada response biasanya TIDAK menjadi breaking change, ASAL client
// toleran terhadap unknown field dan tidak menggunakan strict schema validation.
//
// Penting: Penambahan field SELALU aman? Tidak. Namun JIKA consumer
// mengabaikan unknown field secara default (seperti Go json.Unmarshal),
// maka perubahan ini backward compatible.
func TestAdditiveField_LegacyClientStillWorks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(AdditiveHandler))
	defer server.Close()

	// Legacy client tetap menggunakan struct lama yang tidak punya field `currency`
	resp, err := http.Get(server.URL + "/?id=1001")
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected HTTP 200, got %d", resp.StatusCode)
	}

	// Legacy client decode ke struct LAMA (tanpa field currency)
	var legacyResponse LegacyInvoice
	err = json.NewDecoder(resp.Body).Decode(&legacyResponse)
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
	t.Log("INI adalah contoh additive change yang TIDK menjadi breaking change")
	t.Log("Catatan: Penambahan field aman HANYA bila consumer tolerant unknown field")
}