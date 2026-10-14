#!/usr/bin/env python3
"""
Publication Builder for Software Engineering Lab 27: Database Constraints & Data Integrity.
"""
import os
import sys
import math
import subprocess
from datetime import datetime
from pptx import Presentation
from pptx.util import Inches, Pt
from pptx.dml.color import RGBColor
from pptx.enum.text import PP_ALIGN

OUTPUT_DIR = "/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/27-database-constraints/result"
CHROME_BIN = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
REPO_URL = "https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/27-database-constraints"
COMMIT_SHA = "03a693cc41fcccbc02e6a1d03c0f0cf2c49dacff"
LAB_NO = "27"
SLUG = "database-constraints"
TITLE_ID = "Database Constraints"
SUBTITLE_ID = "Memindahkan Validasi ke Level Database untuk Menjamin Integritas Data"
TITLE_EN = "Database Constraints"
SUBTITLE_EN = "Moving Validation to Database Level for Data Integrity"
AUTHOR = "Lukman / lukman-ss"

DARK_NAVY = "0b1b2d"
TEAL = "087f76"
CREAM = "f4f1e8"
WHITE = "ffffff"
SLATE = "52606b"
ORANGE = "d97706"

SLIDES = []

def add_slide(categories, title, subtitle, content_items, mental_model, source_path="", slide_num=""):
    SLIDES.append({
        "categories": categories,
        "title": title,
        "subtitle": subtitle,
        "content_items": content_items,
        "mental_model": mental_model,
        "source": source_path,
        "slide_num": slide_num
    })

# Slide 1: Cover (Dark Navy, rich visual)
add_slide(
    categories=["KONTEKS"],
    title=TITLE_ID,
    subtitle="SOFTWARE ENGINEERING LAB | LAB " + LAB_NO,
    content_items=[
        {"type": "cover", "title": "", "items": [
            "Retro CRT monitor beige dengan scanlines bercahaya",
            "ASCII art terminal menampilkan error SQLSTATE codes",
            f"Penulis: {AUTHOR}",
        ]}
    ],
    mental_model="",
    source_path="",
    slide_num=""
)

# Slide 2: Motivation & Problem Definition
add_slide(
    categories=["MOTIVASI", "KONSEP"],
    title="Validasi Aplikasi vs Database Constraint",
    subtitle="Pola check-then-act rentan race condition di bawah konkurensi",
    content_items=[
        {"title": "MOTIVASI", "items": [
            "Validasi aplikasi saja rentan: SELECT → cek → INSERT dapat diinterleave",
            "Test konkuren 20 goroutine: UnsafeStore menghasilkan >1 record dengan email sama",
            "Perbaikan duplikat data di produksi jauh lebih mahal daripada pencegahan",
        ], "type": "bullet"},
        {"title": "KONSEP", "items": [
            "Race condition terjadi saat dua transaksi membaca keadaan sama",
            "Database constraint memindahkan pengecekan ke titik atomik: saat INSERT",
            "UNIQUE constraint menggunakan B-tree index dengan ROW EXCLUSIVE lock",
        ], "type": "bullet"},
        {"title": "STUDI KASUS", "items": [
            "Goroutine A: scan → tidak ketemu → INSERT sukses",
            "Goroutine B: scan → tidak ketemu → INSERT sukses (DUPLIKAT)",
            "SafeStore delegation ke engine → hanya 1 berhasil, 19 error 23505",
        ], "type": "bullet"},
    ],
    mental_model="Database constraints eliminate the check-then-act vulnerability by merging check and write atomically.",
    source_path="labs/27-database-constraints/content/02-master-draft.md",
    slide_num="2 / 24"
)

