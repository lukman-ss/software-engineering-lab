# Source Map — Cache Invalidation Strategies

Peta ini menghubungkan setiap bagian artikel ke sumber yang digunakan: research, implementation, tests, dan audit.

---

## Overview / Executive Summary

**Research:**
- `research/05-report.md` (Executive Summary, Verified facts list)
- `research/02-sources.md` (Sources 01–11)
- `research/03-evidence.md` (Evidence 1–20)

**Implementation:**
- `internal/cache/store.go` — MemoryCache, TTLWithJitter
- `internal/cache/patterns.go` — three write policies
- `internal/cache/stampede.go` — stampede services, XFetch, SWR

**Tests:**
- `tests/cache_test.go` — `TestCachePatterns`, `TestStampedeMitigation`, `TestXFetchLogic`, `TestStaleWhileRevalidate`, `TestJitter`

**Audit:**
- `research-audit/07-verdict.md` — APPROVED
- `engineering-audit/06-verdict.md` — APPROVED
- `engineering-audit-opensource/06-verdict.md` — APPROVED

---

## Problem & Why This Matters

**Research:**
- `research/05-report.md` (Finding 4: Cache stampede definition)
- `research/03-evidence.md` (Evidence 5: stampede cascading failure; Evidence 18: synthetic numbers disclaimer)

**Implementation:**
- `cmd/demo/main.go` (`demoStampede`) — demonstrasi naive vs singleflight

**Tests:**
- `tests/cache_test.go` (`TestStampedeMitigation`)

---

## Core Concept: Three Write Policies

### Cache-Aside

**Research:**
- `research/05-report.md` (Finding 1)
- `research/03-evidence.md` (Evidence 1, Evidence 2, Evidence 4)
- `research/02-sources.md` (Source 01 — Microsoft Learn)

**Implementation:**
- `internal/cache/patterns.go:22-47` (`CacheAsideService`)

**Tests:**
- `tests/cache_test.go` (`TestCachePatterns/Cache-Aside_Read_&_Write`)

**Audit:**
- `engineering-audit/02-code-audit.md` (Finding 3 — write-then-invalidate ordering)
- `engineering-audit/04-docs-vs-code.md` (MATCH row)

### Write-Through

**Research:**
- `research/05-report.md` (Finding 2)
- `research/03-evidence.md` (Evidence 3, Evidence 4)
- `research/02-sources.md` (Source 01, Source 06)

**Implementation:**
- `internal/cache/patterns.go:50-89` (`WriteThroughService`)

**Tests:**
- `tests/cache_test.go` (`TestCachePatterns/Write-Through_Read_&_Write`)

**Audit:**
- `engineering-audit/02-code-audit.md` (Finding 4 — ReadDelta semantic mismatch)
- `research/04-contradictions.md` (C6 — "same write operation" interpretation)

### Write-Behind

**Research:**
- `research/05-report.md` (Finding 3)
- `research/03-evidence.md` (Evidence 3)
- `research/02-sources.md` (Source 06)

**Implementation:**
- `internal/cache/patterns.go:92-166` (`WriteBehindService`, `flushWorker`, `Update`)

**Tests:**
- `tests/cache_test.go` (`TestCachePatterns/Write-Behind_Asynchronous_Flush`)

**Audit:**
- `engineering-audit/02-code-audit.md` (Finding 5 — drain-on-close race warning; Finding 6 — drop-on-overflow)
- `engineering-audit/05-gaps.md` (GAP-02, GAP-03)

---

## Core Concept: Cache Stampede & Mitigations

### Stampede Definition

**Research:**
- `research/05-report.md` (Finding 4)
- `research/03-evidence.md` (Evidence 5, Evidence 6, Evidence 7)
- `research/04-contradictions.md` (C2 — thundering herd vs stampede terminology)

**Implementation:**
- `internal/cache/stampede.go:16-41` (`NaiveStampedeService`)
- `cmd/demo/main.go` (`demoStampede`)

**Tests:**
- `tests/cache_test.go` (`TestStampedeMitigation/Naive_Stampede`)

### Single-Flight

**Research:**
- `research/05-report.md` (Finding 5)
- `research/03-evidence.md` (Evidence 8, Evidence 9)
- `research/02-sources.md` (Source 03 — Go singleflight pkg)
- `research/04-contradictions.md` (C4 — distributed vs in-process lock)

**Implementation:**
- `internal/cache/stampede.go:44-84` (`SingleFlightService`)
- `go.mod` dependency: `golang.org/x/sync v0.7.0`

**Tests:**
- `tests/cache_test.go` (`TestStampedeMitigation/SingleFlight_Coalesces`)

**Audit:**
- `engineering-audit/02-code-audit.md` (Finding 7 — double-check inside Do)
- `engineering-audit/04-docs-vs-code.md` (MATCH row)

### XFetch (Probabilistic Early Expiration)

**Research:**
- `research/05-report.md` (Finding 6)
- `research/03-evidence.md` (Evidence 10, Evidence 11)
- `research/02-sources.md` (Source 04, Source 08, Source 09)
- `research/04-contradictions.md` (C1 — formula sign)

**Implementation:**
- `internal/cache/stampede.go:87-169` (`XFetchService`, `ShouldRecompute`, `getRand`, `SetRandFunc`)
- `internal/cache/store.go:79-86` (`TTLWithJitter` — related hygiene)

