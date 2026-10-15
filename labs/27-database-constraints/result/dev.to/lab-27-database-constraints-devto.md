---
title: "Database Constraints: The Last Line of Data Integrity: Preventing Race Conditions, Enforcing Schema Invariants, and Rejecting Duplicates Atomically in Go"
published: false
description: "Explore Database Constraints: The Last Line of Data Integrity in Go: mechanics, architecture, caveats, and testable invariants."
tags: go, testing, performance, softwareengineering
---

# Database Constraints: The Last Line of Data Integrity: Preventing Race Conditions, Enforcing Schema Invariants, and Rejecting Duplicates Atomically in Go

A production system must remain resilient under load. Lab 27 makes that behavior visible with a focused Go simulation and deterministic test assertions.

> Application-level checks cannot bridge the race window between read and write operations under concurrency.

## Validasi Aplikasi Runtuh di Bawah Beban Konkuren

The check-then-act pattern inspecting records via SELECT before issuing an INSERT invariably succeeds in sequential testing. Under concurrent traffic, however, two threads observe the same pre-write state before either commits. This failure does not stem from faulty conditional logic, but from lacking atomic isolation at the persistence tier.

## Mengapa Review Manual dan Pengecekan Memori Gagal?

Code reviews routinely miss integrity vulnerabilities because application-level checks appear impeccably structured on screen. Sequential unit tests pass reliably, fostering false confidence across engineering teams. Running validations purely in application services does not automatically ensure integrity without atomic database constraints.

## Model Mental Dua Lapis: Kolaborasi Aplikasi dan Database

Modern software architectures divide integrity responsibilities into two layers: the application handles user experience, while the database guarantees absolute correctness. Application validation filters invalid formats swiftly. Only along this path, when concurrent collisions occur, does the database reject violations atomically with standard SQLSTATE codes.

## Urutan Evaluasi Constraint: Alur Validasi Internal Mesin

The storage engine evaluates schema constraints in a deterministic sequence: NOT NULL runs first, followed by CHECK boolean expressions, UNIQUE index scans, and FOREIGN KEY verification. If any check fails, execution aborts immediately without mutating tables or indexes. This lab models this sequential evaluation atomically within a unified critical section.

## Evaluasi Atomik Constraint di Bawah Mutex Lock

The InsertUser method executes all constraint checks inside an exclusive mutex lock representing PostgreSQL's ROW EXCLUSIVE lock. Locking the critical section does not mean this lab implements low-level B-tree page latches; it cleanly models atomic guarantees so competing goroutines cannot insert conflicting emails.

```go
func (e *Engine) InsertUser(u model.User, usePartial bool) (model.User, error) {
    e.mu.Lock()
    defer e.mu.Unlock()
    if u.Email == "" { return model.User{}, dberr.NewNotNullViolation("users", "email", "users_email_not_null") }
    if u.Age < 18 { return model.User{}, dberr.NewCheckViolation("users", "users_age_check", "age >= 18") }
    if exists := e.emailIndex[u.Email]; exists {
        return model.User{}, dberr.NewUniqueViolation("users", "users_email_key", "duplicate email")
    }
    u.ID = e.userSeq.Add(1); e.users[u.ID] = u; e.emailIndex[u.Email] = u.ID
    return u, nil
}
```

## Kronologi Tabrakan: Dua Goroutine pada Validasi Naif

When two goroutines execute UnsafeStore.RegisterUser concurrently, both observe that the requested email does not yet exist. CPU scheduling delays expose a vulnerability window, prompting both threads to issue INSERT calls. Only along this path do duplicate records infiltrate storage, proving memory checks without constraints are powerless.

## Matriks Taksonomi: Jenis Constraint vs Karakteristik

This matrix charts the five constraint mechanisms verified in our lab alongside standard PostgreSQL Class 23 SQLSTATE codes. Each mechanism addresses a distinct integrity dimension, ranging from single-column presence to cross-table referential integrity. Inspecting SQLSTATE codes enables backends to map low-level database errors into semantic HTTP responses.