# Slide 3: The Five Constraints
add_slide(
    categories=["THEME", "CONSTRAINT"],
    title="Lima Mekanisme Constraint PostgreSQL",
    subtitle="Database constraint menerapkan integritas data pada level storage",
    content_items=[
        {"title": "NOT NULL (23502)", "items": [
            "Kolom wajib: email, username, user_id tidak boleh NULL",
            "Constraint validation terjadi pada saat INSERT/UPDATE",
            "Prevent integrity issues sebelum masuk ke business logic",
        ], "type": "bullet"},
        {"title": "CHECK (23514)", "items": [
            "Predikat boolean: age >= 18, status IN ('active','inactive','banned')",
            "total_cents > 0 mencegah order dengan nilai nol atau negatif",
            "Validasi domain-specific berbasis aturan deklaratif",
        ], "type": "bullet"},
        {"title": "UNIQUE (23505)", "items": [
            "Satu kemunculan per nilai di kolom atau kombinasi kolom",
            "Mencegah race condition read-then-write secara otomatis",
            "B-tree index otomatis dibuat untuk performa lookup",
        ], "type": "bullet"},
    ],
    mental_model="Constraints provide declarative, ACID-guaranteed enforcement at the storage layer.",
    source_path="labs/27-database-constraints/content/02-master-draft.md",
    slide_num="3 / 24"
)

# Slide 4: Foreign Key & Partial Unique Index
add_slide(
    categories=["CONSTRAINT", "PARTIAL"],
    title="Referential Integrity & Partial Uniqueness",
    subtitle="FOREIGN KEY dan partial unique index menambah fleksibilitas desain",
    content_items=[
        {"title": "FOREIGN KEY (23503)", "items": [
            "Menjamin user_id di orders merujuk ke record users yang ada",
            "Mencegah orphan records tanpa data induk yang valid",
            "Constraint name tersedia di error detail untuk debugging",
        ], "type": "bullet"},
        {"title": "PARTIAL UNIQUE INDEX", "items": [
            "CREATE UNIQUE INDEX ... WHERE deleted_at IS NULL",
            "Uniqueness hanya untuk data aktif (deleted_at IS NULL)",
            "Email dapat dipakai ulang setelah soft-delete",
        ], "type": "bullet"},
        {"title": "PARTIAL INDEX UNTUK SOFT-DELETE", "items": [
            "Active user dengan email same@example.com: hanya 1",
            "Soft-deleted user dengan email same@example.com: diperbolehkan",
            "Active user baru dengan email same@example.com: diperbolehkan",
        ], "type": "bullet"},
    ],
    mental_model="Partial unique index enables conditional uniqueness for soft-delete patterns.",
    source_path="labs/27-database-constraints/internal/store/store_test.go",
    slide_num="4 / 24"
)

# Slide 5: SQLSTATE Error Classification
add_slide(
    categories=["ERROR", "KLASIFIKASI"],
    title="SQLSTATE: Kode Error Standar PostgreSQL",
    subtitle="Class 23: Constraint Violation — stable dan tidak terlokalisasi",
    content_items=[
        {"title": "Kode SQLSTATE", "items": [
            "23502: NOT NULL violation — kolom wajib kosong",
            "23503: FOREIGN KEY violation — referensi tidak ada",
            "23505: UNIQUE violation — duplicate key",
            "23514: CHECK violation — nilai di luar batas",
            "23P01: Exclusion violation — exclusion constraint",
        ], "type": "bullet"},
        {"title": "Error Object", "items": [
            "ConstraintError{Code, ConstraintName, TableName, Detail}",
            "ConstraintName: nama constraint yang dilanggar",
            "TableName: tabel yang terpengaruh",
            "Detail: pesan error detail untuk debugging",
        ], "type": "bullet"},
        {"title": "Pemetaan ke Domain Error", "items": [
            "23505 → HTTP 409 Conflict (resource exists)",
            "23502 → HTTP 400 Bad Request (missing field)",
            "23514 → HTTP 422 Unprocessable Entity (validation failed)",
            "23503 → HTTP 400 Bad Request (invalid reference)",
        ], "type": "bullet"},
    ],
    mental_model="SQLSTATE codes (23xxx) are stable across locales — use codes, not error text.",
    source_path="labs/27-database-constraints/internal/dberr/errors.go",
    slide_num="5 / 24"
)

