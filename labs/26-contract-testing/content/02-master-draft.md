# Contract Testing — Consumer-Driven Contracts, CI Gate, dan Safe API Evolution

## Problem

Ketika Service A (consumer) dan Service B (provider) masing-masing memiliki test suite hijau, integrasi mereka tetap bisa pecah di produksi. Contoh konkret di laboratorium: consumer membaca field `name` dari response JSON, provider melakukan refactor ke `full_name`. Test unit consumer tetap PASS karena mock-nya belum update, test unit provider juga PASS karena logika bisnisnya benar — namun saat deploy ke produksi, consumer menerima `null` untuk field yang diharapkan, dan order gagal diproses.

Masalah ini bukan kegagalan unit test atau integration test biasa. Ini adalah kegagalan *contract*: pemahaman bersama tentang format pesan yang tidak dieksekusi sebelum deploy.

## Why This Matters

Setiap microservice atau bounded context yang berkomunikasi via HTTP/JSON menyimpan asumsi implisit tentang format pesan mitra. Asumsi ini tidak tertangkap unit test (yang menguji logika internal) maupun integration test yang biasanya bersifat end-to-end (yang lambat dan rapuh). Contract testing menutup celah ini dengan mengubah "pemahaman bersama" menjadi *executable assertion* yang berjalan di CI pipeline sebelum deploy.

Dalam laboratorium ini, hal yang diukur konkret: apakah sebuah provider memenuhi kontrak yang didefinisikan consumer sebelum provider diizinkan deploy ke lingkungan produksi.

## Mental Model

Bayangkan kontrak seperti *spec dokumen pengiriman*: consumer (pengirim paket) menyatakan "saya butuh paket berisi `id`, `status`, `nama`, `nominal`". Provider (kurir) tidak perlu tahu bisnis konsumen, tapi harus memastikan paket yang dikirim sesuai spesifikasi itu. Contract testing membuat *spesifikasi* ini executable dan diverifikasi otomatis oleh CI.

Tiga prinsip kunci:
1. **Consumer yang menentukan** ekspektasi minimal — hanya field yang benar-benar dipakai yang diverifikasi. Field tambahan di provider tidak memengaruhi verifikasi.
2. **Kontrak bukan schema JSON** — ia mencakup semantik HTTP (method, path, status, header, tipe, casing enum, error behavior). Header dideklarasikan dalam kontrak dan *verifier memvalidasi response header* terhadap ekspektasi tersebut (`internal/contract/verifier.go:90-98`), selain status code dan body.
3. **Breaking change terdeteksi sebelum production** — provider verification gagal di CI, deployment dicegah.

## Core Concept

**Consumer-Driven Contract Testing (CDC)** adalah teknik di mana consumer menulis test yang mendefinisikan ekspektasi minimalnya terhadap provider. Test ini menghasilkan *contract file* (dalam ekosistem Pact disebut *pact file*). Provider kemudian menjalankan verifikasi terhadap contract file tersebut di pipeline CI-nya: request yang sama dikirim ke provider yang sesungguhnya, dan response diverifikasi apakah memenuhi ekspektasi minimal consumer.

Secara teknis, kontrak terdiri dari satu atau lebih *interactions*, masing-masing berisi: deskripsi, provider state, request (method, path, header), dan response yang diharapkan (status, header, body dengan tipe dan casing tertentu). Provider yang memenuhi kontrak menghasilkan `VerificationResult.Passed = true`. Provider yang melanggar kontrak menghasilkan daftar error terperinci (missing field, type mismatch, enum mismatch).

## Failure Scenario

Laboratorium ini memodelkan tiga breaking change yang umum terjadi:

1. **Enum casing change**: `status` dari `"IN_PROGRESS"` menjadi `"in_progress"`. JSON-nya valid secara sintaksis, tetapi consumer yang mengharapkan casing persis akan gagal.
2. **Field rename**: `customer.name` menjadi `customer.full_name`. Field lama tidak ada di response, consumer mendapatkan nilai default/undefined.
3. **Primitive type mutation**: `total` dari integer `150000` menjadi string `"150000"`. JSON-nya valid, tetapi consumer yang melakukan operasi aritmatika pada `total` akan gagal.

