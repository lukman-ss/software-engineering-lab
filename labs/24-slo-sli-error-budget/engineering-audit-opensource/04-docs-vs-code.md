# Docs vs Code

Target Lab: `labs/24-slo-sli-error-budget`
Compared: `README.md` vs `engineering/01-design.md` vs `engineering/02-implementation-notes.md`
vs `engineering/03-execution-result.md` vs code vs tests vs demo output.

## README.md — PASS

- Structure bullets match actual packages (`internal/metrics`, `internal/slo`, `internal/alerting`,
  `cmd/demo`, `tests/`).
- Commands verified live from the lab dir: `go test ./...`, `go test -race ./...`, `go run ./cmd/demo`.
- "Thread-safety and mathematical correctness" claim is substantiated (race-clean + boundary tests).

## engineering/03-execution-result.md — PASS

- Re-ran all four commands; build exit 0, all 6 tests pass, race clean, demo numbers identical
  (1000/1000/0, 1100/1090/10, 9.09x/9.09x TICKET, Reports -5.00). No invented output.

## Mismatches

### DOC_CODE_MISMATCH-1 — design claims "histogram latency buckets" (MEDIUM)

Location: `engineering/01-design.md` Architecture ("histogram latency buckets & success counts").
Observed: `Bucket` (`tracker.go:15-20`) holds only `TotalCount/GoodCount/BadCount` — no histogram,
no latency distribution. Latency enters only as a boolean via the caller's `isGood` closure.
Notes: README does not repeat this claim; the overclaim is confined to the design doc.

### DOC_CODE_MISMATCH-2 — design claims a "recovery" demo phase (LOW)

Location: `engineering/01-design.md` Architecture ("demonstrating normal traffic, incident budget
burn, alerting, and recovery") and `02-implementation-notes.md` ("What Is Not Demonstrated" omits it).
Observed: `cmd/demo/main.go` has 4 phases ending frozen (both services `CanDeploy=false`); no recovery
or budget-restoration phase exists.
Notes: Honest omission in code; design doc overstates.

### DOC_CODE_MISMATCH-3 — "100% test coverage on core math" unverified (MEDIUM)

Location: `engineering/01-design.md` Success Criteria.
Observed: no coverage report exists; identified uncovered branches (`CalculateBurnRate` zero guards,
100%-error case — see test audit). Claim is unproven, and literal 100% is implausible as stated.
Notes: what IS tested is strong (boundary-exact + negative), but the "100%" quantifier overclaims.

### DOC_CODE_MISMATCH-4 — dead config fields imply unimplemented configurability (MEDIUM)

Locations: `slo.Config.LatencyThreshold` (never read by evaluator); `BurnRateRule.LongWindow`,
`ShortWindow`, `BudgetConsumedPct` (never read by engine).
Observed: docs describe per-endpoint latency policy and per-rule multi-window alerting; code resolves
both at the caller/tracker level instead. Behavior is correct but less configurable than described.

### TEST_CLAIM_MISMATCH-1 — planned "100% errors" test absent (LOW)

Location: `engineering/01-design.md` Test Strategy lists "100% errors" edge case.
Observed: no all-bad test in `tests/slo_test.go`.

## RESEARCH_IMPLEMENTATION_MISMATCH

Out of scope per pipeline override — not assessed. No research-content claims are adjudicated here.
