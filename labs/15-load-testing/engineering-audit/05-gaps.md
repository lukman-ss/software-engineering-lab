# Engineering Gaps

Target Lab: labs/15-load-testing

## Gap 1

- Gap Type: `RESEARCH_MISMATCH`
- Severity: HIGH
- Description: The core research hypothesis claimed that average latency conceals tail latency spikes. However, the mock server uses a deterministic query duration (`20ms`) under a closed-loop virtual user harness. Consequently, queuing delay distributes uniformly across all requests, resulting in an Average latency (~200ms) that degrades virtually identically to P95 (~211ms). The masking effect of averages is not demonstrated by the actual runtime numbers.

## Gap 2

- Gap Type: `DOC_CODE_MISMATCH`
- Severity: MEDIUM
- Description: `engineering/01-design.md` states: *"P95 and P99 latency spikes significantly, while average latency degrades less severely, proving the masking effect of averages."* The actual execution result shows average latency degrading from 21ms to 200ms, which is a 9.5x degradation, virtually matching the P95 degradation (22ms to 212ms, 9.6x).

## Gap 3

- Gap Type: `MISSING_TEST`
- Severity: LOW
- Description: No unit or integration test checks whether the percentile calculation maintains ordering invariants (Min <= P50 <= P90 <= P95 <= P99 <= Max) under non-monotonic or random distributions.
