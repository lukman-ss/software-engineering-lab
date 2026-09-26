# Test Audit

Target Lab: `labs/20-zero-downtime-deployment`

## Coverage Matrix

| Behavior                  | Test                                  | Result |
|---------------------------|---------------------------------------|--------|
| DB invalid get            | `TestDBNotFound`                      | PASS   |
| DB single-name legacy split | `TestDBSingleNameLegacy`             | PASS   |
| DB empty fields expand    | `TestDBSaveExpandEmptyFields`         | PASS   |
| DB legacy→expand overwrite | `TestDBLegacyOverwriteWithExpand`    | PASS   |
| DB legacy→expand read     | `TestExpandContractDatabase`          | PASS   |
| Probes ready/unready      | `TestServerProbes`                    | PASS   |
| In-flight drain 1 req     | `TestServerGracefulShutdown`          | PASS   |
| preStop delay honored     | `TestServerPreStopHook`               | PASS   |
| preStop ctx cancel aborts | `TestServerPreStopContextCancellation`| PASS   |
| Invalid duration fallback | `TestServerInvalidDurationFallback`   | PASS   |
| Ready↔unready toggle      | `TestServerReadyUnreadyTransition`    | PASS   |
| Multi-req drain (3 reqs)  | `TestServerMultiRequestDrain`         | PASS   |
| Client cancel → Active=0  | `TestServerWorkRequestCancellation`   | PASS   |
| Worker concurrency (3×6)  | `TestWorkerConcurrency`               | PASS   |
| Worker drain preserves order | `TestWorkerGracefulShutdown`       | PASS   |
| Enqueue-after-Stop dropped | `TestWorkerEnqueueAfterStop`         | PASS   |
| Concurrent Enqueue+Stop race safety | `TestWorkerConcurrentEnqueueStop` | PASS   |
| Drain timeout aborts long job | `TestWorkerShutdownTimeout`        | PASS   |

Happy path: covered (probes set/read, drain, expand/contract).
Failure path: covered (503 readiness, ErrNotFound, preStop cancel, drain timeout abort, client cancel).
Edge cases: covered (single-name split, empty first/last, invalid duration fallback, enqueue-after-stop).
Transitions: covered (ready↔unready, legacy→expand).
Recovery: covered (drain timeout escalation, client-cancel cleanup).
Concurrency: covered + race detector clean (`TestWorkerConcurrency`, `TestWorkerConcurrentEnqueueStop`, `TestServerMultiRequestDrain`).

## Strengths
- Suite goes beyond happy path: preStop cancellation, drain timeout, client cancellation, concurrent enqueue/stop — uncommonly thorough for a lab.
- `TestWorkerShutdownTimeout` genuinely exercises timeout escalation (slow job finishes, 100ms job aborted at 25ms timeout).
- `TestServerWorkRequestCancellation` asserts `ActiveRequests()==0` cleanup, not just no-crash.

## Weaknesses
- `TestWorkerConcurrentEnqueueStop` asserts only "no panic/race" — does not assert job outcome; weak oracle but correct for its purpose (race regression).
- No test for double `Stop()` (would panic — Finding 5 in 02-code-audit).
- No test asserting server unreachability after full `Shutdown` (new connections refused).
- DB tests are single-goroutine; no concurrent access test (implementation uses RWMutex; race detector would not catch without concurrent test).

## Execution Record (this audit, real runs)

- `go build ./...` → success.
- `go vet ./...` → clean.
- `go test -v ./...` → 18/18 PASS.
- `go test -race ./...` → ok (no data races).
- `go test -race -count=1 ./...` → ok (2.118s, uncached).
- `go run ./cmd/demo` → exit 0, output matches `engineering/03-execution-result.md` (17 lines: readiness → SIGTERM → preStop 1s → job finished → client 200 → drain complete → "Demo finished cleanly. Zero downtime achieved.").

Suite is strong, not merely passing.
