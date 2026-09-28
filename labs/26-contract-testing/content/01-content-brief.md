# Content Brief

Topic: Contract Testing — Consumer-Driven Contracts (CDC), CI Deployment Gate, dan Safe API Evolution tanpa E2E rapuh.

Target Reader: Backend engineer, platform engineer, dan tech lead yang mengelola 2+ service terpisah dengan boundary HTTP/JSON (REST). Pembaca sudah familiar dengan unit/integration test, tetapi sering melihat integrasi pecah di produksi meski semua test hijau.

Problem: Test tiap service bisa PASS namun integrasi gagal karena provider mengubah kontrak tanpa sepengetahuan consumer (rename field, ganti casing enum, ubah tipe primitif). E2E lintas service mahal, lambat, dan rapuh. Perubahan breaking lolos ke produksi karena tidak ada gate yang memvalidasi ekspektasi consumer terhadap provider sebelum deploy.

Core Mental Model: Kontrak = pemahaman bersama yang dieksekusi (executable) tentang pesan request/response yang benar-benar dipakai consumer. Consumer mendefinisikan ekspektasi minimal → provider memverifikasi di CI sebelum deploy (pass = boleh deploy, fail = blokir). Additive (tambah field) biasanya aman; rename/hapus/ganti tipe tanpa migrasi = breaking dan harus lewat expand → migrate → contract.

Approved Research Status: APPROVED (2026-09-26, 0 unsupported claims, 0 contradictions, 2 non-blocking gaps).

Approved Engineering Status: APPROVED (2026-09-26, 5/5 tests PASS, race PASS. Lab scope: Go stdlib, single-resource order contract, in-memory contract exchange).

Main Concepts: CDC vs provider-driven contracts; contract by example (Pact file) vs JSON Schema; test double vs real service; HTTP semantics dalam kontrak (method, path, status, header, body + enum/type/error); anti-pattern over-specification (jangan uji aturan validasi bisnis di kontrak); expand/contract untuk breaking change aman; Pact Broker & CI/CD gate (can-i-deploy) & independent deployability; Message Pact untuk event-driven (Kafka/RabbitMQ); piramida test setelah contract testing.

Verified Behaviors (`tests/contract_test.go`):
- `TestConsumerContractGeneration` (`:14`) — consumer MobileApp menghasilkan kontrak 1 interaksi `GET /v1/orders/ORD-123` dengan body minimal `{id, status, customer.name, total}`.
- `TestProviderV1_ContractVerification_Success` (`:28`) — Provider V1 memenuhi kontrak; verifier `Passed=true`; `FetchOrder` berhasil.
- `TestProviderBreaking_ContractVerification_Fails` (`:51`) — Breaking provider (3 mutasi: `IN_PROGRESS`→`in_progress`, `customer.name`→`customer.full_name`, `total` int→string) menghasilkan `Passed=false` dengan ≥3 error terperinci; `FetchOrder` gagal.
- `TestProviderDual_ContractVerification_Success` (`:75`) — Dual provider (`/v1` + `/v2`) tetap lolos verifikasi V1; `FetchOrder` berhasil.
- `TestConcurrentContractVerification` (`:97`) — 20 goroutine memverifikasi bersamaan tanpa race.
- `TestVerifier_HeaderValidation_And_ErrorBranches` (`:118`) — Verifier header mismatch bad JSON, status mismatch.
- `TestProviderDual_V2Endpoint_DirectAssertion` (`:195`) — Direct `/v2` endpoint test.
- Demo 4 stage (generate → V1 pass → breaking blocked → dual pass) dieksekusi `go run ./cmd/demo` (`cmd/demo/main.go:14-68`).

Available Case Studies: Lab order service: Mobile App sebagai consumer, Order Service sebagai provider (V1 compliant, Breaking, Dual V1+V2) pada routing `/v1/orders/{id}` dan `/v2/orders/{id}`.

Warnings:
- Lab adalah simplifikasi: satu resource order, tidak ada Pact Broker remote, provider state via routing deterministik, tidak ada message broker/Kafka, tidak ada array diff, dan `http.Client` menggunakan timeout 5 * time.Second (sesuai untuk `httptest` dan gate produksi).
- Spring Cloud Contract telah diarsip (Juli 2026) — disebut hanya sebagai konteks historis.
- AsyncAPI hanya diverifikasi sebatas positioning (spec detail tidak diverifikasi).
- Pesan error verifier urutannya tidak deterministik karena iterasi map; test hanya assert jumlah error, bukan urutan.
