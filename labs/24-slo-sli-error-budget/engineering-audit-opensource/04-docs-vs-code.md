# Docs vs Code — labs/24-slo-sli-error-budget

Compared: README.md, engineering/01-design.md, engineering/02-implementation-notes.md, engineering/03-execution-result.md vs code, tests, live demo output.

## Check 1

README structure list vs code: PASS. metrics/slo/alerting/cmd-demo/tests roles match files on disk. Test/demo commands verbatim correct (`go test ./...`, `go test -race ./...`, `go run ./cmd/demo` all reproduced).

## Check 2

Execution-result vs live run: PASS. Re-ran demo; output identical to engineering/03-execution-result.md (Phase1/2/3 numbers, TICKET 9.09x line, DEMO COMPLETE). Test/race transcripts match live runs (4 PASS, no-test-files notices for cmd+internal packages). No FAKE_DEMO, no FAKE_BENCHMARK (no benchmarks claimed).

## Check 3

Design "100% test coverage on core math" vs measured: DOC_CODE_MISMATCH (MEDIUM). CalculateBurnRate 71.4%, NewWindowTracker 66.7% on default/guard branches. Core happy paths are 100%; claim overshoots on edge branches.

## Check 4

Design "histogram latency buckets & success counts" / "ring buffer" vs code: DOC_CODE_MISMATCH (LOW). Tracker is time-bucketed good/bad counters, no latency histogram; storage is append+evict slice, not ring. Behavior as specified; terminology overshoots.

## Check 5

Design "Endpoint Criticality Bucketing (Payment 99.9% vs Reports 95%)" vs code: DOC_CODE_MISMATCH (LOW). `Event.Endpoint` recorded but never used for per-endpoint SLO routing; demo/tests use single SLO. Feature described in design, absent in code; README does not claim it — scoped to design doc only.

## Check 6

Design "demonstration ... and recovery" vs demo: DOC_CODE_MISMATCH (LOW). Demo shows baseline→incident→alert+freeze; no recovery phase (traffic recovering, budget replenishing). Notes correctly list TSDB/PagerDuty as not demonstrated.

No RESEARCH_IMPLEMENTATION_MISMATCH in scope (research excluded per override). No TEST_CLAIM_MISMATCH beyond coverage-percentage wording (tests assert what they claim).
