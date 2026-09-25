# 06 Timeout, Retry, and Exponential Backoff

## Timeout
Bounds execution time. Prevents threads from hanging forever. Without timeouts, circuits cannot detect slow downstream hangs.

## Retry
Attempts operation again under expectation of transient glitch.
Risks: Retry amplification / Retry storms.
Mitigation: Exponential backoff with Full Jitter and max retry counts (e.g. 2-3 retries max).

## Synergy
Timeout bounds one request → Retries attempt recovery → Circuit Breaker trips when failure persists to halt the retries.
