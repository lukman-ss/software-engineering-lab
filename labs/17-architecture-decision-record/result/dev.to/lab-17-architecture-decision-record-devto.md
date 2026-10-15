---
title: "Architecture Decision Records & Structural Validation: Preventing Cyclical Debates and Decision Lineage Rot in Go"
published: false
description: "Explore Architecture Decision Records & Structural Validation in Go: mechanics, architecture, caveats, and testable invariants."
tags: go, testing, performance, softwareengineering
---

# Architecture Decision Records & Structural Validation: Preventing Cyclical Debates and Decision Lineage Rot in Go

A production system must remain resilient under load. Lab 17 makes that behavior visible with a focused Go simulation and deterministic test assertions.

> Undocumented decisions trigger cyclical technical debates and blind reversal of proven architectural trade-offs.

## Konteks Keputusan Hilang Memicu Debat Teknis Berulang

When historical architectural context is lost, incoming engineers cannot deduce the real trade-offs behind foundational designs. Frequent technical debates do not imply productivity when teams repeatedly re-evaluate resolved trade-offs without written records.

## Wiki dan Markdown Manual Cepat Mengalami Drift

Maintaining design docs in disconnected wikis guarantees documentation drift as code evolves through pull requests. Placing markdown files directly in repositories does not automatically ensure referential integrity unless automated linters enforce bidirectional links.

## Empat Pilar Siklus Hidup Architecture Decision Record

ADRs preserve architectural decisions within version control as immutable artifacts once accepted. Only along this path, through automated validation checks, can decision lineages prevent one-sided status updates and non-monotonic gaps.

## Pipeline Validasi: Dari Raw Markdown Menjadi Error List

The validation pipeline scans raw markdown strings, checks required sections, enforces monotonic numbering, and verifies graph references in parallel. While robust in memory, it has not yet proven automated git pre-commit hook integrations.

## Ekstraksi Metadata dan Seksi Wajib Lewat Regex

The parser extracts IDs and enforces non-empty content across Context, Decision, and Consequences headings. Relying on strict regexes does not imply universal markdown compatibility; it deliberately constrains documents to Nygard's structural standard.

