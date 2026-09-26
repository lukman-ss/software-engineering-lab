# Source Map: Backward Compatibility — Expand → Migrate → Contract Pattern

## Definisi dan Konsep Utama

Research:
- research/03-core-concepts.md (Evidence 1: Definisi Backward Compatibility)
- research/06-expand-migrate-contract.md (The Expand → Migrate → Contract Pattern)

Implementation:
- internal/compat/model.go (Struct data dan DTO)
- internal/compat/flags.go (Fitur mode tulis/baca)

Tests:
- internal/compat/service_test.go (TestSerializationBackwardCompatibility)

---

## Expand Phase (Perluasan Skema dan API)

Research:
- research/06-expand-migrate-contract.md (Phase 1: Expand)
- research/05-api-compatibility.md (Additive Fields)

Implementation:
- schema.sql (V2: Expand — CREATE TABLE user_phones)
- internal/compat/store.go (CreateDual, CreateModern)
- internal/compat/service.go (CreateUser dengan mode WriteDual)

Tests:
- tests/migration_test.go (TestFullExpandMigrateContractLifecycle — Step 1: Expand Phase)
- internal/compat/service_test.go (TestSerializationBackwardCompatibility)

---

## Dual-Write Mechanism

Research:
- research/03-core-concepts.md (Evidence 7: Risiko Dual Write)
- research/04-database-migration.md (Dual Write Pattern)

Implementation:
- internal/compat/store.go (CreateDual — atomic write ke kedua tabel)
- internal/compat/service.go (CreateUser dengan mode WriteDual)
- internal/compat/metrics.go (DualWriteCount, DualWriteErrors)

Tests:
- tests/migration_test.go (TestFullExpandMigrateContractLifecycle — Step 2: Verifikasi dual write)
- tests/concurrency_test.go (TestConcurrency)

---

## Migrate Phase (Backfill & Fallback Read)

Research:
- research/06-expand-migrate-contract.md (Phase 2: Migrate)
- research/04-database-migration.md (Backfill Strategies, Fallback Read)

Implementation:
- internal/compat/backfill.go (BackfillWorker — idempotent, resumable batch worker)
- internal/compat/service.go (GetUser dengan mode ReadFallback — lazy backfill)

Tests:
- internal/compat/service_test.go (TestBackfillIdempotentAndResumable, TestFallbackRead)
- tests/migration_test.go (TestFullExpandMigrateContractLifecycle — Step 2 & 3)

---

## Switch Read Path (Read Mode Transition)

Research:
- research/06-expand-migrate-contract.md (Phase 2: Migrate — Consumer Migration)

Implementation:
- internal/compat/flags.go (ReadMode — ReadNewOnly)
- internal/compat/service.go (GetUser dengan mode ReadNewOnly)

Tests:
- tests/migration_test.go (TestFullExpandMigrateContractLifecycle — Step 4: Switch Read Path)

---

## Contract Phase (Drop Legacy Schema & Endpoints)

Research:
- research/06-expand-migrate-contract.md (Phase 3: Contract)
- research/05-api-compatibility.md (API Lifecycle dan Deprecation)

Implementation:
- internal/compat/service.go (ApplyContract — guard dengan metrik traffic)
- internal/compat/store.go (ApplyContractDropLegacyColumn)
- internal/compat/handler.go (GetUserV1 — HTTP 410 Gone setelah contract)
- internal/compat/flags.go (ContractApplied)

Tests:
- internal/compat/service_test.go (TestDeprecationHeadersAndContractEnforcement)
- tests/migration_test.go (TestFullExpandMigrateContractLifecycle — Step 5: Contract)

---

## Rollback Safety

Research:
- research/07-deployment-and-rollback.md (Skenario Rollback, Fundamental Rule)
- research/03-core-concepts.md (Evidence 9: Compatibility Saat Rolling Deployment)

Implementation:
- internal/compat/service.go (CreateUser mode WriteLegacyOnly — menulis ke legacy setelah rollback)
- internal/compat/flags.go (WriteMode — WriteLegacyOnly)

Tests:
- tests/migration_test.go (TestRollbackScenarios — Safe Rollback during Dual-Write, Unsafe Rollback after NewOnly)

---

## Data Drift Reconciliation

Research:
- research/03-core-concepts.md (Evidence 7: Risiko Dual Write)
- research/04-database-migration.md (Dual Write — risiko drift)

Implementation:
- internal/compat/service.go (ReconcileData — membandingkan legacy phone dengan primary number di user_phones)
- internal/compat/metrics.go (DriftDetected)

Tests:
- internal/compat/service_test.go (TestDataReconciliationAndDrift)
- tests/concurrency_test.go (TestConcurrency — goroutine reconciliasi)

---

## API Deprecation Headers (RFC 8594)

Research:
- research/05-api-compatibility.md (Metadata Deprecation — RFC 8594)

Implementation:
- internal/compat/handler.go (GetUserV1 — header Deprecation: true, Sunset)

Tests:
- internal/compat/service_test.go (TestDeprecationHeadersAndContractEnforcement)

---

## Observability and Metrics

Research:
- research/03-core-concepts.md (Evidence 10: Mendeteksi Consumer Lama Masih Aktif)

Implementation:
- internal/compat/metrics.go (Observability — LegacyReadHits, NewReadHits, DualWriteCount, DualWriteErrors, BackfillProcessed, DriftDetected)

Tests:
- internal/compat/service_test.go (semua tes memverifikasi metrik melalui obs.Snapshot())
- tests/concurrency_test.go (TestConcurrency — memverifikasi DualWriteErrors == 0)

---

## Database Schema Evolution (PostgreSQL)

Research:
- research/04-database-migration.md (Zero-Downtime Migration)

Implementation:
- schema.sql (V1: Baseline, V2: Expand, V3: Contract)

Catatan: Perintah PostgreSQL spesifik seperti `CREATE INDEX CONCURRENTLY` dan `ALTER TABLE ... NOT VALID` hanya didokumentasikan dalam research, bukan diimplementasikan dalam kode demo (in-memory store).

---

## Case Studies

Research:
- research/09-case-studies.md
  - Case Study A: Customer Phone (1:1 → 1:N)
  - Case Study B: CMMS Invoice Mechanics (1:1 → N:M)
  - Case Study C: Multi-Currency (Schema Splitting)

---

## Failure Modes

Research:
- research/08-failure-modes.md (Immediate Destructive Schema Alteration, Dual-Write Inconsistency, Missing Backfill Records, Permanent Expand State, Type Change Corruption, Transaction Deadlock)
- research/07-deployment-and-rollback.md (Unsafe Rollback after stopping Dual-Write)

Implementation:
- internal/compat/service.go (ApplyContract guard — mencegah penerapan premature contract)
- tests/migration_test.go (TestRollbackScenarios — Skenario B: rollback setelah WriteNewOnly)

---

## Concurrency Thread Safety

Research:
- research/03-core-concepts.md (Evidence 7: Risiko Dual Write — increased complexity)

Implementation:
- internal/compat/store.go (sync.RWMutex pada semua operasi storage)
- internal/compat/flags.go (atomic.Value, atomic.Bool untuk thread-safe flag)
- internal/compat/backfill.go (sync.Mutex pada checkpoint)
- internal/compat/metrics.go (atomic.Int64 untuk semua counter)

Tests:
- tests/concurrency_test.go (TestConcurrency — 10 writer goroutine, 10 reader goroutine, 1 backfill goroutine, 1 reconciliation goroutine dalam 2 detik)