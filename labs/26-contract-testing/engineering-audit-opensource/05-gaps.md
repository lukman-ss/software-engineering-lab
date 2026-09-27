# Gap Analysis

Target Lab: `labs/26-contract-testing`

Allowed gap types:
- MISSING_TEST
- BROKEN_IMPLEMENTATION
- DOC_CODE_MISMATCH
- RACE_CONDITION
- UNHANDLED_ERROR
- MISSING_EDGE_CASE
- IMPLEMENTATION_OVERCLAIM
- RESEARCH_MISMATCH
- FAKE_DEMO
- FAKE_BENCHMARK
- UNVERIFIED_RESULT

## GAP 1 — MISSING_TEST
Location: No test asserts V2 endpoint `/v2/orders/{id}` is contract-compliant.
Claim: design/01-design.md §4: "V2 contract passes against /v2/orders/{id}"; "introducing V2 DTO while maintaining V1 compatibility."
Test: `TestProviderDual_*` only verifies V1 route on Dual provider; demo Stage 4 prints PASSED but makes no assertion.
Effect: V2 safe-evolution path not covered by automated test; relies solely on demo print.
Severity: MEDIUM
Notes: Could add a test constructing a V2-consumer contract (expecting `currency`, lowercase status, string total) and verify `/v2/orders` against a V2-only contract.

## GAP 2 — MISSING_EDGE_CASE
Location: `tests/contract_test.go:50-72` (`TestProviderBreaking_ContractVerification_Fails`)
Claim: design/01-design.md line 17: "Verification failure with 3 exact breaking change diffs."
Test: Only asserts `len(result.Errors) >= 3`, does not assert which three (enum casing, missing field, type mismatch).
Effect: Regression could yield three different errors yet still pass test.
Severity: LOW
Notes: Strengthen test by asserting presence of each known breaking diff category.

## GAP 3 — MISSING_EDGE_CASE
Location: Contract interaction provider_state field ignored.
Claim: interaction includes `provider_state: "Order ORD-123 exists and is IN_PROGRESS"` (consumer/design).
Implementation: Provider handlers do not honor provider-state; it is a known limitation (implementation-notes.md line 32).
Effect: Demonstrates CDC mechanics but not true provider-state replay; limited to deterministic endpoints.
Severity: LOW
Notes: Documented limitation; acceptable for lab scope but notable as missing edge case.

## GAP 4 — MISSING_EDGE_CASE
Location: No test for non-GET / 404 / unexpected status paths.
Claim: Provider handlers return 405 (method not allowed) and 404 (not found).
Test: No test asserts these negative HTTP paths.
Effect: Minor coverage gap in provider error handling.
Severity: LOW

## Summary

Total gaps identified: 4
- MISSING_TEST: 1 (V2 contract verification)
- MISSING_EDGE_CASE: 3 (which-breaking-diffs, provider-state, negative HTTP paths)

No BROKEN_IMPLEMENTATION, RACE_CONDITION, UNHANDLED_ERROR, IMPLEMENTATION_OVERCLAIM, FAKE_DEMO, FAKE_BENCHMARK, or UNVERIFIED_RESULT found during this implementation+test audit.