```go
var (
    titleRegex  = regexp.MustCompile(`^#\s+(\d+)\.\s+(.+)$`)
    statusRegex = regexp.MustCompile(`(?i)^Status:\s*([A-Za-z]+)(?:\s+by\s+(\d+))?`)
    contextRegex = regexp.MustCompile(`(?i)^##\s+Context`)
    decisionRegex = regexp.MustCompile(`(?i)^##\s+Decision`)
)
// Parse memverifikasi isi seksi Context, Decision, Consequences
if !foundContext || !contextHasContent { return nil, errors.New("context empty") }
```

## Anatomi Kerusakan Graf: Tautan Sepihak & Referensi Hantu

The most hazardous documentation defect arises from one-sided supersession or referencing non-existent records. Only along this path does the linter intercept structural mismatches, failing PR checks before false architectural states mislead engineers.

## Matriks Komparasi: Gejala Anomali vs Pencegahan Mesin

This matrix contrasts symptoms of manual documentation drift against programmatic linter invariants. Bidirectional link verification and three-state graph coloring eliminate duplicate active states and cyclical supersession loops.

Anomali Dokumentasi | Dampak di Tim Engineering | Aturan Verifikasi Linter
|---|---|---|
Tautan Sepihak | Status ganda: dua arsitektur bentrok | Reciprocal Check (rec.Supersedes == old.ID)
Dangling Pointer | Konteks pengganti hilang tanpa jejak | Exist Check (replacement, exists := recordMap[ID])
Siklus Sirkular | Loop penggantian tanpa ujung (1->2->1) | DFS Cycle Detection (3-state coloring)
Nomor Tidak Urut | Ambiguitas file hilang dalam Git history | Monotonic Numbering Check (ids[i] == i+1)
Format Informal | Status 'Draft' / seksi wajib hilang | Parser Guard (IsValid() & Section Check)

## Evolusi Arsitektur SaaS ERP: Dari Monolith ke Service

The decision sequence reflects pragmatic architectural evolution: adopting a Modular Monolith initially, then extracting notification services at scale. Rejecting Event Sourcing does not indicate fear of innovation; standard relational audit logs satisfied requirements.

## Eksekusi Test Suite dan Demo cmd/demo Lolos Bersih

Terminal logs confirm the complete parser and linter test suite passes cleanly with zero data races across 100 concurrent records. The demo run confirms the ERP evolution scenario satisfies all referential invariants.

```go
=== RUN   TestLinter_ValidSequence
--- PASS: TestLinter_ValidSequence (0.00s)
=== RUN   TestLinter_BrokenReferences (7 sub-tests)
--- PASS: TestLinter_BrokenReferences (0.00s)
=== RUN   TestLinter_ConcurrencyStress (100 records)
--- PASS: TestLinter_ConcurrencyStress (0.00s)
=== RUN   TestParse_Valid & TestParse_Invalid (9 sub-tests)
--- PASS: TestParse_Invalid (0.00s)
PASS: ok lab17/tests 1.469s (race detector: clean)
=== Architecture Decision Record (ADR) Lab ===
Successfully parsed 3 ADRs | Linter: Integrity check passed!
```

## Batasan Bukti Lab: Yang Dibuktikan vs Belum

This lab proves in-memory parsing invariants and bidirectional graph validation algorithms under race detection. It does not prove automated git hook integration or empirical velocity gains for organizations exceeding 100 engineers.

## Trade-off Penerapan ADR dan Validasi Terotomasi

Mandating ADRs introduces review friction by requiring formal documentation for major technical shifts. Running an automated linter does not automatically validate technical wisdom; it enforces structural integrity, not business correctness.

## Validasi Graf Konkuren Menggunakan Goroutines & Mutex

Graph edge validation runs concurrently via goroutines with mutex-protected error aggregation. While this in-memory simulation is not identical to reading thousands of physical markdown files from disk, it passes thread-sanitizer checks cleanly.

```go
var wg sync.WaitGroup
for _, r := range records {
    wg.Add(1)
    go func(rec *Record) {
        defer wg.Done()
        // Evaluasi referensi timbal balik rec.Supersedes & SupersededBy
        if len(localErrs) > 0 {
            mu.Lock()
            errs = append(errs, localErrs...)
            mu.Unlock()
        }
    }(r)
}
wg.Wait()
```

## Transformasi: Dari Catatan Liar ke Graf Tervalidasi

Prior to structured ADRs, technical decisions were lost in chat threads and decaying wiki repositories. Adopting version-controlled ADRs with automated linting does not impose rigid bureaucracy; it locks architectural clarity with minimal syntax.

## Enforcing State Machine Siklus Hidup ADR

The Status enum restricts records to five unambiguous lifecycle states. Enforcing these specific statuses directly reflects test coverage in this lab, preventing ambiguous labels like Draft from leaking into production branches.

```go
type Status string
const (
    StatusProposed   Status = "Proposed"
    StatusAccepted   Status = "Accepted"
    StatusSuperseded Status = "Superseded"
    StatusDeprecated Status = "Deprecated"
    StatusRejected   Status = "Rejected"
)
func (s Status) IsValid() bool { /* switch validasi */ }
```

## Invarian Penomoran Monotonik Mencegah Berkas Hilang

The linter sorts collected IDs and asserts each strictly matches its expected sequential index i+1. Enforcing unbroken sequences does not automatically prevent Git rebase mistakes, but it reliably exposes skipped numbers to reviewers.

```go
ids := make([]int, 0, len(records))
for id := range recordMap { ids = append(ids, id) }
sort.Ints(ids)
for i := 0; i < len(ids); i++ {
    if ids[i] != i+1 {
        errs = append(errs, fmt.Errorf("non-monotonic numbering"))
        break
    }
}
```

## Deteksi Siklus Silsilah dengan Algoritma DFS 3-State

Cycle detection uses a depth-first search graph traversal with three visitation states. Only along this path can recursive checks catch multi-hop loops where ADR 1 is superseded by 2, 2 by 3, and 3 mistakenly points back to 1.

```go
state := make(map[int]int) // 0: unvisited, 1: visiting, 2: visited
detectCycle = func(id int) {
    state[id] = 1
    if rec, ok := recordMap[id]; ok && rec.SupersededBy != 0 {
        nextID := rec.SupersededBy
        if state[nextID] == 1 {
            errs = append(errs, fmt.Errorf("cyclical supersession"))
        } else if state[nextID] == 0 { detectCycle(nextID) }
    }
    state[id] = 2
}
```

## Verifikasi Timbal-Balik: Konsistensi Hubungan Penggantian

The linter guarantees bidirectional agreement between SupersededBy and Supersedes declarations. Enforcing reciprocal metadata does not introduce useless redundancy; it ensures readers at either end of the transition observe consistent status.

```go
if rec.Status == StatusSuperseded {
    replacement, exists := recordMap[rec.SupersededBy]
    if !exists { localErrs = append(localErrs, fmt.Errorf("non-existent")) }
    else if replacement.Supersedes != rec.ID {
        localErrs = append(localErrs, fmt.Errorf("mismatched supersession"))
    }
}
if rec.Supersedes != 0 && (old.Status != StatusSuperseded || old.SupersededBy != rec.ID) {
    localErrs = append(localErrs, fmt.Errorf("not marked as superseded"))
}
```

## Penolakan Berkas Cacat: Parser Guard Seksi Wajib

Unit tests prove the parser rejects records containing malformed titles, unofficial statuses, or empty required sections. This strict validation scope directly matches test coverage, guaranteeing every committed ADR contains actionable rationale.

```go
func TestParse_Invalid(t *testing.T) {
    tests := []struct { name, content, wantErr string }{
        {"missing title", "Status: Accepted", "title not found"},
        {"invalid status", "# 1. T\nStatus: Draft", "invalid status: Draft"},
        {"empty context", "# 1. T\nStatus: Accepted\n## Context\n## Decision\nB...",
         "context section missing or empty"},
    }
    for _, tt := range tests { /* assert strings.Contains(err, wantErr) */ }
}
```

## Checklist Tata Kelola Architecture Decision Record

This checklist codifies essential ADR governance: code co-location, gapless monotonic IDs, closed lifecycle states, immutable accepted records, reciprocal supersession links, and continuous integration linter enforcement.

## Anti-Pattern Umum dalam Dokumentasi Arsitektur

Retroactively editing accepted ADRs erases historical reasoning and invites confusion regarding past trade-offs. Bypassing formal ADRs does not accelerate delivery; it merely compounds cognitive debt for future engineering teams.

## Reproduksi Mandiri di Mesin Lokal Menggunakan Go

The verification suite operates purely on the Go standard library with zero external modules. While this confirms in-memory algorithmic correctness, production systems still require file-crawling CLI tooling and Git pre-commit hooks.

```sh
cd labs/17-architecture-decision-record
go test -v ./...    # validasi parser regex dan aturan linter struktural
go test -race ./... # verifikasi validasi graf konkuren bebas data race
go run ./cmd/demo   # eksekusi demo parsing 3 ADR evolusi SaaS ERP
```

## Prinsip Inti: Keputusan Adalah Kode yang Harus Diuji

Architectural decisions must be governed with the same engineering rigor applied to source code. Storing markdown files without automated testing does not automatically guarantee truth; reliability emerges when structure is enforced via CI pipelines.

## Source

[Software Engineering Lab 17 - Architecture Decision Records & Structural Validation](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/17-architecture-decision-record), snapshot commit `36ff6c8`. Source: README, implementation, tests, and execution records. Code excerpts belong in their original package context.
