# Docs vs Code Audit

## Sources compared
- README.md (root)
- engineering/01-design.md
- engineering/02-implementation-notes.md
- engineering/03-execution-result.md
- internal/loadtest/*, internal/server/*, cmd/demo/main.go
- live execution output (auditor run)

## README vs code

| README claim                          | Code match            | Status |
|---------------------------------------|------------------------|--------|
| cmd/demo: Smoke vs Stress entry point | cmd/demo/main.go — exact | PASS |
| internal/server: constrained conn pool| server.go semaphore    | PASS |
| internal/loadtest: runner + percentiles | runner.go + metrics.go | PASS |
| tests: integration tests              | tests/loadtest_test.go | PASS |
| engineering/ dir documented           | present                | PASS |
| Requirements: Go 1.22+                | go.mod go 1.22         | PASS |
| `go test -v ./...` runnable           | yes                    | PASS |
| `go test -race ./...` documented      | yes; passes            | PASS |
| Demo command `go run ./cmd/demo`      | matches                | PASS |

## engineering vs code

Claim: standard library only, no third-party deps — `go.mod` has zero requires. PASS
Claim: Smoke P95 ~= Avg, Stress tail degrades — reproduced live (smoke P95 21.5ms vs Avg 21.3ms; stress P95 772ms vs Avg 463ms). PASS
Claim: 10x pool limit causes queuing — smoke 93 RPS vs stress 89 RPS total (fewer completions due to queue wait). PASS
Claim: All tests pass under `-race` — reproduced. PASS

## Mismatch findings (none blocking)

- DOC_CODE_MISMATCH: README omits `engineering-revision/`, `research/`, `research-audit/`, `research-revision/`, `tests/` already listed. Not a correctness issue; README is high-level. LOW.
- RESEARCH_IMPLEMENTATION_MISMATCH: out of scope per pipeline override; engineering claims align to code regardless. NONE.
- TEST_CLAIM_MISMATCH: `engineering/03` lists test durations (1.36s / 1.657s); auditor run measured 1.63s / 2.240s — environment delta, not claim violation. NONE substantive.
- UNVERIFIED_RESULT: `03-execution-result.md` labels demo numbers "sample run ... actual values vary". README does not reproduce the numbers, only describes. Acceptable.

Verdict: No DOC_CODE_MISMATCH that misleads. Minor omissions in README (no `go vet` mention, no build status badge) — cosmetic.
