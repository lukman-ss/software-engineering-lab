# Gap Analysis

Target Lab: `labs/24-slo-sli-error-budget`

## Identified Gaps

No critical, high, or medium gaps found.

### Minor Observation 1: Out-of-Order Timestamp Handling
- Type: `MISSING_EDGE_CASE`
- Severity: LOW
- Description: `WindowTracker.Record` expects timestamps to generally advance chronologically. If an event with an older timestamp arrives out of order, it will append a new bucket rather than merging into an existing older bucket. In this lab's scope and demo use case, events advance chronologically.
- Mitigation: Out-of-scope for simple sliding window demonstration; acceptable for lab scope.

### Minor Observation 2: Rule-specific Window Configurations
- Type: `IMPLEMENTATION_OVERCLAIM`
- Severity: LOW
- Description: In `alerting.BurnRateRule`, fields `LongWindow` and `ShortWindow` exist on the struct, but `AlertEngine` evaluates against its two pre-bound window trackers (`shortTracker` and `longTracker`) rather than dynamically sizing trackers per rule.
- Mitigation: Engine is configured with 5-minute short and 60-minute long trackers matching the demo rules; acceptable simplification.
