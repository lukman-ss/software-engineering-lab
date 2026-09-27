# Docs vs Code

Target Lab: labs/26-contract-testing

## Comparison

| Source | Claim | Matches Code? | Status |
|--------|-------|---------------|--------|
| README.md | "Consumer declares minimal required schema/interactions" | consumer.GenerateMobileContract() generates exactly {id, status, customer.name, total} | PASS |
| README.md | "Provider verifies implementation against consumer contracts before deployment" | Verifier.Verify() run against ProviderV1/Breaking/Dual | PASS |
| README.md | "Detects enum casing changes, field renames, primitive type mutations" | Verifier detects status IN_PROGRESS vs in_progress, missing customer.name, int vs string total (3 errors) | PASS |
| README.md | "Preserving V1 contract compatibility while exposing V2 schemas" | ProviderDual serves /v1/orders/ (V1-compliant) and /v2/orders/ (V2 schema) | PASS |
| README.md | Commands: go test ./..., go test -race ./..., go run ./cmd/demo | All verified executable | PASS |
| engineering/01-design.md | "Full test suite passing with race detector" | go test -race PASS | PASS |
| engineering/01-design.md | "Verification passes for valid providers and dual-versioned providers" | Tests pass | PASS |
| engineering/03-execution-result.md | Reports 3 breaking errors | Actual demo produces exactly 3 errors | PASS |

## Findings

- No DOC_CODE_MISMATCH observed. README and demo output align.
- No TEST_CLAIM_MISMATCH: all claimed tests exist and pass.
- No RESEARCH_IMPLEMENTATION_MISMATCH: design claims match implementation.
- NOTE: ProviderState field in contract struct is documented as part of interaction but unused in implementation (not a mismatch, design-level simplification).