Jenis Constraint | SQLSTATE | Cakupan Evaluasi | Pencegahan Kegagalan | HTTP Domain Mapping
|---|---|---|---|---|
NOT NULL | 23502 | Kolom Tunggal | Mencegah record kosong tak beridentitas | 400 Bad Request
CHECK | 23514 | Kolom / Baris | Menolak nilai di luar batas bisnis logis | 422 Unprocessable Entity
UNIQUE | 23505 | Antar Baris (Index) | Mencegah race condition registrasi ganda | 409 Conflict
FOREIGN KEY | 23503 | Lintas Tabel | Mencegah orphan record (induk tidak ada) | 400 Reference Error
PARTIAL UNIQUE | 23505 | Baris Bersyarat | Mengizinkan reuse email pasca soft-delete | 409 Conflict

## Tiga Keputusan Kunci Merancang Integritas Skema

Robust schema design begins by enforcing mandatory fields via NOT NULL and permissible ranges via CHECK. When handling unique entities that support soft deletion, implement partial unique indexes so recycled identifiers remain valid. Adding foreign keys does not mean database optimization is complete; referenced foreign keys require manual indexing to sustain query performance.

## Verifikasi Eksekusi: Test Suite dan Stress Test 50 Goroutine

Execution output confirms that all test suites pass with zero warnings from the Go race detector. In the demo stress test featuring 50 concurrent goroutines submitting an identical email, the storage engine admitted exactly 1 registration while rejecting 49 attempts with error 23505. This demonstrates flawless 100% data integrity under sustained parallel load.

```go
=== RUN   TestNotNullConstraints
--- PASS: TestNotNullConstraints (0.00s)
=== RUN   TestCheckConstraints
--- PASS: TestCheckConstraints (0.00s)
=== RUN   TestUniqueConstraint
--- PASS: TestUniqueConstraint (0.00s)
=== RUN   TestForeignKeyConstraint
--- PASS: TestForeignKeyConstraint (0.00s)
=== RUN   TestPartialUniqueIndex
--- PASS: TestPartialUniqueIndex (0.00s)
=== RUN   TestConcurrentRegistration_Safe_EnforcesUniqueness
--- PASS: TestConcurrentRegistration_Safe_EnforcesUniqueness (0.00s)
=== RUN   TestErrorClassification
--- PASS: TestErrorClassification (0.00s)
PASS: ok  labs/27-database-constraints/internal/store  1.330s (race: clean)
[DEMO] CONCURRENCY STRESS TEST: 50 GOROUTINES FOR SAME EMAIL
Results: Total=50 | Successful=1 | Rejected 23505=49 | Integrity Intact=true
```

## Batasan Bukti Lab: Yang Dibuktikan vs Belum

This lab proves SQLSTATE semantics, atomic locking efficacy, and partial unique index mechanics using a pure Go in-memory simulator. The test suite has not yet proven distributed constraints across multi-node databases or GiST exclusion operators. Maintaining this empirical boundary ensures engineers do not conflate simulator testing with production database benchmarking.

## Kompromi Rekayasa: Constraints vs Fleksibilitas Skema

Enforcing database constraints introduces minor write overhead during INSERTs as the engine inspects B-tree indexes and parent records. Furthermore, attaching constraints to massive production tables demands phased migration workflows to prevent prolonged table locks. However, this overhead does not degrade overall throughput when weighed against the catastrophic cost of data repair.

## Pipeline Validasi Konkuren: Goroutine vs Mutex Lock

This diagram visualizes how the storage engine handles 50 concurrent goroutines racing to register an identical email. The mutex serializes execution at the write boundary, allowing the initial winner to commit its record. All subsequent goroutines are deterministically rejected by the unique index lookup, proving that atomic isolation eliminates race windows entirely.

## Transformasi: Dari UnsafeStore ke SafeStore

Prior to implementing database constraints, UnsafeStore relied on in-memory loops that collapsed under 20 concurrent goroutines, allowing duplicate rows. Migrating to SafeStore delegates integrity to the storage engine, guaranteeing exactly one winner among 50 parallel requests. Retaining application validation upfront does not represent redundancy; it provides an essential bridge for friendly client feedback.