Ketika ketiga perubahan ini dilakukan provider tanpa koordinasi, contract verification di CI menghasilkan:
- Missing expected field `'customer.name'`
- Type mismatch pada `'total'` (expected `json.Number`, got `string`)
- Value mismatch pada `'status'` (expected `"IN_PROGRESS"`, got `"in_progress"`)

Hasilnya: `VerificationResult.Passed = false`. CI gate menolak deployment. Produksi terlindungi.

## How It Works

Alur kontrak dalam laboratorium:

1. **Consumer test** menjalankan mock server, menyatakan ekspektasi minimal, dan menghasilkan contract JSON.
2. **Contract file** disimpan sebagai artefak yang merepresentasikan kebutuhan consumer.
3. **Provider CI** mengambil contract tersebut dan menjalankan verifikasi terhadap provider server yang sesungguhnya (via `httptest.Server` dalam laboratorium, atau instance produksi dalam pipeline nyata).
4. **Verifier** mengirim request sesuai kontrak ke provider, menerima response, membandingkan **status code, response header, dan body** secara *subset*: hanya field yang dideklarasikan consumer yang diperiksa; response header divalidasi terhadap ekspektasi kontrak (`verifier.go:90-98`).
5. **Hasil**: pass → deployment diizinkan; fail → deployment diblokir dengan daftar error terperinci.

Mekanisme verifikasi inti menggunakan *recursive map comparison* dengan `json.Number` untuk menjaga presisi tipe numerik, dan `reflect` untuk membandingkan tipe primitif non-numerik.

## Architecture

```text
[Consumer (MobileApp Client)]
             │
             │ (generates contract via GenerateMobileContract)
             ▼
[contract.json — minimal expectations]
             │
             │ (verifier fetches contract)
             ▼
[Contract Verifier — Verify(baseURL, contract)]
             │
     ┌───────┼────────┬──────────┐
     ▼       ▼        ▼          ▼
[Provider V1] [Breaking Provider] [Dual Provider /v1 & /v2]
     │              │                  │
     ▼              ▼                  ▼
  PASS          FAIL (blocked)      PASS (V1 route)
```

Komponen inti:
- `internal/consumer/client.go`: `GenerateMobileContract()` membangun contract JSON minimal dan `MobileOrderClient` yang berinteraksi dengan provider.
- `internal/contract/verifier.go`: `Verifier.Verify()` menjalankan verifikasi; `diffValues()` membandingkan expected vs actual secara recursive.
- `internal/model/order.go`: DTO untuk `OrderResponseV1`, `OrderResponseBreaking`, dan `OrderResponseV2`.
- `internal/provider/server.go`: `ProviderV1`, `ProviderBreaking`, `ProviderDual` sebagai HTTP handler.
- `tests/contract_test.go`: 7 test mencakup generasi, verifikasi sukses, verifikasi gagal, dual provider, konkurensi, header validation + error branches, dan v2 direct assertion.
- `cmd/demo/main.go`: Eksekusi 4 stage demo.

## Implementation

Laboratorium diimplementasikan dalam Go standar library — tidak bergantung pada eksternal Pact daemon, Ruby binary, atau proses jarak jauh. `net/http` dan `httptest` digunakan untuk simulasi server. Verifier adalah engine custom yang melakukan subset matching: hanya field yang dideklarasikan consumer yang diverifikasi; field tambahan provider (seperti `notes`, `created_at`) diabaikan.

Keputusan ini disengaja untuk reproduktibilitas lokal dan CI tanpa dependensi runtime eksternal. Komprominya: verifier tidak mengimplementasikan seluruh spesifikasi Pact (misalnya provider state callbacks, message pact, broker client).

## Code Walkthrough

### Consumer Contract Definition

File: `internal/consumer/client.go` — fungsi `GenerateMobileContract()`

