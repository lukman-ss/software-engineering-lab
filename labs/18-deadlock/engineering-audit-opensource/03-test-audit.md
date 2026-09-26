# Test Audit

Target Lab: labs/18-deadlock

## Coverage

- Happy path: PASS — TestLockOrderingPreventsDeadlock asserts both succeed, balances 100/100. Demo ordered/retry paths return nil.
- Failure path: PASS — TestDeadlockOccurrence asserts >=1 ErrDeadlock victim under bidirectional naive transfer.
- Edge cases: PARTIAL — 2-account bidirectional covered. No 3+ lock, same-account, zero-amount cases. Out of scope per engineering notes.
- Transitions: PASS — lock acquire -> sleep -> second acquire -> commit/abort verified via errs and balances.
- Recovery: PASS — TestRetryRecoversDeadlock asserts both retry transfers succeed under contention.
- Rollback: NOT_APPLICABLE — balance mutation occurs only after both locks held; no partial state to roll back.
- Concurrency: PASS — all 4 tests use paired goroutines + WaitGroup + opposing directions.
- Negative cases: PASS — victim path asserted, duration-correlation asserted.

## Strength

Suite proves claimed behavior, not just passes. 30x TestDeadlockOccurrence, 30x TestRetryRecoversDeadlock, 20x TestTransactionDurationImpact, 20x race retry — all green. Duration test uses weak `>=` assertion, stable by design, not vacuous: long-delay window deterministically widens circular-wait.

## Execution Record

- `go test ./...` — PASS, 4 passed.
- `go test -race ./...` — PASS, 4 passed.
- `go test -run TestDeadlockOccurrence -count=30` — PASS, 30 passed.
- `go test -run TestRetryRecoversDeadlock -count=30` — PASS, 30 passed.
- `go test -race -run TestRetryRecoversDeadlock -count=20` — PASS, 20 passed.
- `go test -run TestTransactionDurationImpact -count=20` — PASS, 20 passed.
- `go test -v -count=1 ./...` — PASS.
- `go run ./cmd/demo` — PASS, output matches engineering/03-execution-result.md modulo victim line order (nondeterministic, expected).

## Assessment

PASS. No weak-test flag. Missing single-transfer no-contention test and maxRetries-exhaustion test are LOW gaps, core claims fully proven.
