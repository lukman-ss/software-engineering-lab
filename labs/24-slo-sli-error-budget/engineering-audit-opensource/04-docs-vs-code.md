# Docs vs Code

## README vs Code
- README: "internal/slo: Evaluator calculating SLI ratios, remaining Error Budget, and release freeze policy enforcement." Matches code.
- README: "Multi-window burn rate alert calculator." Matches code.
- README does not mention LatencyThreshold, matching actual usage (dead field).

## Engineering Design vs Code
- Design: "In-memory ring/time-bucketed window tracking." Matches implementation.
- Design: "standard library only." Confirmed (time, sync, math, fmt).

## Execution Result vs Actual Run
- Recorded demo output matches actual `go run ./cmd/demo` output exactly.
- Recorded test output matches actual `go test` output exactly.

## Mismatch Summary
No DOC_CODE_MISMATCH in external behavior. Dead field Config.LatencyThreshold not documented as functional.

Result: PASS.