```go
func GenerateMobileContract() *contract.Contract {
	return &contract.Contract{
		Consumer: "MobileApp",
		Provider: "OrderService",
		Interactions: []contract.Interaction{
			{
				Description:   "A request for order details by ID",
				ProviderState: "Order ORD-123 exists and is IN_PROGRESS",
				Request: contract.RequestDefinition{
					Method: http.MethodGet,
					Path:   "/v1/orders/ORD-123",
				},
				Response: contract.ResponseDefinition{
					Status: http.StatusOK,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
					Body: map[string]interface{}{
						"id":     "ORD-123",
						"status": "IN_PROGRESS",
						"customer": map[string]interface{}{
							"name": "Budi Santoso",
						},
						"total": json.Number("150000"),
					},
				},
			},
		},
	}
}
```

Kontrak menyatakan ekspektasi *minimal*: `id`, `status`, `customer.name`, `total`. Provider bebas menambahkan field lain (seperti `notes` yang ada di V1 response) tanpa memengaruhi verifikasi.

### Provider Implementations

File: `internal/provider/server.go` — tiga handler (`ProviderV1`, `ProviderBreaking`, `ProviderDual`)

```go
// ProviderV1 — memenuhi kontrak
func (p *ProviderV1) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// ...
	resp := model.OrderResponseV1{ID: id, Status: "IN_PROGRESS", Customer: model.CustomerResponseV1{Name: "Budi Santoso"}, Total: 150000}
	// ...
}

// ProviderBreaking — 3 mutasi breaking
func (p *ProviderBreaking) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// ...
	resp := model.OrderResponseBreaking{ID: id, Status: "in_progress", Customer: model.CustomerResponseBreaking{FullName: "Budi Santoso"}, Total: "150000"}
	// ...
}
```

`ProviderDual` mendukung routing `/v1/orders/{id}` (respons V1 compliant) dan `/v2/orders/{id}` (respons V2 dengan field tambahan `currency`). V1 route tetap lolos verifikasi consumer. **Catatan**: Endpoint V2 belum memiliki consumer contract atau verifikasi test — hanya route V1 yang diuji dalam lab (lihat GAP-02 di `engineering-audit-opensource/06-verdict.md`).

### Verifier Engine

File: `internal/contract/verifier.go` — `Verifier.Verify()` dan `diffValues()`

```go
func (v *Verifier) Verify(baseURL string, c *Contract) VerificationResult {
	result := VerificationResult{Passed: true}
	for _, interaction := range c.Interactions {
		targetURL := strings.TrimRight(baseURL, "/") + interaction.Request.Path
		req, _ := http.NewRequest(interaction.Request.Method, targetURL, nil)
		resp, err := v.Client.Do(req)
		// ... status code check ...
		var actualBody map[string]interface{}
		decoder := json.NewDecoder(bytes.NewReader(bodyBytes))
		decoder.UseNumber()
		decoder.Decode(&actualBody)
		diffs := diffValues("", interaction.Response.Body, actualBody)
		if len(diffs) > 0 { result.Passed = false; /* append errors */ }
	}
	return result
}
```

`diffValues()` melakukan *recursive map comparison*: jika expected adalah map, dicek keberadaan setiap key di actual; jika expected adalah `json.Number`, actual juga harus `json.Number` (mencegah string numeric lolos sebagai integer). Ini memungkinkan deteksi tipe primitif yang berubah tanpa bergantung pada skema JSON formal.

## What the Tests Prove

Lima test dalam `tests/contract_test.go` membuktikan perilaku konkret berikut:

