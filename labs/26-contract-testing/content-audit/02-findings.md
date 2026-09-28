# Content Audit Findings & Inaccuracies

Target Lab: `labs/26-contract-testing`
Audit Date: Mon Sep 28 2026

## Critical Findings

### F-001: False Claim Regarding Response Header Validation (CRITICAL)
- **Location**: `content/02-master-draft.md` (lines 21, 52, 204), `content/04-diagrams.md` (line 158), `content/05-key-takeaways.md` (line 17), `content/07-revision-record.md` (line 8)
- **Claim**: The content repeatedly asserts that "Response header validation is not implemented", "verifier saat ini hanya memvalidasi status code dan body", and references "GAP-01" in open-source engineering audit as a planned enhancement.
- **Actual Implementation**: `internal/contract/verifier.go` (lines 90-98) **fully implements** response header validation:
  ```go
  // Validate response headers if expected by contract
  for expectedHeaderKey, expectedHeaderVal := range interaction.Response.Headers {
      actualHeaderVal := resp.Header.Get(expectedHeaderKey)
      if !strings.EqualFold(actualHeaderVal, expectedHeaderVal) {
          result.Passed = false
          result.Errors = append(result.Errors, fmt.Sprintf("[%s] header mismatch for '%s': expected %q, got %q",
              interaction.Description, expectedHeaderKey, expectedHeaderVal, actualHeaderVal))
      }
  }
  ```
  Furthermore, `tests/contract_test.go` (lines 118-193) contains an explicit test `TestVerifier_HeaderValidation_And_ErrorBranches` validating this behavior.
- **Impact**: Misleads readers into believing the custom verifier lacks header assertion capabilities when it actually supports it.

### F-002: Fabricated Engineering Gap Identifiers (MEDIUM)
- **Location**: Throughout `content/02-master-draft.md`, `content/04-diagrams.md`, `content/05-key-takeaways.md`, `content/07-revision-record.md`
- **Claim**: The content references `GAP-01` (response headers unasserted), `GAP-02` (V2 unverified), and `GAP-06` (map iteration nondeterminism) as official gaps from `engineering-audit-opensource/06-verdict.md`.
- **Actual Implementation/Audit**: The actual open-source engineering audit (`engineering-audit-opensource/05-gaps.md`) lists 3 unnamed non-blocking observations without any "GAP-01/02/06" numbering. While V2 body schema is indeed untested and map iteration is nondeterministic in Go maps, the identifiers and header assertion claim are fabricated/mismatching.
- **Impact**: Process traceability failure; references non-existent formal IDs.

## Minor Findings (LOW Severity)

### F-003: Client Timeout Discrepancy
- **Location**: `content/01-content-brief.md` (line 28), `content/02-master-draft.md` (line 203)
- **Claim**: Claims `http.Client` has no timeout in lab.
- **Actual Implementation**: `internal/consumer/client.go` (line 31) and `internal/contract/verifier.go` (line 55) both configure `Timeout: 5 * time.Second`.
- **Impact**: Minor technical inaccuracy regarding default configuration.

### F-004: Test Suite Count & Scope
- **Location**: `content/01-content-brief.md`, `content/02-master-draft.md`
- **Claim**: Mentions 5 test functions.
- **Actual Implementation**: `tests/contract_test.go` actually contains 7 test functions (including `TestVerifier_HeaderValidation_And_ErrorBranches` and `TestProviderDual_V2Endpoint_DirectAssertion`).
- **Impact**: Incomplete reflection of test coverage.
