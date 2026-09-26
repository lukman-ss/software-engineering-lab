# Engineering Revision Plan

Target Lab: labs/20-zero-downtime-deployment
Previous Verdict: APPROVED (8 non-blocking gaps; no blocking issues)

## Blocking Issues
None.

## Non-Blocking Issues
1. GAP-01 (LOW): No test for GetUser with non-existent ID returning ErrNotFound.
2. GAP-02 (LOW): No test for InsertLegacy with single-name (no space).
3. GAP-03 (LOW): No test for SaveExpand with empty firstName or lastName.
4. GAP-04 (MEDIUM): Enqueue after Stop panics; no runtime guard.
5. GAP-05 (LOW): Worker tests use concurrency=1 only.
6. GAP-06 (LOW): /work with invalid duration fallback untested.
7. GAP-07 (LOW): engineering/01-design.md describes Shutdown signature incorrectly.
8. GAP-08 (LOW): READY→UNREADY probe transition not tested as explicit probe assertion.

## Files To Change
- `tests/db_test.go` — add GAP-01, GAP-02, GAP-03 tests
- `tests/server_test.go` — add GAP-06, GAP-08 tests
- `tests/worker_test.go` — add GAP-05 test
- `internal/worker/worker.go` — add Enqueue-after-Stop guard (GAP-04)
- `engineering/01-design.md` — fix Shutdown signature description (GAP-07)

## Tests To Add/Modify
- TestDBNotFound (GAP-01)
- TestDBSingleNameLegacy (GAP-02)
- TestDBSaveExpandEmptyFields (GAP-03)
- TestWorkerConcurrency (GAP-05)
- TestServerInvalidDurationFallback (GAP-06)
- TestServerReadyUnreadyTransition (GAP-08)

## Validation Commands
```bash
go test -v ./...
go test -race -count=1 ./...
go run ./cmd/demo
```