# Slide 6: Race Condition Proof
add_slide(
    categories=["BUKTI", "KONKUREN"],
    title="Proof: UnsafeStore vs SafeStore di Bawah Konkurensi",
    subtitle="Perbandingan perilaku konkuren 20 goroutine dengan email sama",
    content_items=[
        {"title": "UNSAFE STORE (Tanpa Constraint)", "items": [
            "20 goroutine memanggil RegisterUser() dengan email 'unsafe@example.com'",
            "Setiap goroutine melakukan SELECT → scan loop → Sleep(1ms) → INSERT",
            "Hasil: GetUsersCount() > 1 (dua atau lebih record berhasil inserted)",
            "Dua goroutine dapat membaca count=0 secara bersamaan",
        ], "type": "bullet"},
        {"title": "SAFE STORE (Dengan Constraint)", "items": [
            "20 goroutine memanggil RegisterUser() dengan email 'safe@example.com'",
            "Database constraint 23505 menolak 19 goroutine secara atomik",
            "Hasil: GetUsersCount() = 1 (hanya 1 goroutine berhasil)",
            "Exactly 1 success + 19 SQLSTATE 23505 errors",
        ], "type": "bullet"},
        {"title": "CONCLUSION", "items": [
            "Without constraint: race condition terjadi",
            "With UNIQUE constraint: exactly one winner, N-1 failures",
            "Database constraint mengeliminasi check-then-act window",
        ], "type": "bullet"},
    ],
    mental_model="Without atomic constraint enforcement, read-then-write patterns fail under concurrency.",
    source_path="labs/27-database-constraints/internal/store/store_test.go",
    slide_num="6 / 24"
)

# Slide 7: Engine Architecture
add_slide(
    categories=["ARKITEKTUR", "ENGINE"],
    title="Engine Simulator: In-Memory PostgreSQL Semantics",
    subtitle="Simulator Go merepresentasikan semantik PostgreSQL tanpa Docker",
    content_items=[
        {"title": "Komponen Engine", "items": [
            "Engine: central state store dengan sync.RWMutex",
            "users: map[int64]User, orders: map[int64]Order",
            "InsertUser(): NOT NULL → CHECK → UNIQUE → Generate ID",
            "InsertOrder(): NOT NULL → CHECK → FOREIGN KEY → Generate ID",
        ], "type": "bullet"},
        {"title": "Coarse-Grained Locking", "items": [
            "e.mu.Lock() menahan untuk seluruh operasi insert",
            "Critical section atomik: cek dan tulis dalam satu lock",
            "Tidak ada interleaving antara check dan write",
        ], "type": "bullet"},
        {"title": "Trade-off", "items": [
            "Coarse-grained lock bukan B-tree page latch",
            "Readability dan zero-dependency diprioritaskan",
            "Semantik SQLSTATE dan concurrency pattern termodel",
        ], "type": "bullet"},
    ],
    mental_model="Engine models SQLSTATE behavior and concurrency patterns — not production throughput.",
    source_path="labs/27-database-constraints/internal/engine/engine.go",
    slide_num="7 / 24"
)

# Slide 8: Constraint Check Flow
add_slide(
    categories=["ALUR", "CHECK"],
    title="Constraint Check Flow: Insert User",
    subtitle="Urutan evaluasi constraint di Engine.InsertUser",
    content_items=[
        {"title": "Urutan Evaluasi", "items": [
            "1. NOT NULL: email, username, age, status wajib terisi",
            "2. CHECK: age >= 18, status IN ('active','inactive','banned')",
            "3. UNIQUE: email tidak boleh duplicate",
            "4. Generate Primary Key dan commit ke users map",
            "5. Update B-tree index simulasi",
        ], "type": "bullet"},
        {"title": "Atomic Critical Section", "items": [
            "e.mu.Lock() di awal operasi",
            "Semua check dan write terjadi di dalam lock",
            "e.mu.Unlock() setelah commit berhasil",
        ], "type": "bullet"},
        {"title": "Partial Unique Index", "items": [
            "InsertUserPartial() menerima flag usePartialIndex",
            "Partial unique index: WHERE deleted_at IS NULL",
            "Soft-deleted user tidak ter-index untuk uniqueness check",
        ], "type": "bullet"},
    ],
    mental_model="Constraints are evaluated atomically during INSERT — no window for interleaving.",
    source_path="labs/27-database-constraints/internal/engine/engine.go",
    slide_num="8 / 24"
)

