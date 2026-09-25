# Lab 13: Backward Compatibility — Cara Mengembangkan Sistem Lama Tanpa Merusak Production

## Overview
Implementasi arsitektur **Backward Compatibility** menggunakan pola **Expand -> Migrate -> Contract (Parallel Change)**. Lab ini mendemonstrasikan bagaimana skema database dan contract API bertransisi dari relasi 1:1 (`users.phone`) ke 1:N (`user_phones` table & array payload) tanpa downtime dan tanpa merusak client lawas (V1).

## Project Structure
```text
labs/13-backward-compatibility/
├── cmd/
│   └── demo/
│       └── main.go              # Runnable interactive lifecycle demonstration
├── internal/
│   └── compat/
│       ├── model.go             # Data models and DTOs (V1 legacy, V2 modern, Enriched)
│       ├── store.go             # Thread-safe storage with dual-write and contract drop
│       ├── flags.go             # Feature flags controlling WriteMode and ReadMode
│       ├── metrics.go           # Observability counters for traffic, backfill, and drift
│       ├── backfill.go          # Resumable, idempotent batch backfill worker
│       ├── service.go           # Domain coordinator and data reconciliation
│       ├── handler.go           # HTTP endpoints with Deprecation/Sunset headers
│       └── service_test.go      # Unit tests for core patterns
├── tests/
│   ├── migration_test.go        # Full lifecycle integration tests & rollback scenarios
│   └── concurrency_test.go      # Concurrent reads/writes/backfill under -race detector
├── engineering/
│   ├── 01-design.md             # Engineering design document
│   ├── 02-implementation-notes.md # Architecture notes and limitations
│   └── 03-execution-result.md   # Recorded execution output
├── research/                    # Approved research documentation
├── research-audit/              # Audit verdict and claim review
├── schema.sql                   # SQL schema evolution steps
└── go.mod                       # Go module configuration
```

## Implemented Features
1. **Parallel Change (Expand-Migrate-Contract)**:
   - **Expand**: Model modern (`user_phones`) ditambahkan; payload API diperkaya secara aditif (`phones` array) berdampingan dengan legacy field (`phone`).
   - **Migrate**: Dual-write menulis ke kedua skema; backfill memigrasi data historis dalam batch; fallback read mencegah data starvation.
   - **Contract**: Legacy field/kolom di-drop setelah traffic legacy bernilai 0.
2. **Resumable & Idempotent Backfill**:
   - Memproses data dengan checkpoint `last_processed_id` dan mencegah duplikasi data jika diulang.
3. **Data Drift Reconciliation**:
   - Audit otomatis untuk membandingkan data legacy dan tabel modern sebelum cutover.
4. **Safe Rollback**:
   - Skema dual-write mempertahankan kompatibilitas penuh dengan versi aplikasi N saat terjadi rollback.
5. **Observability & Deprecation**:
   - Counter metrik melacak akses endpoint lawas dan header RFC 8594 (`Deprecation`, `Sunset`).

## Running the Lab

### Run Automated Tests
```bash
go test -v ./...
```

### Run Concurrency Race Detector
```bash
go test -race ./...
```

### Run Demonstration
```bash
go run ./cmd/demo
```
