# Engineering Audit — Gap Analysis

Gap types per audit spec. Evidence: code read + commands executed (build OK, `go test -v ./...` PASS 5/5, `go test -race -count=1 ./...` PASS, `go vet` clean, `go run ./cmd/demo` OK exit 0).

## GAPS

### GAP-01 — UNVALIDATED_RESPONSE_HEADERS (BROKEN_IMPLEMENTATION / DOC_CODE_MISMATCH)

Severity: HIGH
Location: internal/contract/verifier.go:30, 58-115 (Verify never reads Response.Headers)
Details: `ResponseDefinition.Headers` declared, contract populates it (client.go:96-98), docs claim validation, but code ignores it. A provider returning wrong/missing Content-Type passes.

### GAP-02 — MISSING_TEST for V2 endpoint

Severity: HIGH
Location: tests/contract_test.go (no /v2/orders/ interaction)
Details: Dual provider exposes /v2/orders/ (server.go:112-128) with V2 schema (full_name, string total, currency). No contract, test, or demo interaction exercises it. Design success criterion #5 ("Verification passes for ... dual-versioned providers") is only partially met.

### GAP-03 — IMPLEMENTATION_OVERCLAIM (response header validation)

Severity: HIGH
Location: docs (design 01-design.md:71; implementation-notes 02:8,23)
Details: Docs list response headers among verified dimensions; code does not validate them.

### GAP-04 — RACE_CONDITION risk surface (LOW, not observed)

Severity: LOW
Location: internal/contract/verifier.go:52 (NewVerifier returns singleton Client)
Details: No timeout configured on http.Client; -race passes (stateless per-call), but missing timeout is an operational hazard, not a data race. Included as hardening note.

### GAP-05 — UNDER_ASSERTED breaking-change detection

Severity: LOW
Location: tests/contract_test.go:62 (`len(result.Errors) < 3` lower bound)
Details: Test asserts >=3 errors, not their categories. A regression swapping one error type for another would pass. Weaker than design intent ("exact field diffs").

### GAP-06 — NONDETERMINISTIC_ERROR_ORDERING (DOC_CODE_MISMATCH)

Severity: LOW
Location: internal/contract/verifier.go:131 (map iteration), engineering/03-execution-result.md:104-108
Details: Recorded demo transcript shows errors in an order that differs from a fresh run. Map iteration order is randomized. Correctness intact; transcript should be footnoted or diffs sorted.

### GAP-07 — UNHANDLED_ERROR in demo

Severity: MEDIUM
Location: cmd/demo/main.go:20 (`rawJSON, _ := json.MarshalIndent`)
Details: Demo ignores MarshalIndent error. Low impact (static input), but violates error-propagation rigor claimed in design notes. Not a correctness failure for this lab.

## Gap Summary Table

| ID | Type | Severity | File(s) |
|---|---|---|---|
| GAP-01 | BROKEN_IMPLEMENTATION | HIGH | verifier.go:30,58-115 |
| GAP-02 | MISSING_TEST | HIGH | contract_test.go |
| GAP-03 | IMPLEMENTATION_OVERCLAIM | HIGH | docs (design, impl-notes) |
| GAP-04 | RACE_CONDITION (hardening) | LOW | verifier.go:52 |
| GAP-05 | MISSING_EDGE_CASE (assert strength) | LOW | contract_test.go:62 |
| GAP-06 | DOC_CODE_MISMATCH | LOW | verifier.go:131 + 03-execution-result.md |
| GAP-07 | UNHANDLED_ERROR | MEDIUM | main.go:20 |

## Blocking issues for APPROVED

- GAP-01 (response headers) — HIGH, breaks claimed verification semantics.
- GAP-03 (overclaim) — HIGH, aligns with GAP-01.
- GAP-02 (V2 untested) — HIGH, breaks claimed safe-evolution criterion.

Two of three blocking issues stem from the same root: response-header validation absent. V2 test gap is independent.

## Non-blocking

GAP-04, GAP-05, GAP-06, GAP-07.