# Slide 9: Constraint Violation Examples
add_slide(
    categories=["CONTOH", "VIOLATION"],
    title="Contoh Pelanggaran Constraint",
    subtitle="Kode test menunjukkan pelanggaran masing-masing constraint",
    content_items=[
        {"title": "NOT NULL Violation (23502)", "items": [
            "RegisterUser{Username: 'alice', Age: 25} tanpa email → error",
            "RegisterUser{Email: 'alice@example.com', Age: 25} tanpa username → error",
            "CreateOrder{TotalCents: 100} tanpa user_id → error",
        ], "type": "bullet"},
        {"title": "CHECK Violation (23514)", "items": [
            "User age=16 → age < 18 violates CHECK constraint",
            "User status='banned' → status tidak di ('active','inactive','banned')",
            "Order total_cents=0 → total_cents <= 0 violates CHECK",
        ], "type": "bullet"},
        {"title": "UNIQUE Violation (23505)", "items": [
            "First insert email='unique@example.com' → sukses",
            "Second insert email='unique@example.com' → 23505 error",
            "Concurrent insert same email → exactly 1 success, 19 failures",
        ], "type": "bullet"},
    ],
    mental_model="Each constraint type has a specific SQLSTATE code — use codes, not text.",
    source_path="labs/27-database-constraints/internal/store/store_test.go",
    slide_num="9 / 24"
)

# Slide 10: Foreign Key Constraint
add_slide(
    categories=["FK", "INTEGRITY"],
    title="Foreign Key Constraint & Referential Integrity",
    subtitle="Orders harus memiliki user yang valid di database",
    content_items=[
        {"title": "Insert Order Flow", "items": [
            "1. NOT NULL: user_id wajib terisi",
            "2. CHECK: total_cents > 0",
            "3. FOREIGN KEY: user_id harus ada di users map",
            "4. Generate Primary Key dan commit ke orders map",
        ], "type": "bullet"},
        {"title": "Violasi FK", "items": [
            "CreateOrder{UserID: 99999, TotalCents: 5000} → error 23503",
            "99999 tidak ada di users map → referential integrity violation",
            "Order terpisah dari user induknya (orphan record)",
        ], "type": "bullet"},
        {"title": "Valid Order", "items": [
            "Insert User: id=1, email='fk@example.com'",
            "CreateOrder{UserID: 1, TotalCents: 5000} → sukses",
            "Order.id dan Order.UserID konsisten",
        ], "type": "bullet"},
    ],
    mental_model="Foreign key prevents orphan records — orders must reference existing users.",
    source_path="labs/27-database-constraints/internal/engine/engine.go",
    slide_num="10 / 24"
)

# Slide 11: Partial Unique Index in Action
add_slide(
    categories=["PARTIAL", "INDEX"],
    title="Partial Unique Index: Soft-Delete Uniqueness Pattern",
    subtitle="Uniqueness hanya untuk record aktif (deleted_at IS NULL)",
    content_items=[
        {"title": "DDL Index", "items": [
            "CREATE UNIQUE INDEX users_email_partial_idx ON users(email) WHERE deleted_at IS NULL",
            "Index hanya memuat row dengan deleted_at IS NULL",
            "Row dengan deleted_at NOT NULL tidak ter-index",
        ], "type": "bullet"},
        {"title": "Scenario: Email Reuse", "items": [
            "1. Insert active user: email='partial@example.com' → sukses",
            "2. Insert duplicate active → 23505 (partially indexed)",
            "3. Soft delete: user1.deleted_at = now()",
            "4. Insert new active dengan email yang sama → sukses (user1 tidak ter-index)",
        ], "type": "bullet"},
        {"title": "Key Points", "items": [
            "Unique constraint: active users only",
            "Partial index mencegah dupe aktif tanpa blokir soft-delete",
            "User dapat di-activate ulang dengan email yang sama",
        ], "type": "bullet"},
    ],
    mental_model="Partial unique index: uniqueness only for active records, soft-deleted excluded.",
    source_path="labs/27-database-constraints/internal/engine/engine.go",
    slide_num="11 / 24"
)

