# Documentation vs Code

Target Lab: labs/15-load-testing

## Documented Claims
1. "Smoke Load (low VUs): All requests process within normal latency limits. Average and P95 are close." (01-design.md)
2. "Stress Load (high VUs): ... P95 and P99 latency spikes significantly, while average latency degrades less severely" (01-design.md)
3. "Failure Scenario: Under excessive concurrent load, queuing ... causes high latency and timeouts for the 95th percentile." (01-design.md)

## Code Reality
1. MATCH: Smoke test shows Avg 24ms, P95 39ms. No queuing observed.
2. MATCH: Stress test shows Avg 474ms, P95 736ms, P99 856ms. Non-linear degradation strictly proven.
3. MISMATCH (DOC_CODE_MISMATCH): The design claims "timeouts for the 95th percentile", but the demo output shows 0 errors. The `Runner` client timeout is 5 seconds, but the test duration is only 2 seconds. Therefore, timeouts are technically impossible to trigger in the demo.

## Result
Mismatch found. Documentation overclaims timeouts in the stress test, which the implementation does not actually trigger due to test duration constraints. Latency degradation is perfectly proven, but timeouts are not.