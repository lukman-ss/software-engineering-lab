# Source Map

## Master Draft: 02-master-draft.md

### Pendahuluan — Problem Definition
Research:
- `research/01-research-plan.md`
- `research/02-sources.md`
- `research/03-core-concepts.md`

Definition backward compatibility:
- `research/03-core-concepts.md:11-15`
- `research/11-final-research.md:10-15`

### Kompatibilitas Mundar-melit vs Forward Compatibility
Research:
- `research/11-final-research.md:14-15`
- `research/03-core-concepts.md` (implicit definition)

### Identifikasi Breaking Change
Research:
- `research/11-final-research.md:17-23`
- `research/03-core-concepts.md`

### Pola Expand → Migrate → Contract
Research:
- `research/11-final-research.md:25-29`
- `research/04-database-migration.md`
- `research/05-api-compatibility.md`
- `engineering/01-design.md:Introduction`

### Migrasi Skema Tanpa Downtime
Research:
- `research/11-final-research.md:31-35`
- `research/04-database-migration.md:3`

### Dual Read & Fallback Read
Research:
- `research/11-final-research.md:37-38`
- `research/05-api-compatibility.md`

### Risiko Dual Write
Research:
- `research/11-final-research.md:40-43`
- `research/08-failure-modes.md:8-13`

### Backfill Aman
Research:
- `research/11-final-research.md:45-48`

### Rolling Deployment & Kompatibilitas Skema
Research:
- `research/11-final-research.md:50-51`
- `research/07-deployment-and-rollback.md:1-14`

### Pengamatan Legacy Traffic
Research:
- `research/11-final-research.md:54-57`
- `research/11-final-research.md:61` (caveat: 30-day heuristic)

### Penghapusan Field Lama
Research:
- `research/11-final-research.md:59-61`

### Observabilitas untuk Migrasi
Research:
- `research/11-final-research.md:64-67`
- `research/05-api-compatibility.md`

### Feature Flag untuk Rollout/Rollback
Research:
- `research/11-final-research.md:69-70`

### Rollback & Desain Migrasi
Research:
- `research/11-final-research.md:72-75`
- `research/07-deployment-and-rollback.md:25-33`

### Failure Mode Umum
Research:
- `research/08-failure-modes.md` (semua)
- `research/11-final-research.md:77-81`

### Migrasi Resumable & Idempotent
Research:
- `research/11-final-research.md:83-85`

### Pengujian Sebelum Contract
Research:
- `research/11-final-research.md:87-90`

### Perbedaan Database Compatibility vs API Compatibility
Research:
- `research/11-final-research.md:92-94`
- `research/05-api-compatibility.md`

---

## Implementasi: internal/compat/

### model.go — Struktur Data
Implementation:
- `internal/compat/model.go:1-39`

### store.go — Penyimpanan & Dual-Write
Implementation:
- `internal/compat/store.go:1-40` (struktur)
- `internal/compat/store.go:32-51` (CreateLegacy)
- `internal/compat/store.go:53-85` (CreateDual, CreateModern)
- `internal/compat/store.go:120-145` (GetUser, GetPhones)
- `internal/compat/store.go:147-171` (SavePhoneEntry)
- `internal/compat/store.go:173-196` (GetUserIDs)
- `internal/compat/store.go:198-202` (TotalUsers)
- `internal/compat/store.go:204-212` (ApplyContractDropLegacyColumn)
- `internal/compat/store.go:214-218` (IsLegacyDropped)

### flags.go — Fitur Flag
Implementation:
- `internal/compat/flags.go:7-20` (konstanta mode)
- `internal/compat/flags.go:22-31` (struktur FeatureFlags)
- `internal/compat/flags.go:33-62` (konstruktor & metode)

### metrics.go — Observabilitas
Implementation:
- `internal/compat/metrics.go:7-28`

### backfill.go — Worker Backfill
Implementation:
- `internal/compat/backfill.go:9-33` (struktur)
- `internal/compat/backfill.go:35-107` (RunBatch, RunAll, GetProgress)

### service.go — Layanan Domain
Implementation:
- `internal/compat/service.go:9-53` (konstruktor, accessor)
- `internal/compat/service.go:55-95` (CreateUser — routing write mode)
- `internal/compat/service.go:97-148` (GetUser — routing read mode)
- `internal/compat/service.go:150-168` (GetLegacyUser)
- `internal/compat/service.go:170-184` (GetModernUser)
- `internal/compat/service.go:186-226` (ReconcileData)
- `internal/compat/service.go:228-245` (ApplyContract)
- `internal/compat/service.go:247-249` (SerializeResponse)

### handler.go — HTTP Handler
Implementation:
- `internal/compat/handler.go:10-16` (struktur APIHandler)
- `internal/compat/handler.go:18-41` (GetUserV1 — legacy + deprecation headers)
- `internal/compat/handler.go:43-56` (GetUserV2 — modern)

---

## Tests — Verifikasi Perilaku

### Unit Tests (internal/compat/service_test.go)
Tests:
- `TestSerializationBackwardCompatibility:11-45` — V1/V2 unmarshaling
- `TestBackfillIdempotentAndResumable:47-91` — idempotent, resumable
- `TestFallbackRead:93-125` — fallback read + lazy backfill
- `TestDataReconciliationAndDrift:127-154` — drift detection
- `TestDeprecationHeadersAndContractEnforcement:156-205` — header, contract guard

### Integration Tests (tests/migration_test.go)
Tests:
- `TestFullExpandMigrateContractLifecycle:10-94` — seluruh siklus 6 tahap
- `TestRollbackScenarios:96-139` — rollback aman vs tidak aman

### Concurrency Tests (tests/concurrency_test.go)
Tests:
- `TestConcurrency:13-78` — thread safety, race detector

---

## Demo — Verifikasi End-to-End
Command: `go run ./cmd/demo`
Source:
- `cmd/demo/main.go:1-129`

Output Verified:
- Step 1: Baseline (V1 Production State)
- Step 2: Expand Phase (V2 deployed, DualWrite enabled)
- Step 3: Migrate Phase (Backfill + Reconciliation)
- Step 4: Switch Read Path (ReadNewOnly)
- Step 5: Rollback Demonstration
- Step 6: Contract Phase
- Metrics Snapshot

---

## Case Studies
Research:
- `research/09-case-studies.md` — 3 contoh migrasi

---

## Audit Research — Verifikasi Status
Research Audit Verdict:
- `research-audit/07-verdict.md` — APPROVED

Non-Blocking Issues (dihargankan):
- `research/04-database-migration.md` — batch backfill tidak eksplisit di Tier 1 sources
- `research/11-final-research.md:61` — 30-day quiet period heuristic (NOT VERIFIED)

---

## Audit Engineering — Verifikasi Implementasi
Engineering Audit Verdict:
- `engineering-audit/06-verdict.md` — APPROVED

Non-Blocking Issues:
- `engineering-audit/05-gaps.md` — missing HTTP endpoint tests (`net/http/httptest`)
- `engineering/02-implementation-notes.md` — in-memory storage limitation