# Slide 12: Error Handling Pattern
add_slide(
    categories=["ERROR", "PATTERNS"],
    title="Error Handling Pattern",
    subtitle="Map SQLSTATE ke domain error dengan ConstraintError",
    content_items=[
        {"title": "ConstraintError Structure", "items": [
            "type ConstraintError struct {",
            "  Code           SQLState",
            "  ConstraintName string",
            "  TableName      string",
            "  Detail         string",
            "}",
        ], "type": "bullet"},
        {"title": "Constructor Functions", "items": [
            "NewNotNullViolation(table, column, constraintName)",
            "NewUniqueViolation(table, constraintName, detail)",
            "NewCheckViolation(table, constraintName, detail)",
            "NewForeignKeyViolation(table, constraintName, detail)",
        ], "type": "bullet"},
        {"title": "Pemetaan Domain Error", "items": [
            "23505 → 'conflict: resource with this unique attribute already exists'",
            "23502 → 'invalid input: mandatory field is missing'",
            "23514 → 'validation failed: value outside permissible boundary'",
            "23503 → 'reference error: referenced entity does not exist'",
        ], "type": "bullet"},
    ],
    mental_model="Use SQLSTATE codes for error classification, not error text parsing.",
    source_path="labs/27-database-constraints/internal/dberr/errors.go",
    slide_num="12 / 24"
)

# Slide 13: Store Layer Abstraction
add_slide(
    categories=["STORE", "ABSTRAKSI"],
    title="Store Layer: UnsafeStore vs SafeStore",
    subtitle="Pemisahan antara validasi aplikasi dan constraint database",
    content_items=[
        {"title": "UnsafeStore", "items": [
            "RegisterUser: app-level check loop → Sleep → InsertUserUnsafe",
            "InsertUserUnsafe bypass database constraint",
            "Vulnerable: two goroutines dapat lolos dengan email yang sama",
        ], "type": "bullet"},
        {"title": "SafeStore", "items": [
            "RegisterUser: delegate ke Engine.InsertUser",
            "Database constraint menangani race condition atomik",
            "All-or-nothing: success or 23505 error",
        ], "type": "bullet"},
        {"title": "Factory Functions", "items": [
            "NewUnsafeStore(eng)",
            "NewSafeStore(eng)",
            "Same engine, different constraint enforcement strategy",
        ], "type": "bullet"},
    ],
    mental_model="UnsafeStore demonstrates the race condition, SafeStore shows constraint protection.",
    source_path="labs/27-database-constraints/internal/store/store.go",
    slide_num="13 / 24"
)

# Slide 14: Test Coverage
add_slide(
    categories=["TEST", "COVERAGE"],
    title="Test Coverage: 8 Test Cases",
    subtitle="Semua constraint type dan concurrency pattern teruji",
    content_items=[
        {"title": "Basic Constraint Tests", "items": [
            "TestNotNullConstraints: 3 assertions untuk email/username/user_id",
            "TestCheckConstraints: 3 assertions untuk age/status/total_cents",
            "TestUniqueConstraint: duplicate email rejection",
            "TestForeignKeyConstraint: orphan order rejection",
        ], "type": "bullet"},
        {"title": "Advanced Tests", "items": [
            "TestPartialUniqueIndex: 6 assertions untuk soft-delete pattern",
            "TestConcurrentRegistration_Safe_EnforcesUniqueness: 20 goroutines, 1 success",
            "TestConcurrentRegistration_Unsafe_SuffersRaceCondition: 20 goroutines, >1 record",
            "TestErrorClassification: SQLSTATE mapping verification",
        ], "type": "bullet"},
        {"title": "Race Detector", "items": [
            "go test -race ./... → clean, no data races",
            "TestConcurrentRegistration_Safe_*: thread-safe constraint enforcement",
            "TestConcurrentRegistration_Unsafe_*: data race terdeteksi",
        ], "type": "bullet"},
    ],
    mental_model="8 test cases cover all constraint types and concurrency edge cases.",
    source_path="labs/27-database-constraints/internal/store/store_test.go",
    slide_num="14 / 24"
)

