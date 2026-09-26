# Content Brief

Topic: Backward Compatibility — Expand → Migrate → Contract (Parallel Change) untuk evolusi skema database 1:1 → 1:N dan kontrak API tanpa downtime
Target Reader: Backend engineer dan DevOps yang memigrasi skema/API produksi dengan klien lama yang belum bisa update serentak
Problem: Perubahan struktur (kolom `users.phone` → tabel `user_phones`, field `phone` → array `phones`) bersifat breaking jika dilakukan langsung: klien V1 crash, migrasi berat menyebabkan downtime, rollback setelah hapus kolom menyebabkan data loss
Core Mental Model: Pecah satu breaking change menjadi tiga fase non-breaking — Expand (dukung keduanya), Migrate (pindahkan data/traffic bertahap), Contract (hapus lama hanya setelah traffic legacy = 0)
Approved Research Status: APPROVED (research-audit/07-verdict.md; audit awal APPROVED_WITH_WARNINGS)
Approved Engineering Status: APPROVED (engineering-audit/06-verdict.md) dan APPROVED_WITH_WARNINGS (engineering-audit-opensource/06-verdict.md, 7 temuan LOW)
Main Concepts:
- Parallel Change / Expand → Migrate → Contract
- Additive schema & payload (kolom/tabel baru + field baru, lama tidak disentuh)
- Dual-write atomic dalam satu critical section (lab, in-memory)
- Backfill batch idempotent + resumable via `last_processed_id`
- Fallback read (dual-read) + lazy backfill
- Feature flags WriteMode/ReadMode untuk dekopling deploy vs aktivasi dan rollback instan
- Observability counters + header Deprecation/Sunset untuk keputusan cutover
Verified Behaviors:
- Payload enriched terbaca V1 (`phone`) dan V2 (`phones`) — TestSerializationBackwardCompatibility
- Backfill batch 3/3/4 dari 10 record, rerun = 0 — TestBackfillIdempotentAndResumable
- Fallback read + lazy backfill — TestFallbackRead
- Drift 1 sebelum backfill, 0 sesudah — TestDataReconciliationAndDrift
- Header Deprecation/Sunset, guard tolak contract saat traffic legacy > 0, 410 Gone pasca-contract, V2 tetap 200 — TestDeprecationHeadersAndContractEnforcement
- Siklus penuh Expand→Migrate→Contract + rollback aman saat dual-write + data loss saat rollback setelah NewOnly — tests/migration_test.go
- Konkuren tulis/baca/backfill/rekonsiliasi lolos -race — tests/concurrency_test.go
Available Case Studies:
- Lab mengimplementasikan: Customer Phone 1:1 → 1:N (Alice/Bob/Charlie). research/09-case-studies.md memuat dua studi tambahan (Invoice 1:1 → N:M, Multi-Currency) yang TIDAK diimplementasikan di lab ini
Warnings:
- Contract irreversible; guard lab memakai cumulative counter tanpa reset/sliding-window, lolos via `force=true` (documented escape, LOW)
- Jendela observasi "30 hari" ilustratif, unverified; GitHub public API memakai ≥24 bulan — bukan rekomendasi universal
- Header `Deprecation` = RFC 9224, `Sunset` = RFC 8594 (README/design doc menyingkat keduanya sebagai RFC 8594 — kosmetik, LOW)
- Store in-memory (`sync.RWMutex`), bukan PostgreSQL; pola DDL Postgres (`CONCURRENTLY`, `NOT VALID`) hanya di research, tidak dieksekusi runtime
- Dual-write menaikkan latensi tulis; skala lintas-microservice butuh outbox/CDC — tidak didemonstrasikan
- LOW lain: handler abaikan error `strconv.Atoi` (404 bukan 400), `GetUserIDs` bubble sort O(n²), `SavePhoneEntry` idempotency dan cabang `WriteNewOnly` tidak di-unit-test langsung