## Verifikasi Test: Membuktikan Kebocoran Data pada UnsafeStore

This unit test empirically demonstrates that application-level validation collapses under concurrent workloads. Twenty goroutines submitting an identical email produce multiple duplicate records in the underlying table. The scope of this failure aligns strictly with lab tests, proving that systems lacking database constraints suffer inevitable data corruption under traffic spikes.

```go
func TestConcurrentRegistration_Unsafe_SuffersRaceCondition(t *testing.T) {
    eng := engine.NewEngine(); s := store.NewUnsafeStore(eng)
    goroutines := 20; var wg sync.WaitGroup; wg.Add(goroutines)
    for i := 0; i < goroutines; i++ {
        go func(idx int) {
            defer wg.Done()
            _, _ = s.RegisterUser(ctx, model.User{Email: "unsafe@example.com", Username: fmt.Sprintf("u_%d", idx)})
        }(i)
    }
    wg.Wait()
    if eng.GetUsersCount() <= 1 { t.Errorf("expected race condition causing >1 user entries, got %d", eng.GetUsersCount()) }
}
```

## Verifikasi Test: SafeStore Menjamin Tepat Satu Pemenang

This test verifies the resilience of SafeStore backed by engine constraints. Among 20 concurrently spawned goroutines, exactly one registration succeeds in persisting to storage, while the remaining 19 receive immediate rejections. Only along this path is table integrity strictly preserved at exactly 1 row, confirming the power of database-level protection.

```go
func TestConcurrentRegistration_Safe_EnforcesUniqueness(t *testing.T) {
    eng := engine.NewEngine(); s := store.NewSafeStore(eng)
    goroutines := 20; var wg sync.WaitGroup; wg.Add(goroutines)
    successCount, errorCount := 0, 0; var mu sync.Mutex
    for i := 0; i < goroutines; i++ {
        go func(idx int) {
            defer wg.Done()
            _, err := s.RegisterUser(ctx, model.User{Email: "safe@example.com", Username: fmt.Sprintf("u_%d", idx), Age: 25, Status: "active"})
            mu.Lock(); defer mu.Unlock()
            if err == nil { successCount++ } else { errorCount++ }
        }(i)
    }
    wg.Wait()
    if successCount != 1 || errorCount != 19 || eng.GetUsersCount() != 1 { t.Fatalf("integrity violated!") }
}
```

## Solusi Soft-Delete: Partial Unique Index WHERE deleted_at IS NULL

Applying a standard UNIQUE constraint to tables with soft deletion creates severe friction: deactivated emails can never be re-registered. By leveraging a conditional Partial Unique Index with WHERE deleted_at IS NULL, the database enforces uniqueness exclusively on active rows. Deactivated records are evicted from the index, allowing identifier recycling without conflict.

```go
// PARTIAL UNIQUE: CREATE UNIQUE INDEX users_active_email_idx ON users (email) WHERE deleted_at IS NULL;
if usePartialUniqueIndex {
    if u.DeletedAt == nil {
        if existingID, exists := e.activeEmails[u.Email]; exists {
            return model.User{}, dberr.NewUniqueViolation("users", "users_active_email_idx", "duplicate active email")
        }
    }
}
// Saat soft delete: delete(e.activeEmails, u.Email) -> email langsung dapat didaftarkan kembali
// Baris dengan deleted_at != nil tidak diindeks, sehingga multi-duplikat soft-delete tetap valid.
```

## Mencegah Data Yatim Piatu dengan Foreign Key Constraint

The InsertOrder function sequentially evaluates three constraints: NOT NULL on UserID, CHECK on TotalCents, and referential FOREIGN KEY validation. If a UserID of 99999 is submitted without existing in the users table, the engine returns SQLSTATE 23503. This mechanism prevents orphaned order records that could compromise financial reconciliation.