- **`TestConsumerContractGeneration`**: Kontrak yang dihasilkan memiliki consumer `"MobileApp"`, provider `"OrderService"`, dan 1 interaksi dengan path `/v1/orders/ORD-123`.
- **`TestProviderV1_ContractVerification_Success`**: Provider V1 memenuhi kontrak consumer secara penuh. `Verifier.Verify()` menghasilkan `Passed=true`. `FetchOrder` berhasil dengan nama, status, dan total sesuai kontrak.
- **`TestProviderBreaking_ContractVerification_Fails`**: Provider dengan 3 breaking change menghasilkan `Passed=false` dengan minimal 3 error (missing `customer.name`, type mismatch `total`, enum mismatch `status`). `FetchOrder` gagal karena consumer mendeteksi pelanggaran kontrak.
- **`TestProviderDual_ContractVerification_Success`**: Provider dengan dual routing `/v1` + `/v2` tetap lolos verifikasi kontrak V1. `FetchOrder` pada route `/v1` berhasil.
- **`TestConcurrentContractVerification`**: 20 goroutine menjalankan verifikasi secara bersamaan tanpa data race — membuktikan verifier thread-safe.

## Recovery / Rollback

Laboratorium menggambarkan dua jalur recovery:

1. **Breaking change terdeteksi di CI (dicegah)**: Provider verification gagal, deployment diblokir. Provider merancang perubahan dengan *expand → migrate → contract* pattern: tambahkan field baru atau endpoint baru, update consumer, baru hapus field lama.
2. **Dual provider (V1 + V2)**: Provider mempertahankan endpoint V1 yang compliant sambil menambahkan endpoint V2 dengan skema baru. Consumer dapat bermigrasi ke V2 sesuai ritme masing-masing tanpa memecah kontrak yang ada.

## Production Considerations

Laboratorium menggunakan `httptest.Server` sebagai provider — dalam pipeline nyata, verifikasi mengarah ke instance provider yang berjalan. Pertimbangan produksi:

- **Provider state management**: Dalam lab, provider state ditangani melalui routing deterministik. Dalam produksi, provider perlu menyiapkan state tertentu (misalnya "order ORD-123 exists") sebelum verifikasi — ini biasanya dilakukan melalui API state-setup atau fixture database.
- **Contract Broker**: Dalam lab, kontrak dipertukarkan secara in-memory. Dalam produksi, Pact Broker (atau implementasi sejenis) menyimpan kontrak dan hasil verifikasi, memungkinkan `can-i-deploy` check lintas pasangan consumer-provider.
- **Timeout dan konektivitas**: Verifier laboratorium menggunakan `http.Client` dengan timeout 5 detik (`verifier.go:55`) — cocok untuk `httptest` maupun pipeline produksi. Dalam pipeline nyata, timeout mungkin perlu disesuaikan dengan latensi jaringan.
- **Header assertion**: Dalam lab, `Content-Type` header dideklarasikan dalam kontrak dan diverifikasi oleh verifier (`verifier.go:90-98`). Gate produksi dapat menambah assertion tambahan sesuai kebutuhan.
- **Asynchronous systems**: Lab membatasi cakupan ke REST HTTP. Sistem event-driven (Kafka, RabbitMQ, webhooks) memerlukan Message Pact dengan prinsip serupa namun di level message payload.

## Common Mistakes

1. **Over-specification**: Mengujikan aturan validasi bisnis dalam kontrak (misalnya "username max 20 chars") membuat kontrak terlalu rapuh. Provider yang melonggarkan validasi (misalnya memperbolehkan nomor) akan memecah kontrak padahal tidak ada dampak pada consumer. Kontrak harus fokus pada *format pesan* dan *error behavior* (misalnya "400 dengan pesan error"), bukan *mengapa* error terjadi.
2. **Tidak membedakan additive vs breaking**: Senior engineer membedakan additive change (tambah field, aman) dan breaking change (rename, hapus, ganti tipe) sebelum merge. Contract testing membuat perbedaan ini terdeteksi otomatis.
3. **Kontrak terlalu ketat**: Consumer meminta field yang tidak dibutuhkan; provider terblokir menambah fitur apa pun. Prinsip CDC: consumer hanya mendeklarasikan apa yang dibutuhkan.
4. **Mengandalkan JSON Schema saja**: Skema JSON valid secara sintaksis namun tidak menangkap semantik (integer vs string, casing enum). Contract testing menangkap ini lewat *concrete examples*.
5. **Menganggap contract testing menggantikan semua test**: Contract test menggantikan *sekelas* integration test yang memvalidasi penggunaan API, bukan *inti* logika bisnis. Provider tetap memerlukan unit test dan functional test. Piramida yang sehat: banyak unit test → contract tests → beberapa integration test → sedikit critical E2E test.

