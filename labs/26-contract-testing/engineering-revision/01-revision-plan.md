# Engineering Revision Plan

Target Lab: labs/26-contract-testing
Previous Verdict: APPROVED / APPROVED_WITH_WARNINGS

## Blocking Issues
None.

## Non-Blocking Issues
- GAP-1 (LOW): Response headers declared in contract are not validated by `Verifier.Verify`.
- GAP-2 (LOW): Dual provider `/v2` route has no dedicated assertion in contract tests.
- GAP-3 (LOW): Negative/error branches (status mismatch, invalid JSON body) in `Verifier.Verify` are unproven by direct unit test.
- GAP-4 (LOW): HTTP clients use bare `http.Client{}` without explicit timeout.
- GAP-5 (LOW): Cosmetic doc mismatch regarding contract persistence file path vs in-memory representation.

## Files To Change
- `internal/contract/verifier.go`
- `internal/consumer/client.go`
- `tests/contract_test.go`
- `engineering/01-design.md`
- `engineering/02-implementation-notes.md`

## Tests To Add/Modify
- Add test `TestVerifier_HeaderValidation_And_ErrorBranches` in `tests/contract_test.go`.
- Add test `TestProviderDual_V2Endpoint_DirectAssertion` in `tests/contract_test.go`.

## Validation Commands
```bash
go test -v ./...
go test -race ./...
go run ./cmd/demo
```