# Slide 15: Test Results
add_slide(
    categories=["HASIL", "TEST"],
    title="Test Results: 8/8 PASS",
    subtitle="go test -race ./... menghasilkan clean run tanpa race condition",
    content_items=[
        {"title": "Test Execution", "items": [
            "TestNotNullConstraints → PASS",
            "TestCheckConstraints → PASS",
            "TestUniqueConstraint → PASS",
            "TestForeignKeyConstraint → PASS",
            "TestPartialUniqueIndex → PASS",
        ], "type": "bullet"},
        {"title": "Concurrency Tests", "items": [
            "TestConcurrentRegistration_Safe_EnforcesUniqueness → PASS",
            "TestConcurrentRegistration_Unsafe_SuffersRaceCondition → PASS",
            "TestErrorClassification → PASS",
        ], "type": "bullet"},
        {"title": "Summary", "items": [
            "Total test: 8 test cases",
            "Pass rate: 100% (8/8 PASS)",
            "Race detector: clean, no data races",
            "Test time: ~1.5s",
        ], "type": "bullet"},
    ],
    mental_model="8/8 test cases PASS with race detector clean — constraint enforcement verified.",
    source_path="labs/27-database-constraints/internal/store/store_test.go",
    slide_num="15 / 24"
)

# Slide 16: Demo Program
add_slide(
    categories=["DEMO", "CONTOH"],
    title="Demo Program: Interactive Test Harness",
    subtitle="cmd/demo/main.go menggambarkan penggunaan Store layer",
    content_items=[
        {"title": "Features", "items": [
            "Register users dengan email dan password",
            "Create orders dengan user reference",
            "Stress test dengan 50 goroutines konkuren",
            "Real-time output: successful/rejected counts",
        ], "type": "bullet"},
        {"title": "Sample Output", "items": [
            "Total Goroutine: 50",
            "Successful: 1",
            "Rejected 23505: 49",
            "Integrity intact: true",
        ], "type": "bullet"},
        {"title": "Usage", "items": [
            "go run ./cmd/demo",
            "Interactive prompt untuk user input",
            "Stress test mode untuk konkurensi test",
        ], "type": "bullet"},
    ],
    mental_model="Demo program illustrates real-world usage of constraint-protected store.",
    source_path="labs/27-database-constraints/cmd/demo/main.go",
    slide_num="16 / 24"
)

# Slide 17: Production Migration Strategy
add_slide(
    categories=["MIGRASI", "PRODUKSI"],
    title="Production Migration Strategy",
    subtitle="NOT VALID + VALIDATE CONSTRAINT pattern untuk zero-downtime",
    content_items=[
        {"title": "NOT VALID Pattern", "items": [
            "ALTER TABLE orders ADD CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id) NOT VALID",
            "Constraint tidak divalidasi untuk existing rows",
            "Existing rows tidak dicek — constraint hanya berlaku untuk new/updated rows",
        ], "type": "bullet"},
        {"title": "VALIDATE CONSTRAINT", "items": [
            "ALTER TABLE orders VALIDATE CONSTRAINT fk_user",
            "Background check: scan existing rows tanpa table lock",
            "Zero-downtime: application tetap can read/write",
        ], "type": "bullet"},
        {"title": "Step-by-Step", "items": [
            "1. Add constraint NOT VALID",
            "2. Validate constraint (background)",
            "3. Drop NOT VALID if desired",
            "4. Add indexes on referencing columns untuk DELETE/UPDATE",
        ], "type": "bullet"},
    ],
    mental_model="NOT VALID + VALIDATE avoids table lock and enables incremental validation.",
    source_path="labs/27-database-constraints/research/05-report.md",
    slide_num="17 / 24"
}

