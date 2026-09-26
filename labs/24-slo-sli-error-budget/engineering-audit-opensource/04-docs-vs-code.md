# Documentation vs Code

## Comparison Matrix

| Artifact | Reference | Status |
|----------|-----------|--------|
| README running instructions | go test ./... ; go test -race ./... ; go run ./cmd/demo | Matches code exactly. PASS |
| README `tests/` description (unit + concurrency tests) | tests/slo_test.go | Matches. PASS |
| README `cmd/demo` description (baseline, incident, alert) | cmd/demo/main.go phases 1-3 | Matches. PASS |
| engineering/02-implementation-notes.md "Files Added" | actual files present | Matches. PASS |
| engineering/02-implementation-notes.md "Standard library only" | go.mod has no requires | Matches. PASS |
| engineering/01-design.md "100% test coverage" success criterion | go test -coverpkg => 94.1% | MISMATCH — claim not met |
| engineering/01-design.md WindowTracker described as "ring buffer" | tracker.go uses grow/shrink slice | Terminology mismatch |
| engineering/03-execution-result.md demo output | cmd/demo/main.go | MISMATCH — recorded output omits Phase 4 |

## Finding D1: DOC_CODE_MISMATCH — Stale demo output in execution result

Location: engineering/03-execution-result.md (lines 53-65) vs cmd/demo/main.go (lines 117-150)
Claimed Behavior: Recorded demo run ends after [PHASE 3] alert, then "DEMO COMPLETE".
Observed Implementation: Current cmd/demo/main.go contains an additional [PHASE 4] Endpoint Criticality Comparison that runs between the alert and DEMO COMPLETE.
Recorded Output (my run):
```
[PHASE 4] Endpoint Criticality Comparison (Payment 99.9% vs Reports 95.0%)...
Reports Target SLO: 95.0% | Current SLI: 90.0% | Budget Remaining: -5.00
Payment CanDeploy: false | Reports CanDeploy: false (Reports has wider 5% error tolerance)
```
Assessment: FAIL (DOC_CODE_MISMATCH)
Severity: MEDIUM
Notes: The recorded execution artifact does not match the current code. Either the code was extended without updating the record, or the record is from an earlier revision. The real demo runs fine; only the saved artifact is stale.

## Finding D2: RESEARCH_IMPLEMENTATION_MISMATCH — False 100% coverage claim

Location: engineering/01-design.md (line 21) vs go test -coverpkg measurement (94.1%)
Claimed Behavior: "100% test coverage on core math and sliding window calculations" as a success criterion.
Observed Implementation: Actual coverage is 94.1%, with uncovered branches in CalculateBurnRate and NewWindowTracker.
Assessment: FAIL (claim overstated)
Severity: MEDIUM
Notes: The success criterion was not verified; no cover profile was produced to substantiate the 100% figure.

## Finding D3: IMPLEMENTATION_OVERCLAIM — "Recovery" demonstrated

Location: engineering/02-implementation-notes.md (line 29) + 01-design.md (line 29)
Claimed Behavior: Demo demonstrates "...alerting, and recovery."
Observed Implementation: cmd/demo has no recovery phase. Phase 2 depletes budget; Phase 3 alerts; Phase 4 compares criticality. Budget never recovers; no test covers recovery.
Assessment: WARNING (IMPLEMENTATION_OVERCLAIM)
Severity: LOW
Notes: "Recovery" listed as a design concept but neither demoed nor tested.

## Finding D4: Misleading in-code note on criticality comparison

Location: cmd/demo/main.go (line 145)
Claimed Behavior: "(Reports has wider 5%% error tolerance)" implies Reports should remain deployable where Payment is frozen.
Observed Implementation: Reports (95% SLO, 10% error rate) is also CanDeploy=false (budget -5.00), so the parenthetical contradicts actual output.
Assessment: WARNING (DOC_CODE_MISMATCH in note)
Severity: LOW
Notes: Logic is correct; the human-readable note is misleading. The 10% error rate exceeds even the lenient 5% Reports tolerance.

## Finding D5: Terminology drift — "ring buffer"

Location: engineering/02-implementation-notes.md (line 21) vs tracker.go
Claimed Behavior: "In-memory ring/time-bucketed ring buffer."
Observed Implementation: A Go slice is appended to and evicted by re-slicing. No fixed-size ring.
Assessment: WARNING (terminology mismatch)
Severity: LOW
Notes: Functionally equivalent sliding window; the "ring buffer" label is inaccurate but harmless.

## Summary
README aligns with code (PASS). Engineering notes have one stale demo artifact (MEDIUM), one false coverage claim (MEDIUM), and minor terminology/over-claim issues (LOW).