```go
func (e *Engine) InsertOrder(o model.Order) (model.Order, error) {
    e.mu.Lock(); defer e.mu.Unlock()
    if o.UserID == 0 { return model.Order{}, dberr.NewNotNullViolation("orders", "user_id", "orders_user_id_not_null") }
    if o.TotalCents <= 0 { return model.Order{}, dberr.NewCheckViolation("orders", "orders_total_cents_check", "total > 0") }
    // FOREIGN KEY: CONSTRAINT fk_orders_user FOREIGN KEY (user_id) REFERENCES users(id)
    if _, exists := e.users[o.UserID]; !exists {
        return model.Order{}, dberr.NewForeignKeyViolation("orders", "fk_orders_user", fmt.Sprintf("key (user_id)=(%d) not present", o.UserID))
    }
    o.ID = e.orderSeq.Add(1); e.orders[o.ID] = o; return o, nil
}
```

## Inspeksi SQLSTATE: Mengubah Error Database Menjadi Domain Error

Many engineers parse raw error strings like duplicate key, which break whenever database locales change. Class 23 SQLSTATE codes provide a standardized, permanent machine contract. Employing MapToDomainError allows applications to translate engine codes into clean domain errors without leaking internal database schemas to clients.

```go
func MapToDomainError(err error) error {
    var cErr *ConstraintError
    if errors.As(err, &cErr) {
        switch cErr.Code {
        case SQLStateUniqueViolation: // 23505
            return fmt.Errorf("conflict: resource with this unique attribute already exists (rule: %s)", cErr.ConstraintName)
        case SQLStateNotNullViolation: // 23502
            return fmt.Errorf("invalid input: mandatory field is missing (rule: %s)", cErr.ConstraintName)
        case SQLStateCheckViolation: // 23514
            return fmt.Errorf("validation failed: value outside boundary (rule: %s)", cErr.ConstraintName)
        case SQLStateForeignKeyViolation: // 23503
            return fmt.Errorf("reference error: referenced entity does not exist (rule: %s)", cErr.ConstraintName)
        }
    }
    return err
}
```

## Checklist Penerapan Database Constraints di Produksi

This checklist summarizes mandatory engineering steps before releasing schemas to production: declaring NOT NULL on mandatory fields, applying UNIQUE constraints to block race conditions, leveraging partial indexes for soft deletes, defining CHECK constraints for row limits, indexing foreign key columns, and inspecting SQLSTATE codes in application error handlers.

## Anti-Pattern Umum dalam Menjaga Integritas Data

The most perilous anti-pattern is assuming safety from application-level SELECT-then-INSERT flows. Other prevalent traps include parsing localized error strings rather than SQLSTATE codes and assuming foreign keys automatically index referencing columns. Overlooking these realities creates ticking performance and data corruption time-bombs.

## Reproduksi Mandiri di Mesin Lokal Menggunakan Go

The entire lab simulation is engineered strictly with the Go standard library, requiring no external packages or Docker containers. While local tests prove constraint semantics deterministically, production environments still require phased migration strategies such as ADD CONSTRAINT NOT VALID.

```sh
cd labs/27-database-constraints
go test -v ./...    # menjalankan 7 unit test untuk NOT NULL, CHECK, UNIQUE, FK, & Partial Index
go test -race ./... # memvalidasi eksekusi konkuren bebas dari data race (PASS in 1.33s)
go run ./cmd/demo   # eksekusi demonstrasi CLI interaktif: SQLSTATE demo & 50-worker stress test
```

## Prinsip Inti: Database Adalah Benteng Terakhir Kebenaran

Application layers will inevitably be refactored and microservices will proliferate, but persisted data outlives software lifecycles. Entrusting data integrity solely to application code is a high-risk architectural gamble. Database constraints guarantee that regardless of application-layer regressions, persistent data remains uncorrupted.

## Source

[Software Engineering Lab 27 - Database Constraints: The Last Line of Data Integrity](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/27-database-constraints), snapshot commit `c317c1d`. Source: README, implementation, tests, and execution records. Code excerpts belong in their original package context.
