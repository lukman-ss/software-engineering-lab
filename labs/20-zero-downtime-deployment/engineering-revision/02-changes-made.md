## Revision 1

Audit Issue: GAP-01 — No test for GetUser with non-existent ID
Severity: LOW
Files Changed: `tests/db_test.go`
Action: Added `TestDBNotFound` asserting `errors.Is(err, db.ErrNotFound)` for unknown ID.
Verification: PASS
Status: RESOLVED

## Revision 2

Audit Issue: GAP-02 — No test for InsertLegacy with single-name (no space)
Severity: LOW
Files Changed: `tests/db_test.go`
Action: Added `TestDBSingleNameLegacy` inserting "Madonna" and asserting FirstName="Madonna", LastName="".
Verification: PASS
Status: RESOLVED

## Revision 3

Audit Issue: GAP-03 — No test for SaveExpand with empty firstName or lastName
Severity: LOW
Files Changed: `tests/db_test.go`
Action: Added `TestDBSaveExpandEmptyFields` covering SaveExpand("","Smith") and SaveExpand("Jane","") cases.
Verification: PASS
Status: RESOLVED

## Revision 4

Audit Issue: GAP-04 — Enqueue after Stop panics (send on closed channel)
Severity: MEDIUM
Files Changed: `internal/worker/worker.go`
Action: Added `stopped atomic.Bool` field; `Stop()` sets it before closing channel; `Enqueue()` checks flag and logs+returns instead of panicking.
Verification: `go build ./...` PASS; `go test -race -count=1 ./...` PASS
Status: RESOLVED

## Revision 5

Audit Issue: GAP-05 — Worker tests use concurrency=1 only
Severity: LOW
Files Changed: `tests/worker_test.go`
Action: Added `TestWorkerConcurrency` starting 3 workers, enqueuing 6 jobs, asserting all 6 complete.
Verification: PASS
Status: RESOLVED

## Revision 6

Audit Issue: GAP-06 — /work invalid duration fallback untested
Severity: LOW
Files Changed: `tests/server_test.go`
Action: Added `TestServerInvalidDurationFallback` sending `?d=INVALID`, asserting HTTP 200 and elapsed >= 50ms fallback.
Verification: PASS
Status: RESOLVED

## Revision 7

Audit Issue: GAP-07 — engineering/01-design.md describes Shutdown signature incorrectly
Severity: LOW
Files Changed: `engineering/01-design.md`
Action: Updated line 47 from `Shutdown(ctx, preStopDelay)` to `Shutdown(ctx context.Context) error` with note that preStopDelay is a constructor parameter.
Verification: Documentation matches implementation.
Status: RESOLVED

## Revision 8

Audit Issue: GAP-08 — READY→UNREADY probe transition not tested explicitly
Severity: LOW
Files Changed: `tests/server_test.go`
Action: Added `TestServerReadyUnreadyTransition` asserting SetReady(true) returns 200, then SetReady(false) returns 503.
Verification: PASS
Status: RESOLVED
