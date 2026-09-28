# Docs vs Code

## README vs Code
- README commands: `go test ./...`, `go test -race ./...`, `go run ./cmd/demo` — all match actual commands that succeed.
- README structure descriptions match actual package layout. PASS

## Engineering Notes vs Code
- engineering/01-design.md: "100% test coverage on core math and sliding window calculations" — not met (missing edge cases). WARNING
- engineering/02-implementation-notes.md: describes count-based budget correctly. PASS
- engineering/03-execution-result.md: demo output matches actual execution output byte-for-byte. PASS

## Mismatches
1. DOC_CODE_MISMATCH: README says "Multi-Window Multi-Burn-Rate alerting" but each BurnRateRule uses a single BurnRateFactor applied to both short and long windows — no separate short/long burn factors.
2. RESEARCH_IMPLEMENTATION_MISMATCH: Design doc claims 100% test coverage; actual coverage has gaps (see 03-test-audit.md).
3. TEST_CLAIM_MISMATCH: None found — tests match what they claim.

## Overall
Documentation is mostly accurate; two minor mismatches identified.