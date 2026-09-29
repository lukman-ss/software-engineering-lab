## Gaps

- Missing negative/error tests (DB error, context cancellation) for all services.
- WriteBehind queue overflow silently drops writes; no test verifies this behavior.
- SWR revalidation error handling not exercised (DB failure during async refresh).
- XFetch edge case when TTL remaining exactly equals computed threshold not tested.
- SingleFlight only intra-process; cross-process stampede not addressed.
- TTL jitter uses non-crypto RNG – acceptable for demo but noted.
- Documentation notes these limitations; but audit gaps require explicit test coverage for robustness.