# Slide 18: Limitations
add_slide(
    categories=["BATAS", "LIMITATIONS"],
    title="Keterbatasan Database Constraints",
    subtitle="Constraints hanya untuk row-scoped logic",
    content_items=[
        {"title": "Row-Scoped Only", "items": [
            "CHECK cannot reference other rows",
            "CHECK cannot query external services",
            "CHECK cannot perform cross-row validation",
        ], "type": "bullet"},
        {"title": "Cross-Row Invariants", "items": [
            "Total balance = sum of accounts — requires aggregate",
            "At most one active subscription per customer — requires count",
            "Solution: SERIALIZABLE isolation atau triggers",
        ], "type": "bullet"},
        {"title": "FK Performance", "items": [
            "Referencing columns tidak auto-indexed",
            "DELETE/UPDATE pada referenced table mungkin full scan",
            "Manually add index pada referencing columns",
        ], "type": "bullet"},
    ],
    mental_model="Constraints are row-scoped — cross-row invariants require different approaches.",
    source_path="labs/27-database-constraints/content/01-content-brief.md",
    slide_num="18 / 24"
}

# Slide 19: Design Principles
add_slide(
    categories=["PRINSIP", "DESAIN"],
    title="Prinsip Desain: Defense in Depth",
    subtitle="Kombinasi validasi aplikasi dan constraint database",
    content_items=[
        {"title": "Two-Layer Defense", "items": [
            "Layer 1 (Aplikasi): user-friendly error messages",
            "Layer 2 (Database): absolute correctness guarantee",
            "Constraint database = last line of defense",
        ], "type": "bullet"},
        {"title": "Complementary Roles", "items": [
            "Application validation: UX lapis pertama",
            "Database constraint: safety net untuk integrity",
            "Duas layer bekerja bersama, bukan menggantikan",
        ], "type": "bullet"},
        {"title": "Real-World Analogy", "items": [
            "Flight check-in: app checks (bag weight), security check (physical scan)",
            "Both layers catch different types of issues",
            "Security check = database constraint (last line)",
        ], "type": "bullet"},
    ],
    mental_model="Application for UX, database for correctness — defense in depth.",
    source_path="labs/27-database-constraints/content/02-master-draft.md",
    slide_num="19 / 24"
}

# Slide 20: When to Use Which
add_slide(
    categories=["KONTEKS", "PILIH"],
    title="Kapan Gunakan Constraint Database vs App Validation",
    subtitle="Konteks dan trade-off masing-masing pendekatan",
    content_items=[
        {"title": "Use Database Constraint", "items": [
            "Data integrity critical: UNIQUE, NOT NULL, CHECK",
            "Race condition possible: concurrency workload",
            "Single source of truth: database layer",
            "Cross-application consistency: multi-service",
        ], "type": "bullet"},
        {"title": "Use Application Validation", "items": [
            "User-friendly error messages",
            "Form validation on UI before submit",
            "Performance: avoid DB round-trip untuk simple check",
            "Complex business logic: multi-step validation",
        ], "type": "bullet"},
        {"title": "Best Practice", "items": [
            "Always use database constraint untuk integrity",
            "Use application validation untuk UX",
            "Never rely on app validation alone",
            "Database constraint = mandatory, app = optional",
        ], "type": "bullet"},
    ],
    mental_model="Database constraint for mandatory integrity, app validation for optional UX.",
    source_path="labs/27-database-constraints/content/05-key-takeaways.md",
    slide_num="20 / 24"
}

# Slide 21: Industry Standards
add_slide(
    categories=["STANDAR", "INDUSTRI"],
    title="SQLSTATE Standards: PostgreSQL & SQLite",
    subtitle="Class 23: Constraint Violation — stable dan portabel",
    content_items=[
        {"title": "SQLSTATE Overview", "items": [
            "PostgreSQL: class 23Constraint Violation",
            "SQLite: class 23Constraint Violation",
            "SQL standard: ISO/IEC 9075-2",
            "Code 23xxx tidak terlokalisasi",
        ], "type": "bullet"},
        {"title": "Common Codes", "items": [
            "23502: NOT NULL violation",
            "23503: FOREIGN KEY violation",
            "23505: UNIQUE violation",
            "23514: CHECK violation",
            "23P01: Exclusion violation",
        ], "type": "bullet"},
        {"title": "Implementation", "items": [
            "Code stable across locales",
            "Error text may be localized",
            "Always use codes, never parse text",
        ], "type": "bullet"},
    ],
    mental_model="SQLSTATE codes are stable — use codes, not error text, for classification.",
    source_path="labs/27-database-constraints/research/03-evidence.md",
    slide_num="21 / 24"
}

