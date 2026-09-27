# Gap Analysis

MISSING_TEST
- No test for relay double-start / multiple goroutine leak.
- No test for Stop() called twice (panic on close of closed channel).
- No test for PurgeProcessedOutbox called concurrently with relay polling.
- No test for broker failure >1 time in relay loop (only SetFailNext used once in demo/test).
- No test verifying that a crashed relay after Publish but before MarkOutbox results in message redelivery (at-least-once) and that consumer deduplicates (requires simulating crash — hard in unit test, but could assert state).
- No test asserting message ordering (not required but worth noting nonexistent).

BROKEN_IMPLEMENTATION
- None: all core behavior compiles, tests pass, demo runs.

DOC_CODE_MISMATCH
- design.md claims SQLite table-based persistence; actual code uses in-memory maps (intentional per impl notes but not reflected in design).
- design.md architecture diagram shows Broker -> Consumer delivery; actual code: consumer pulls from broker manually.

TEST_CLAIM_MISMATCH
- design.md test strategy claims "Broker failure retry mechanism" tested; no test injects broker failure into relay loop and asserts eventual success after recovery (SetFailNext not used in relay path).

RACE_CONDITION
- None detected by `go test -race`; however Relay Start/Stop race on channel close exists (but not exercised by tests).

UNHANDLED_ERROR
- Relay.PollAndDispatch logs but does not propagate broker or DB mark errors upstream; caller ignores return value (dispatched count only). This could hide persistent failures.

MISSING_EDGE_CASE
- Very short poll interval (0 ns) not tested; could cause CPU spin.
- Very long poll interval not tested; no context cancellation.

IMPLEMENTATION_OVERCLAIM
- README claims "decoupled polling relay dispatch to a message broker" (true) and "downstream consumer idempotency" (true) — no overclaim.

RESEARCH_MISMATCH
- Design doc describes SQLite; code is in-memory map DB. Override waived research audit but we note for completeness.

FAKE_DEMO
- Demo output genuine; matches observed state.

FAKE_BENCHMARK
- No benchmarks present.

UNVERIFIED_RESULT
- All results verified by re-running tests and demo.