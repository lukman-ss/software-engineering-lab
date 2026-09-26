# Docs vs Code

Target Lab: labs/18-deadlock

## Sources Compared

- README.md
- engineering/01-design.md, 02-implementation-notes.md, 03-execution-result.md
- internal/bank/account.go, internal/transfer/transfer.go, cmd/demo/main.go
- tests/transfer_test.go

## Findings

1. README file map — MATCH. Paths and function names (`TransferNaive`, `TransferOrdered`, `TransferWithRetry`, `ErrDeadlock`) match code exactly.
2. README run commands — MATCH. `go test -v ./...`, `go test -race ./...`, `go run ./cmd/demo` all execute clean.
3. Design claims vs code — MATCH. Circular-wait deadlock, victim abort via ctx timeout, ID-ordered prevention, duration impact, retry recovery all implemented and tested.
4. Execution-result vs rerun — MATCH. Demo output reproduced: naive 1 victim + 1 nil, ordered nil/nil balances 1100/900, retry nil/nil balances 1100/900. Victim line order varies run to run — scheduling nondeterminism, not mismatch.
5. `TransferOrdered` "Completely prevents deadlocks" — WARNING, IMPLEMENTATION_OVERCLAIM scoped. True for 2-account lab scope; notes already limit scope to A-to-B, no 3+ table claims. No code change needed, wording scope only.
6. `TransferWithRetry` checks `context.DeadlineExceeded` — dead branch, LOW. `Lock` maps all ctx expiry to `ErrDeadlock`, so `DeadlineExceeded` never surfaces. Harmless, retry still correct.

No DOC_CODE_MISMATCH. No TEST_CLAIM_MISMATCH. No RESEARCH_IMPLEMENTATION_MISMATCH in audited scope. No FAKE_DEMO. No FAKE_BENCHMARK.
