# Timeout, Retry, Backoff

## Timeout

### Why outbound calls require timeout
**Claim**: Go `http.Client` without an explicit `Timeout` field can block indefinitely on a hung request.

**Evidence**: Go standard library `net/http` documentation: https://pkg.go.dev/net/http#Client — "Timeout specifies a time limit for requests made by this Client. The timeout includes connection time, any redirects, and reading the response body. The timer remains running after Get, Head, Post, or Do return and will interrupt reading of the Response.Body." Default is zero (no timeout).

**Confidence**: HIGH — official Go documentation.

**Note on lab timeouts**: The numeric timeouts used in this lab (e.g., HTTP client timeout, breaker cooldown) are illustrative values chosen for fast, deterministic demonstration. They are not universal production recommendations. Production timeouts must be tuned against actual service SLAs, P99 latency, and dependency recovery profiles.

### Timeout vs Circuit Breaker (distinct purposes)
**Claim**: Timeout bounds the wait duration of a single request. Circuit Breaker bounds the number of calls sent to a failing dependency. They are complementary, not substitutes.

**Evidence**:
- Azure: "The Circuit Breaker pattern serves a different purpose than the Retry pattern... The Circuit Breaker pattern prevents an application from performing an operation that's likely to fail." Timeout is described separately as a mechanism that "can block concurrent requests... until the time-out period expires. These blocked requests might hold critical system resources." — https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker
- Azure explicitly warns: "Inappropriate time-outs on external services: A circuit breaker might not fully protect applications from failures in external services that have long time-out periods."

**Confidence**: HIGH.

---

## Retry

### Retryable vs non-retryable failures
**Claim**: Retry is appropriate only for transient faults, not for non-transient/business-logic errors.

**Evidence**:
- Azure Retry: "Cancel. If the fault indicates that the failure isn't transient or is unlikely to be successful if repeated, the application should cancel the operation and report an exception." — https://learn.microsoft.com/en-us/azure/architecture/patterns/retry
- Azure Retry: "not useful... For handling failures that aren't due to transient faults, such as internal exceptions caused by errors in the business logic of an application."

**Confidence**: HIGH.

### Bounded retry
**Claim**: Retry attempts must be bounded (maximum attempts / retry budget), otherwise a failing dependency receives amplified load.

**Evidence**:
- Google SRE: "we implement a per-request retry budget of up to three attempts." — https://sre.google/sre-book/handling-overload/
- Google SRE: "we implement a per-client retry budget. Each client keeps track of the ratio of requests that correspond to retries. A request will only be retried as long as this ratio is below 10%."
- Azure: "If a request still fails after a significant number of retries, it's better for the application to prevent further requests going to the same resource and report a failure immediately."

**Confidence**: HIGH.

### Exponential backoff and jitter
**Claim**: Retry delays should increase (exponentially) and include randomization (jitter) to prevent synchronized retry bursts (thundering herd).

**Evidence**:
- Azure Retry: "the period between retries should be chosen to spread requests from multiple instances of the application as evenly as possible to reduce the chance of a busy service continuing to be overloaded." and "If necessary, this process can be repeated with increasing delays between retry attempts... The delay can be increased incrementally or exponentially."
- AWS Builders Library "Timeouts, retries, and backoff with jitter": URL https://aws.amazon.com/builders-library/timeouts-retries-and-backoff-with-jitter/ — published guidance on jitter algorithms. *(Note: page content could not be fully rendered during this research session; cited per topic spec. Confidence: MEDIUM — source URL is authoritative but body not fully inspected.)*

**Confidence**: HIGH for exponential backoff + spread-evenly guidance (Azure, inspected). MEDIUM for full jitter algorithm specifics (AWS source URL known, body not fully retrieved).

### Retry amplification
**Claim**: Retries without budgets multiply traffic to a failing service, worsening overload.

**Evidence**:
- Google SRE: "Let X be the total rate of requests attempted against the datacenter... the number of requests will grow significantly, to somewhere just below 3X." — https://sre.google/sre-book/handling-overload/
- Google SRE on layered stacks: "requests should only be retried at the layer immediately above the layer that is rejecting them... If multiple layers retried, we'd have a combinatorial explosion."
- Azure: "An aggressive retry policy with minimal delay between attempts, and a large number of retries, could further degrade a busy service that's running close to or at capacity."

**Confidence**: HIGH.

### Idempotency
**Claim**: Retries are only safe for idempotent operations, or with idempotency keys.

**Evidence**:
- Azure Retry: "Consider whether the operation is idempotent. If so, it's inherently safe to retry. Otherwise, retries could cause the operation to be executed more than once, with unintended side effects." — https://learn.microsoft.com/en-us/azure/architecture/patterns/retry

**Confidence**: HIGH.

---

## Interaction: Retry × Circuit Breaker
**Claim**: Circuit breaker should stop retries when the dependency is in a non-recovered state; retry logic should inspect circuit-breaker exceptions.

**Evidence**:
- Azure Retry: "the retry logic should be sensitive to any exceptions that the circuit breaker returns and stop retry attempts if the circuit breaker indicates that a fault isn't transient." — https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker
- Azure Circuit Breaker: "The Retry pattern enables an application to retry an operation with the expectation that it eventually succeeds. The Circuit Breaker pattern prevents an application from performing an operation that's likely to fail."

**Confidence**: HIGH.

## NOT VERIFIED
- Universal production numeric timeout recommendations (none exists; depends on SLA).
- Exact jitter algorithm formulas from AWS page (page body not fully retrieved this session).