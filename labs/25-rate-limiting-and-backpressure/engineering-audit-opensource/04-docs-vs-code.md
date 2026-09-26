# Docs vs Code

Files: README.md, engineering/01-design.md, 02-implementation-notes.md, 03-execution-result.md, code, tests, demo.

## README Accuracy

Matching claims:

- TokenBucket + LeakyBucket burst/refill/leak behavior: code matches.
- Bounded queue TrySubmit/ErrQueueFull: code matches.
- Retry jitter (Full, Equal, No, Decorrelated): code matches.
- HTTP 429 middleware + Retry-After header: code matches.
- Structure diagram: matches file tree.
- Execution/test commands: correct.

No DOC_CODE_MISMATCH found in README.

## Design Doc Accuracy

- Component layout (ratelimit, backpressure, retry, httputil, cmd/demo): all exist.
- Success criteria (behavior + race + demo): all satisfied.
- Architecture diagram: implemented (tenant extractor, bucket gate, bounded queue, jitter retry).

Minor deviation: design mentions deterministic simulated time helpers for tests; tests use real time.Sleep, not simulated clocks. Behavior still proven, but statement inaccurate. Classification: DOC_CODE_MISMATCH (LOW).

## Implementation Notes Accuracy

- File inventory: matches actual tree.
- Design decisions (mutex, lazy refill, select-default, AWS jitter): match code.
- Known limitations (in-memory only, no distributed state): accurate.
- What Is/Is Not Demonstrated: correctly scoped.

## Execution Result Accuracy

- Build/tests/race demo claimed green: verified green on 2026-09-26.
- Demo log deterministic sections (token bucket, leaky bucket, backpressure): byte-identical to fresh run. Jitter values differ across runs (expected: math/rand randomness). Classification: sampling artifact, not FAKE_DEMO.

## Research Alignment

Per pipeline override: research/content audit excluded. Implementation rationale (HTTP 429 RFC 6585, jitter, per-tenant keys) consistent with design docs reviewed at face value.

Verdict on mapping: PASS with LOW warning on the simulated-time wording.
