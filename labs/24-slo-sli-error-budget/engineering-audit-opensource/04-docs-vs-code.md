# Docs vs Code Audit

## Finding 1: Time window compression claim mismatch

Location:
- Document: `engineering/01-design.md:31` — "In-memory ring/time-bucketed window tracking to simulate 30-day/rolling windows in compressed real-time (e.g. 1-second = 1-hour scale for demo)"
- Document: `engineering/02-implementation-notes.md:21` — Same claim about compressed real-time windows
- Code: `cmd/demo/main.go:24-30` — Uses literal `30 * time.Minute`, `5 * time.Minute`, `60 * time.Minute` windows
- README: No mention of compression
Claimed: 30-day windows simulated using time compression (1s = 1h scale)
Observed: No time compression; literal 30-minute window used with real-time timestamps
Type: DOC_CODE_MISMATCH
Assessment: WARNING
Severity: MEDIUM
Notes: Neither README nor code implements the claimed time compression. The engineering design document's stated approach does not match the actual implementation. This could mislead readers about the system's ability to demonstrate long-term (30-day) behavior in a compressed format.

## Finding 2: README structure description matches code

Location:
- README: `README.md:7-11` — Describes four components: metrics, slo, alerting, cmd/demo, tests
- Code: Actual file structure matches: `internal/metrics/`, `internal/slo/`, `internal/alerting/`, `cmd/demo/`, `tests/`
Claimed: Structure matches
Observed: All described directories and their stated purposes align with code
Type: None
Assessment: PASS
Severity: N/A
Notes: README accurately describes the repository structure.

## Finding 3: "100% test coverage" success criterion not substantiated

Location:
- Document: `engineering/01-design.md:21` — "100% test coverage on core math and sliding window calculations" (Success Criteria)
- Document: `engineering/03-execution-result.md:11-28` — Test results shown, but no coverage percentage reported
Claimed: 100% test coverage achieved
Observed: No coverage metrics in execution result; test output shows 6 tests passing but does not quantify coverage
Type: IMPLEMENTATION_OVERCLAIM
Assessment: WARNING
Severity: MEDIUM
Notes: The success criterion claims "100% test coverage" but the execution result does not verify or report coverage metrics. Without actual `go test -cover` output, this claim is unverifiable. Manual analysis shows gaps in edge case coverage (see test audit Finding 7).

## Finding 4: README test/run commands match actual commands

Location:
- README: `README.md:15-18` — `go test ./...` and `go test -race ./...`
- README: `README.md:22-24` — `go run ./cmd/demo`
- Actual: Commands run successfully and produce expected output
Claimed: Commands are documented correctly
Observed: Commands match and execute successfully
Type: None
Assessment: PASS
Severity: N/A
Notes: README instructions are accurate and executable.

## Finding 5: Execution result typo in Phase 1 description

Location:
- Document: `engineering/03-execution-result.md:49` — "Simulating Baseline Traffic (1,000 requests, 1,000 requests, 100% success)..."
- Actual: `cmd/demo/main.go:55` — "Simulating Baseline Traffic (1,000 requests, 100% success)..."
Claimed: Phase 1 simulates 1000 requests with 100% success rate
Observed: The recorded output contains a duplicated "1,000 requests" in the description, while the actual demo output does not
Type: DOC_CODE_MISMATCH
Assessment: WARNING
Severity: LOW
Notes: Minor typo in recorded execution results. When the demo was actually re-run during this audit, the output read correctly as "(1,000 requests, 100% success)...". This suggests the execution result was either manually transcribed or from a different version.

## Finding 6: Demo output matches execution result (with minor discrepancy)

Location:
- Document: `engineering/03-execution-result.md:44-65` — Recorded demo output (Phases 1-4)
- Actual: Demo output verified during audit execution on 2026-09-26
Claimed: Demo produces exact output shown in execution results
Observed: Phase 3 output in execution result shows only 1 alert (TICKET/Slow Burn), while actual demo execution produces identical Phase 3 output. Phase 4 output in execution result is NOT shown (the result file ends after Phase 3), but actual demo does produce Phase 4 output.
Type: TEST_CLAIM_MISMATCH (incomplete recording)
Assessment: WARNING
Severity: LOW
Notes: The execution result omits Phase 4 output entirely, yet the demo code includes it and it executes correctly. The recorded execution result is incomplete compared to actual runtime.

## Finding 7: Research claims align with implementation intent

Location:
- Research: `research/05-report.md` — Findings 1-11 define SLI, SLO, error budget, burn rate concepts
- Implementation: Code implements good/total ratio for SLI, (1-SLO)*total for budget, burn-rate alerting
- Engineering: `engineering/01-design.md` — "Concept To Prove" maps concepts 1-4 to implementation
Claimed: Implementation follows research definitions
Observed: Core concepts are correctly translated:
- Finding 1 (SLI as ratio) → Evaluator computes good/total ratio
- Finding 5 (Error budget = 100% - SLO) → Evaluator uses (1 - target) for budget
- Finding 8 (Burn rate as decision tool) → AlertEngine calculates burn rate
- Claim 4 (Endpoint criticality) → Demo compares 99.9% vs 95% SLO
Type: None
Assessment: PASS
Severity: N/A
Notes: Implementation faithfully follows research-backed definitions for core SLO/SLI/error budget concepts. The Datadog-specific burn rate threshold formula mentioned in Finding 12 is NOT implemented; instead, Google SRE-style multi-window burn rate alerting is used.

## Finding 8: Implementation does not use Datadog's "error budget remaining" formula

Location:
- Research: `research/05-report.md:158-167` (Finding 12) — Claims error budget remaining = `100 * (current_status - target) / (100 - target)`
- Implementation: `internal/slo/evaluator.go:50-52` — Computes `budgetRemaining = totalErrorBudget - budgetConsumed` = `(1 - target) * total - bad`
Claimed: Error budget remaining formula from Datadog (Finding 12) would be implemented
Observed: The implementation uses a different formula entirely: remaining = (allowed failure count) - (actual bad count). The Datadog formula computes remaining as a percentage of the total budget.
Type: RESEARCH_IMPLEMENTATION_MISMATCH
Assessment: WARNING
Severity: MEDIUM
Notes: The research Finding 12 explicitly documents the Datadog formula, but the implementation does not use it. Instead, it uses a count-based budget. This is an implementation-specific design decision that should be documented. The research noted this formula was "Datadog-specific implementation" (Finding 12, Confidence: MEDIUM), so the deviation may be intentional, but it is not documented as a deviation in the engineering notes.

## Summary

Three key mismatches identified:
1. **High**: Time compression claim in design docs vs literal short windows in code
2. **Medium**: Datadog error budget formula documented in research vs different formula in implementation
3. **Medium**: "100% test coverage" claim unsubstantiated by coverage metrics

The core implementation correctly follows Google SRE principles as documented in research, and all executable commands work as described in README.