# Docs vs Code

README.md
- Claims: internal/metrics sliding-window tracker, internal/slo evaluator & release policy, internal/alerting multi-window burn-rate alerting, demo baseline/incident/alerting.
- Matches: Implementation aligns exactly with README component responsibilities and demo narrative. Verified via source inspection and demo run.
- Mismatches: None observed.

engineering/01-design.md
- Claims: Research APPROVED. Concepts to prove (1) SLI ratio, (2) error budget calc, (3) multi-window burn-rate alerting, (4) endpoint criticality. Expected: accurate SLI, budget depletion triggering alerts & freeze, demo of baseline/incident/alerting/recovery.
- Matches: SLI ratio correct; budget calc correct; multi-window alerting via two trackers + two rules (short/long windows, fast/slow factors); baseline/incident/alerting demonstrated; endpoint criticality shown via dual-SLO comparison.
- Mismatches: 
  - Recovery: design states demo includes recovery; demo contains no recovery phase (no return to CanDeploy=true after incident).
  - 100% coverage claim: design says test coverage on core math and sliding window calculations; actual coverage: evaluator Evaluate 100%, metrics Record 96.9%, alert Check 100%, CalculateBurnRate 71.4%. Some defensive branches (bucketSize<=0, total==0 in burn rate) uncovered.
  - Unused LatencyThreshold: design field LatencyThreshold in SLO.Config not used by Evaluator; demo passes it in but consumes latency via tracker's isGood callback. Clarifies ownership but does not alter outcome.

engineering/02-implementation-notes.md
- Claims: Go 1.22+, stdlib only, time-bucketed ring buffer, bucket resolution vs memory tradeoff.
- Matches: Implementation uses time, sync, testing, math, fmt only; bucketing via truncate; no external deps.
- Mismatches: None.

engineering/03-execution-result.md
- Claims: Recorded build, test, race, demo outputs.
- Matches: Exact verbatim match for build (PASS), tests (all PASS), race (no warnings), and demo output (lines 49-73 line-for-line identical to actual run). No fabrication detected.
- Mismatches: None.

tests/slo_test.go
- Claims: Unit and concurrency tests ensuring thread-safety and mathematical correctness.
- Matches: 6 tests exercising window tracker, evaluator logic, burn rate math, out-of-order, zero traffic, 20-goroutine race. All pass under race detector.
- Mismatches: None.

Summary: Documentation accurately reflects code behaviour. No HIGH/CRITICAL mismatches. LOW-MEDIUM items are recoverable via doc tweaks (recovery claim, coverage overstatement, unused field note).