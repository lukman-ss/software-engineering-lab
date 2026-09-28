# Gap Analysis: Contract Testing (lab-26)

## Gap 1 — Response headers declared but not enforced

- **Type:** MISSING_TEST
- **Severity:** LOW
- **Location:** internal/contract/verifier.go:Verify (status/body checked; headers ignored)
- **Description:** Contract's expected `Content-Type: application/json` header not asserted. Provider returning wrong content-type with valid JSON body would pass verification.
- **Recommendation:** non-blocking; either enforce header subset check or document as intentional (header-agnostic) design.

## Gap 2 — No request timeout / cancellation context

- **Type:** UNHANDLED_ERROR
- **Severity:** LOW
- **Location:** internal/contract/verifier.go:74 (v.Client.Do), internal/consumer/client.go:36
- **Description:** Default `http.Client` has no timeout; hung provider would block CI gate indefinitely instead of failing fast.
- **Recommendation:** non-blocking; lab-scoped simplification, acceptable for httptest demo. Production gate should set `Timeout` / context deadline.

## Gap 3 — V2 endpoint behavior asserted only indirectly

- **Type:** MISSING_TEST
- **Severity:** LOW
- **Location:** tests/contract_test.go (Dual test verifies V1 path only)
- **Description:** Dual provider's `/v2/orders/{id}` route has no dedicated test asserting its evolved schema (status lowercase, total string, currency) is served correctly.
- **Recommendation:** non-blocking; current tests prove the lab's core claim (V1 compatibility preserved). V2-shape test would strengthen evolution narrative.

## Gap 4 — Source-map line drift (demo orchestrator)

- **Type:** DOC_CODE_MISMATCH
- **Severity:** LOW
- **Location:** content/06-source-map.md:60 (`cmd/demo/main.go:14-68` vs actual 14-71)
- **Description:** Line range stale by 3 lines after final print statement added.
- **Recommendation:** non-blocking cosmetic fix.

## Excluded findings

No MISSING_EDGE_CASE, RACE_CONDITION, BROKEN_IMPLEMENTATION, IMPLEMENTATION_OVERCLAIM, RESEARCH_MISMATCH, FAKE_DEMO, FAKE_BENCHMARK, or UNVERIFIED_RESULT. All core claims executed and reproduced live.