**Tests:**
- `tests/cache_test.go` (`TestXFetchLogic` — pure formula + sign-error check)

**Demo:**
- `cmd/demo/main.go` (`demoXFetch`) — end-to-end with deterministic rand functions

**Audit:**
- `engineering-audit/02-code-audit.md` (Finding 8 — formula correctness; Finding 9 — Get behavior; Finding 11 — mutex on rand)
- `engineering-audit/05-gaps.md` (GAP-01 — integration not directly in test suite)
- `research-audit/06-gaps.md` (Gap 1 — PDF unparsed)

### Stale-While-Revalidate (SWR)

**Research:**
- `research/05-report.md` (Finding 7)
- `research/03-evidence.md` (Evidence 12, Evidence 13)
- `research/02-sources.md` (Source 02 — RFC 5861)
- `research/04-contradictions.md` (C5 — request-triggered vs background job)

**Implementation:**
- `internal/cache/stampede.go:171-254` (`SWRService`, `triggerRevalidate`)

**Tests:**
- `tests/cache_test.go` (`TestStaleWhileRevalidate`)

**Demo:**
- `cmd/demo/main.go` (`demoSWR`)

**Audit:**
- `engineering-audit/02-code-audit.md` (Finding 10 — dedup guard)
- `engineering-audit/04-docs-vs-code.md` (MATCH row)

### Jitter

**Research:**
- `research/05-report.md` (Finding 8)
- `research/03-evidence.md` (Evidence 16)
- `research/04-contradictions.md` (C2 — terminology)

**Implementation:**
- `internal/cache/store.go:79-86` (`TTLWithJitter`)

**Tests:**
- `tests/cache_test.go` (`TestJitter`)

**Demo:**
- `cmd/demo/main.go` (`demoJitter`)

**Audit:**
- `engineering-audit/02-code-audit.md` (Finding 1 — range [base, base+maxJitter))

---

## Architecture

**Research:** none primary (local design doc)
**Implementation:**
- `engineering/01-design.md` (Architecture section)
- `engineering/02-implementation-notes.md` (Core Design Decisions, Implementation-Specific Choices, Known Limitations, Trade-offs, What Is Demonstrated, What Is Not Demonstrated)
- `internal/cache/store.go`, `repo.go`, `patterns.go`, `stampede.go`

**Tests:**
- `tests/cache_test.go`

---

## Implementation Walkthrough / Code

**Primary source files (see code snippets):**
- `internal/cache/patterns.go` — Snippets 1, 2, 3
- `internal/cache/stampede.go` — Snippets 4, 5, 6
- `internal/cache/store.go` — Snippet 7

**Supporting:**
- `internal/cache/repo.go` — MockDB helpers
- `cmd/demo/main.go` — runnable demo for all scenarios

---

## What the Tests Prove

**Implementation:**
- `tests/cache_test.go` — full test matrix

**Audit:**
- `engineering-audit/03-test-audit.md` — Coverage by Feature table
- `engineering-audit/04-docs-vs-code.md` — NO TEST_CLAIM_MISMATCH

**Execution Result:**
- `engineering/03-execution-result.md` (test outputs, race detector, demo output)

---

## Failure Scenarios

**Research:**
- `research/05-report.md` (Limitations section)
- `research/03-evidence.md` (Evidence 18 — synthetic numbers disclaimer)

**Implementation:**
- `engineering/02-implementation-notes.md` (Known Limitations)
- `engineering-audit/05-gaps.md` (GAP-02 — drop on overflow)

**Tests:**
- `tests/cache_test.go` (`TestStampedeMitigation/Naive_Stampede` demonstrates stampede condition)

---

## Production Considerations

**Research:**
- `research/05-report.md` (Conclusion — single-flight scope, XFetch beta, jitter limitations)

**Implementation:**
- `engineering/02-implementation-notes.md` (Known Limitations, Trade-offs table)
- `engineering-audit/02-code-audit.md` (Finding 5 — drain-on-close warning)

---

## Checklist / Key Takeaways

**Research:**
- `research/05-report.md` (Conclusion)

**Implementation:**
- `engineering/03-execution-result.md` (Demo output)

---

## Sources — External References

**Tier 1 (primary/vendor):**
- `research/02-sources.md` Source 01: Microsoft Learn — Cache-Aside Pattern
- `research/02-sources.md` Source 02: IETF RFC 5861
- `research/02-sources.md` Source 03: Go `golang.org/x/sync/singleflight`
- `research/02-sources.md` Source 08: Vattani et al. PVLDB 2015 (DOI verified)

**Tier 3 (synthesis/cross-check):**
- `research/02-sources.md` Source 04: Wikipedia — Cache stampede
- `research/02-sources.md` Source 05: Wikipedia — Thundering herd problem
- `research/02-sources.md` Source 06: Wikipedia — Cache (computing)
- `research/02-sources.md` Source 07: Wikipedia — Cache invalidation

**Local references (pedagogical context only):**
- `research/02-sources.md` Source 11: prior lab `labs/04-caching/*`

**Audit records:**
- `research-audit/*` — audit plan, source audit, claim audit, contradictions, gaps, verdict
- `engineering-audit/*` — audit plan, code audit, test audit, docs-vs-code, gaps, verdict
- `engineering-revision/03-revision-result.md` — no revisions needed
