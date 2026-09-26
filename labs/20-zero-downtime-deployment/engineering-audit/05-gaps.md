# Gaps and Findings

Target Lab: labs/20-zero-downtime-deployment

## Gap 1: BROKEN_IMPLEMENTATION / UNHANDLED_ERROR
- **Type**: BROKEN_IMPLEMENTATION
- **Severity**: MEDIUM
- **Description**: In `internal/worker/worker.go`, `w.Stop()` calls `w.cancel()` immediately alongside closing `w.jobChan`. As a result, worker goroutines selecting on `<-w.ctx.Done()` exit as soon as their active job finishes, leaving any already enqueued/buffered jobs unhandled. A proper queue drain should exhaust the remaining items in `w.jobChan` before termination, unless a forced timeout context occurs.

## Gap 2: MISSING_TEST
- **Type**: MISSING_TEST
- **Severity**: MEDIUM
- **Description**: In `tests/worker_test.go`, the test checks `len(completed) < 1` instead of asserting the exact expected behavior of draining enqueued tasks. This masks the dropping of buffered tasks during worker shutdown.

## Gap 3: DOC_CODE_MISMATCH
- **Type**: DOC_CODE_MISMATCH
- **Severity**: LOW
- **Description**: README claims the worker "stops pulling new jobs but continues processing the current active job until completion", which is literally what it does, but it does not specify whether already-buffered jobs in the internal channel will be abandoned or drained.