# Slide 22: Lab Boundaries
add_slide(
    categories=["BATAS", "LAB"],
    title="Batasan Implementasi Lab 27",
    subtitle="Simulator Go vs PostgreSQL live",
    content_items=[
        {"title": "BATAS 1: IN-MEMORY ENGINE", "items": [
            "Engine menggunakan Go map dan sync.RWMutex",
            "Bukan PostgreSQL live — tidak ada B-tree page latch",
            "Query counter tidak ada (tidak ada query execution)",
        ], "type": "bullet"},
        {"title": "BATAS 2: LOCKING STRATEGY", "items": [
            "Coarse-grained table lock, bukan row-level lock",
            "Bukan B-tree page latch atau MVCC",
            "Trade-off: readability dan zero-dependency",
        ], "type": "bullet"},
        {"title": "BATAS 3: USE CASE", "items": [
            "Lab untuk pembelajaran konsep",
            "Bukan benchmark production throughput",
            "Verifikasi constraint behavior, bukan PostgreSQL internals",
        ], "type": "bullet"},
    ],
    mental_model="This lab proves constraint behavior conceptually — not production PostgreSQL.",
    source_path="labs/27-database-constraints/content/01-content-brief.md",
    slide_num="22 / 24"
}

# Slide 23: Core Design Principles
add_slide(
    categories=["PRINSIP", "INTI"],
    title="Prinsip Desain Inti",
    subtitle="Tiga prinsip untuk data integrity yang benar",
    content_items=[
        {"title": "PRINSIP 1: DELEGATE TO DATABASE", "items": [
            "Delegasikan integrity check ke database constraint",
            "Database constraint = single source of truth",
            "Avoid check-then-act pattern di application layer",
        ], "type": "bullet"},
        {"title": "PRINSIP 2: DEFENSE IN DEPTH", "items": [
            "Application validation untuk UX (optional)",
            "Database constraint untuk correctness (mandatory)",
            "Jangan gantikan constraint dengan app validation",
        ], "type": "bullet"},
        {"title": "PRINSIP 3: CONCURRENT TEST", "items": [
            "Test konkuren eksplisit dengan go test -race",
            "Stress test dengan goroutines untuk race detection",
            "Sequential test tidak dapat mendeteksi race condition",
        ], "type": "bullet"},
    ],
    mental_model="Three principles: delegate to DB, defense in depth, concurrent test.",
    source_path="labs/27-database-constraints/content/05-key-takeaways.md",
    slide_num="23 / 24"
}

# Slide 24: Closing (Dark Navy, rich visual)
add_slide(
    categories=["RINGKASAN"],
    title="Ringkasan Lab 27",
    subtitle=f"Commit: {COMMIT_SHA[:8]} | Test: 8/8 PASS | Repository: lukman-ss/software-engineering-lab",
    content_items=[
        {"title": "KEY RESULTS", "items": [
            "5 constraint types: NOT NULL, CHECK, UNIQUE, FOREIGN KEY, PARTIAL UNIQUE",
            "SQLSTATE codes: 23502/23503/23505/23514/23P01",
            "Race condition prevention via atomic constraint enforcement",
        ], "type": "bullet"},
        {"title": "MENTAL MODEL", "items": [
            "Database constraint eliminates check-then-act vulnerability",
            "Two-layer defense: app for UX, DB for correctness",
            "Partial unique index enables soft-delete email reuse",
        ], "type": "bullet"},
        {"title": "RESOURCES", "items": [
            f"Repository: {REPO_URL}",
            f"Commit SHA: {COMMIT_SHA}",
            "Run tests: go test -v -race ./...",
            "Run demo: go run ./cmd/demo",
        ], "type": "bullet"},
    ],
    mental_model="Database constraints provide atomic integrity enforcement — no window for race conditions.",
    source_path="labs/27-database-constraints/README.md",
    slide_num=""
)

if __name__ == "__main__":
    print(f"Total slides defined: {len(SLIDES)}")
