## DOC_CODE_MISMATCH 1

Location: labs/25-rate-limiting-and-backpressure/engineering/03-execution-result.md lines 71-77
Claimed (from demo output): 
```
Job #1: ACCEPTED into bounded buffer
Job #2: ACCEPTED into bounded buffer
Job #3: ACCEPTED into bounded buffer
Job #4: REJECTED (Backpressure Shedding: backpressure: queue capacity exceeded)
Job #5: REJECTED (Backpressure Shedding: backpressure: queue capacity exceeded)
Job #6: REJECTED (Backpressure Shedding: backpressure: queue capacity exceeded)
Stats: Accepted=3, Rejected=3, Processed=1
```
Observed Actual Demo Output (live run):
```
Job #1: ACCEPTED into bounded buffer
Job #2: ACCEPTED into bounded buffer
Job #3: ACCEPTED into bounded buffer
Job #4: REJECTED (Backpressure Shedding: backpressure: queue capacity exceeded)
Job #5: ACCEPTED into bounded buffer
Job #6: REJECTED (Backpressure Shedding: backpressure: queue capacity exceeded)
Stats: Accepted=4, Rejected=2, Processed=1
```
Assessment: MISMATCH (Job #5 differed; stats off)
Severity: MEDIUM
Notes: Race condition in demo stats read: `Stats()` called after 100ms sleep while worker may still be processing; accepted count can exceed expected because worker not yet incremented processed. Backpressure logic correct; demo timing artifact.

## DOC_CODE_MISMATCH 2

Location: README.md lines 40-43 (Token Bucket)
Claim: `TokenBucket`: Allows bursts up to capacity $B$, refilling continuously at rate $R$.
Code: internal/ratelimit/bucket.go lines 25-46: exact match.
Assessment: PASS

## DOC_CODE_MISMATCH 3

Location: README.md lines 44-45 (Leaky Bucket)
Claim: `LeakyBucket`: Enforces constant drain rate $R$, rejecting bursts when water reaches capacity.
Code: internal/ratelimit/bucket.go lines 81-115: matches.
Assessment: PASS

## DOC_CODE_MISMATCH 4

Location: README.md line 48 (Bounded Queue)
Claim: Non-blocking submission (`TrySubmit`) returning fast `ErrQueueFull` when buffer capacity is reached.
Code: internal/backpressure/queue.go lines 64-81: select-default; matches.
Assessment: PASS

## DOC_CODE_MISMATCH 5

Location: README.md line 49 (Retry)
Claim: Full Jitter, Equal Jitter, No Jitter, and Decorrelated Jitter algorithms matching AWS Architecture specifications.
Code: internal/retry/backoff.go lines 24-62: formulas match Marc Brooker AWS spec.
Assessment: PASS

## DOC_CODE_MISMATCH 6

Location: README.md lines 51-52 (HTTP 429 Middleware)
Claim: Enforces rate limits returning standard RFC 6585 `429 Too Many Requests` response with `Retry-After` header.
Code: internal/httputil/middleware.go lines 17-40: sets header, status 429, JSON body.
Assessment: PASS

## RESEARCH_IMPLEMENTATION_MISMATCH

Out of scope per pipeline override.

## TEST_CLAIM_MISMATCH

Location: internal/ratelimit/bucket_test.go:83-98 concurrency test lacks postcondition.
Claim: test verifies concurrency safety.
Observed: runs 50 goroutines x 10 Allow() but never checks final token count or correctness; only checks no panic.
Assessment: Test weak but not mismatch — passes; see TEST_AUDIT Finding 5.