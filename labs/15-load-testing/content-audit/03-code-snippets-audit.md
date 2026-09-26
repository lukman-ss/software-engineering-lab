# Audit of 03-code-snippets.md

## Summary
This file provides three code snippets with explanations: Server Connection Pool Semaphore, Lock-Free Concurrent Load Runner, and Percentile Calculation.

## Accuracy Assessment

### Snippet 1 - Server Connection Pool Semaphore (lines 1-48)
- **Source**: internal/server/server.go
- **Code Match**: EXACT - lines 7-45 match server.go lines 56-94
- **Explanation**: Correctly describes semaphore blocking, context cancellation handling, time.Timer vs time.Sleep for cancellation support, and the 10% random slowdown amplifier. Note: "25x lebih lama (20ms → 500ms)" is correct for the demo config (20ms * 25 = 500ms).

### Snippet 2 - Lock-Free Concurrent Load Runner (lines 50-101)
- **Source**: internal/loadtest/runner.go
- **Code Match**: EXACT - lines 55-98 match runner.go lines 58-101
- **Explanation**: Correctly describes per-VU slice approach avoiding mutex contention, and the critical limitation that only successful (HTTP 2xx) requests have latencies recorded.

### Snippet 3 - Percentile Calculation (lines 102-136)
- **Source**: internal/loadtest/metrics.go
- **Code Match**: EXACT - lines 107-132 match metrics.go lines 23-65 (with minor elision marked by "// ...")
- **Explanation**: Correctly describes O(N log N) sorting approach, percentile index calculation, and the ponytail note about upgrading to streaming histograms for production scale.

## Issues Found

No issues found. All snippets are exact copies of the source code with accurate explanations.

## Conclusion
The code snippets are faithful representations of the actual implementation with accurate technical explanations.