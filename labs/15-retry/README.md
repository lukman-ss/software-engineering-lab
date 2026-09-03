# Lab 15: Retry Patterns

This lab explores safe retry mechanisms for handling transient failures.

## Problem Statement

When operations fail, naive retries can amplify problems rather than solve them.
Poor retry behavior leads to retry storms and can take down healthy systems.

## Prompt 045 — Mock Provider

The mock provider (`provider.go`) can simulate various failure modes:

- Fail N times before succeeding
- Return HTTP 500 (Internal Server Error)
- Return HTTP 429 (Too Many Requests - rate limited)
- Timeout
- Succeed

```go
provider := NewMockProvider().
    WithMaxFailures(2).
    WithStatusCode(InternalServerError)

resp, err := client.Get(ctx, url)
```

## Prompt 046 — Exponential Backoff with Jitter

The `RetryableClient` implements exponential backoff with configurable jitter:

```go
client := NewRetryableClient(provider,
    WithRetryAttempts(3),
    WithBaseDelay(100 * time.Millisecond),
    WithJitter(0.5), // Full jitter: 0% to 50% additional delay
)
```

### Retry Only Transient Errors

The client checks HTTP status codes to determine if retry is appropriate:

| Status Code | Name | Retry |
|-------------|------|-------|
| 200 | OK | No |
| 400 | Bad Request | No (validation error) |
| 401 | Unauthorized | No (credential issue) |
| 403 | Forbidden | No (permission issue) |
| 404 | Not Found | No (resource doesn't exist) |
| 429 | Too Many Requests | Yes (rate limited) |
| 500 | Internal Server Error | Yes (transient) |
| 503 | Service Unavailable | Yes (transient) |

## Prompt 047 — Retry Storm

### The Problem

When multiple clients experience failures simultaneously, they all retry at nearly the same time, creating a thundering herd that overwhelms the recovering service.

### Solution: Jitter

Jitter randomizes the backoff delay:

- **Full Jitter**: `0ms to baseDelay * 2^attempt`
- **Equal Jitter**: `0 to jitter + baseDelay`
- **Decorrelated Jitter**: Uses previous delay in calculation

The `waitDuration()` function implements full jitter with configurable factor.

## Prompt 048 — Retry Budget

### Retry Amplification Factor

In a multi-layered system, retries multiply exponentially:

```
Client retries API Gateway 3 times
  -> 1 * 3 = 3 requests

API Gateway retries Service A 3 times  
  -> 3 * 3 = 9 requests

Service A retries Service B 3 times
  -> 9 * 3 = 27 requests

Result: 27 requests hit the failing backend for every 1 original request!
```

### Retry Budget

A retry budget limits total retry attempts per time window:

```go
budget := NewRetryBudget(20, time.Minute)  // Allow 20 retries per minute

for _, req := range requests {
    if !budget.TryConsume() {
        // Budget exhausted - fail fast
        return ErrBudgetExceeded
    }
    // ... proceed with retry
}
```

### Best Practices

1. **Retry at only one layer** - Usually at the edge (client or API gateway)
2. **Use retry budgets** - Limit total retries system-wide
3. **Pass retry headers** - `X-Retry-Attempt: N` to inform downstream services
4. **Fail fast when budget exhausted** - Don't add to the storm
5. **Monitor retry metrics** - Alert on budget exhaustion patterns

## Running the Tests

```bash
go test ./labs/15-retry/... -v
```

Example output:
```
=== RUN   TestMockProviderFailModes
--- PASS: TestMockProviderFailModes (0.00s)
=== RUN   TestRetryingTransientErrors
--- PASS: TestRetryingTransientErrors (0.00s)
```

---

## Navigasi

- **Previous**: [Lab 14 — Outbox Pattern](../14-outbox-pattern/)
- **Next**: None