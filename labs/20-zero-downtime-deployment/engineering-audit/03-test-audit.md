# Test Audit

Target Lab: labs/20-zero-downtime-deployment

## Test Suite Coverage Overview

- `TestExpandContractDatabase`: Tests happy path for reading legacy rows and expanded rows. Correctly validates logic of db.go fallback.
- `TestServerProbes`: Tests live and ready probes change state correctly.
- `TestServerGracefulShutdown`: Tests that in-flight requests are allowed to complete during server shutdown using `time.Sleep` to simulate long work.
- `TestServerPreStopHook`: Tests that `Shutdown` correctly pauses for the preStop duration.
- `TestWorkerGracefulShutdown`: Tests worker starts, executes two short jobs, and shuts down gracefully.

## Weaknesses

1. **Worker Shutdown Flaw Not Caught**: The worker test enqueues two 50ms jobs, waits 10ms, and calls `Stop()`. Since concurrency is 1, `job-1` is processing. `job-2` is in the buffer. If `Stop()` is called, `w.cancel()` happens. `job-1` finishes. Then `w.ctx.Done()` is selected over `w.jobChan`, dropping `job-2`. The test asserts `len(completed) < 1` and uses `len(completed) >= 1` as passing criteria (via `if len(completed) < 1 { t.Fatalf... }`). It does not assert that *both* jobs complete, thereby hiding the bug that `job-2` is dropped.
2. **Race condition in worker test logic**: Testing for strictly 1 completed job means if the worker was faster, it might complete 2. But the test does not check if the system safely drains *all* accepted buffered work.

## Verdict
PASS but with WARNING for weak assertions masking a dropped-work bug in worker shutdown.
