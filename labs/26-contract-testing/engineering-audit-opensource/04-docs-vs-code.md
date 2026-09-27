# Engineering Audit — Docs vs Code

## Source Set

- README.md (user-facing)
- engineering/01-design.md (design)
- engineering/02-implementation-notes.md (notes)
- engineering/03-execution-result.md (claimed results)
- code under internal/ and cmd/demo
- tests/contract_test.go

## DOC_CODE_MISMATCH

### Mismatch 1 (HIGH): Response header validation

- README line 21: contract/verifier.go = "CDC engine and verification runner".
- design `## Components` item 4: "comparing interactions (method, **path, headers**, status, response schema/types/enums)".
- implementation-notes line 8: "comparing expected contract interactions, **headers**, status codes, and subset JSON schemas".
- implementation-notes line 23: "Validates ... type safety ... enum case ... nested field existence".
- Code verifier.go:30 defines `ResponseDefinition.Headers`; populated in contract (client.go:96-98) and demo. `Verify` (lines 58-115) never reads `Response.Headers`. Only request headers are SET (line 70-72). Response headers are silently ignored.

Verdict: DOC_CODE_MISMATCH — docs claim header validation; code does not validate response headers.

### Mismatch 2 (MEDIUM): V2 endpoint verification

- README line 11: "Preserving V1 contract compatibility while exposing V2 schemas".
- design Expected Behavior line 18: "V1 contract passes against /v1/orders/{id}, while **V2 contract passes against /v2/orders/{id}**".
- design Architecture line 63: `ProviderDual ===> PASS`.
- implementation-notes line 46: "Safe API evolution using **dual DTO routing (/v1 and /v2)**".

Code: ProviderDual handles /v2/orders/ (server.go:112-128). But:
- consumer GenerateMobileContract() only emits ONE interaction for /v1/orders/ORD-123.
- No V2 interaction/contract is generated or verified anywhere.
- Demo Stage 4 verifies only the V1 interaction against dual provider.
- No test touches /v2/orders/{id}.

Verdict: DOC_CODE_MISMATCH — "V2 contract passes" is claimed but no V2 contract exists to pass.

## TEST_CLAIM_MISMATCH

### Mismatch 3 (MEDIUM): "comprehensive test suite" coverage claims

- implementation-notes line 12: "Comprehensive test suite verifying ... breaking change ... dual provider compatibility, and concurrent execution safety."
- design Test Strategy lists 4 specific tests + 1 concurrency test. All 5 exist and PASS.
- BUT design line 18 V2-constraint and line 46 V1-pass-only: test suite does not actually verify V2. Claim of "dual provider compatibility" overstated.

Verdict: TEST_CLAIM_MISMATCH — suite covers V1; not V2.

## RESEARCH_IMPLEMENTATION_MISMATCH

Per pipeline override, research not audited. Skimmed design only. design `## Expected Behavior` and `## Success Criteria` are the binding implementation contract. Mismatches 1 & 2 above are impl-vs-impl (design vs code), recorded as documented gaps.

## Execution-result verification

- 03-execution-result.md:59 claims `go run ./cmd/demo` output includes Stage 3 errors in order [total type mismatch, status value mismatch, customer.name missing].
- Actual run produces order [status value mismatch, customer.name missing, total type mismatch].
- Map iteration in diffValues (verifier.go:131 `for key, expVal := range expMap`) is randomized across runs, so error ordering is nondeterministic.
- Content (all 3 diffs present) matches; ORDER differs.

Verdict: DOC_CODE_MISMATCH (LOW) — ordering in execution-result is a single captured run; nondeterministic. Not a correctness failure but the recorded transcript should be footnoted or the verifier should sort diffs for determinism.

## Summary

| ID | Type | Severity | Subject |
|---|---|---|---|
| 1 | DOC_CODE_MISMATCH | HIGH | response Header validation claimed, not implemented |
| 2 | DOC_CODE_MISMATCH | MEDIUM | V2 contract verification claimed, V2 contract absent |
| 3 | TEST_CLAIM_MISMATCH | MEDIUM | "comprehensive"/"dual" coverage overstated |
| 4 | DOC_CODE_MISMATCH | LOW | demo error ordering nondeterministic vs recorded |