## Case Study

**Lab Order Service — Mobile App Consumer**

- **Consumer**: Mobile App yang membutuhkan `id`, `status`, `customer.name`, `total` dari endpoint order.
- **Provider**: Order Service yang melayani `/v1/orders/{id}`.
- **Kontrak**: Satu interaksi — `GET /v1/orders/ORD-123` dengan expected body `{id, status:"IN_PROGRESS", customer:{name}, total:150000}`.

**Stage 1 — Compliant Provider**: Provider V1 merespons sesuai kontrak. Verifikasi PASS. CI gate mengizinkan deployment.

**Stage 2 — Breaking Provider**: Provider menerapkan 3 perubahan tanpa koordinasi: casing status diubah, `customer.name` diganti `customer.full_name`, `total` diubah ke string. Verifikasi GAGAL dengan 3 error terperinci. CI gate memblokir deployment — produksi terlindungi.

**Stage 3 — Dual Provider (Safe Evolution)**: Provider menambahkan endpoint `/v2` dengan skema baru (termasuk field tambahan `currency`) sambil mempertahankan `/v1` yang compliant. Verifikasi kontrak V1 PASS. Consumer dapat bermigrasi ke V2 secara independen. CI gate mengizinkan deployment terpisah.

## Checklist

- [ ] Consumer mendefinisikan kontrak minimal (hanya field yang dibutuhkan).
- [ ] Kontrak mencakup method, path, status, header, body dengan tipe dan casing.
- [ ] Provider CI menjalankan verifikasi kontrak sebelum deploy.
- [ ] Verifikasi gagal → deployment diblokir dengan error terperinci.
- [ ] Additive changes (tambah field opsional) tidak memecah kontrak yang ada.
- [ ] Breaking changes menggunakan expand → migrate → contract pattern.
- [ ] Dual-version provider memungkinkan migrasi consumer yang independen.
- [ ] Over-specification (validasi bisnis) dihindari dalam kontrak.
- [ ] Contract test melengkapi — tidak menggantikan — unit test dan functional test.
- [ ] Hasil verifikasi dipublikasikan ke contract broker untuk `can-i-deploy` check.

## Key Takeaways

1. Contract testing menutup celah antara test unit/integration yang hijau dan kegagalan integrasi produksi.
2. Consumer-driven contracts menggeser definisi kontrak ke consumer — hanya field yang digunakan yang diuji, membebaskan provider untuk berkembang.
3. Kontrak mencakup semantik HTTP lengkap (status, header, tipe, casing enum, error behavior) — bukan sekadar JSON Schema.
4. Breaking changes terdeteksi di CI *sebelum* deploy; expand/contract pattern membuat evolusi API yang aman.
5. Provider verification yang gagal adalah sinyal bukan failure — ia mencegah production outage.
6. Kontrak harus minimal dan fokus pada format pesan, bukan aturan validasi bisnis (anti-pattern over-specification).
7. Contract testing berlaku untuk REST dan event-driven (Message Pact).
8. Kontrak menggeser piramida test: lebih banyak contract test, lebih sedikit E2E yang rapuh.
9. Implementasi dalam laboratorium adalah *minimal viable engine* — untuk produksi diperlukan broker, provider state management, timeout, dan message pact.
10. Over-specification adalah musuh utama CDC: kontrak yang terlalu ketat menghambat evolusi provider.

## Sources

- Pact Docs: Introduction, How Pact works, Contract Tests vs Functional Tests, FAQ, CI/CD Setup Guide
- Martin Fowler: Consumer-Driven Contracts (2006), ContractTest bliki (2011)
- Danilo Sato: Parallel Change (expand/contract pattern) (2014)
- Spring Cloud Contract (archived, Jul 2026) — konteks historis
- AsyncAPI Initiative — positioning event-driven contract analogue
- Laboratorium: internal source code, tests, dan execution output lab 26-contract-testing
