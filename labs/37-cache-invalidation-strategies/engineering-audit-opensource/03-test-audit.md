## Test Coverage

- Cache patterns: verified read/write counts, cache hits/misses.
- Stampede mitigation: naive >1 DB queries, singleflight ==1.
- XFetch: deterministic rand tests confirm formula.
- SWR: stale serve and async refresh verified.
- Jitter: range validation test.

All tests pass with race detector.

## Gaps
- WriteBehind queue overflow silently drops writes (MEDIUM risk).
- SWR test relies on sleeps; could be flaky under load (LOW).
- No distributed stampede protection (out of scope).
