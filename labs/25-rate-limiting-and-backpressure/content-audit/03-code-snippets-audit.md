# Code Snippets Audit

## File: content/03-code-snippets.md

### Accuracy Assessment: PASS

### Files Verified: 12 snippets across 5 source files + 1 demo

### Detailed Verification:

#### Snippet 1 — Token Bucket AllowN ✓
- **Expected**: `internal/ratelimit/bucket.go:29-46`
- **Found**: Exact match — verbatim, including mutex, elapsed calculation, cap, debit
- **Severity**: PASS
- **Notes**: Clean verbatim copy

#### Snippet 2 — RetryAfterSeconds ✓
- **Expected**: `internal/ratelimit/bucket.go:55-79`
- **Found**: Exact match — verbatim, including ceil-rounding, zero-charge check
- **Severity**: PASS
- **Notes**: Correctly documents the `refillRate > 0` precondition (discloses the zero-rate panic gap)

#### Snippet 3 — Leaky Bucket Allow ✓
- **Expected**: `internal/ratelimit/bucket.go:98-115`
- **Found**: Exact match — verbatim, water drain, floor at 0, capacity check
- **Severity**: PASS
- **Notes**: Clean verbatim copy

#### Snippet 4 — Registry per-tenant (Double-Checked Locking) ✓
- **Expected**: `internal/ratelimit/registry.go:21-37`
- **Found**: Exact match — verbatim, RLock fast path + Lock fallback with re-check
- **Severity**: PASS
- **Notes**: Correct pattern identification

#### Snippet 5 — BoundedQueue TrySubmit ✓
- **Expected**: `internal/backpressure/queue.go:67-85`
- **Found**: Exact match — verbatim, stopped check, select with ctx.Done/queue/default
- **Severity**: PASS
- **Notes**: Correctly identifies select-default as core backpressure signal

#### Snippet 6 — Worker Loop ✓
- **Expected**: `internal/backpressure/queue.go:49-63`
- **Found**: Exact match — verbatim
- **Severity**: PASS
- **Notes**: Transparently notes `_ = job(...)` error-discard as known gap

#### Snippet 7 — Full Jitter ✓
- **Expected**: `internal/retry/backoff.go:28-42`
- **Found**: Exact match — verbatim, temp = min(cap, base×2^attempt), rand.Float64×temp
- **Severity**: PASS
- **Notes**: Correctly notes AWS default constants (50ms/1000ms/20s) are AWS-specific

#### Snippet 8 — Equal Jitter ✓
- **Expected**: `internal/retry/backoff.go:44-48`
- **Found**: Exact match — verbatim, half + rand half
- **Severity**: PASS
- **Notes**: Clean copy with correct boundary description

#### Snippet 9 — Decorrelated Jitter ✓
- **Expected**: `internal/retry/backoff.go:50-58`
- **Found**: Exact match — verbatim, prev×3 range, min(cap)
- **Severity**: PASS
- **Notes**: Accurately flags decorrelated statistical properties as not proved by tests (LOW gap)

#### Snippet 10 — Middleware 429 RFC 6585 + Retry-After ✓
- **Expected**: `internal/httputil/middleware.go:17-40`
- **Found**: Exact match — verbatim, X-API-Key extraction, fallback to anonymous, 429 + Retry-After + JSON body
- **Severity**: PASS
- **Notes**: Correctly notes header vs body field consistency gap exists and is not yet tested (LOW)

#### Snippet 11 — Demo: Token Bucket Burst ✓
- **Expected**: `cmd/demo/main.go:15-22`
- **Found**: Exact match — verbatim, cap 3 refill 5/s sequence
- **Severity**: PASS
- **Notes**: Correctly ties 300ms pause to ~0.5-token refill (matches design doc Execution Result)

#### Snippet 12 — Demo: Backpressure Load Shedding ✓
- **Expected**: `cmd/demo/main.go:32-51`
- **Found**: Exact match — verbatim, cap 3 workers 1, TrySubmit loop with 50ms sleep jobs
- **Severity**: PASS
- **Notes**: Correctly identifies ErrQueueFull as the rejection signal

### Formatting & Language:
- **Line 3**: Accurately states "Semua snippet diambil verbatim dari implementasi yang lulus Engineering Audit (APPROVED)" — verified, all snippets are verbatim and engineering audit verdict is APPROVED
- **All explanations**: Use Indonesian consistently, no platform bias, correctly attribute known gaps without overstating them

### Issues Found: None

### Hallucinations: None

### Missing Sections: None — snippet coverage spans all major components

### Verdict: PASS