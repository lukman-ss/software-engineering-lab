# 06 Research Gaps: Circuit Breaker Lab

## Gap 1

Type:
OUTDATED_SOURCE

Severity:
LOW

Location:
`research/02-sources.md:9-10`, `README.md:162`

Problem:
The link to AWS Builders Library (`https://aws.amazon.com/builders-library/timeouts-retries-and-backoff-with-jitter/`) returns an HTTP 301 permanent redirect to the newer AWS Builder Center portal URL (`https://builder.aws.com/content/3EumjoZascWd1oZiEgL8ORlv3qE/timeouts-retries-and-backoff-with-jitter`).

Required Revision:
Update the URL to the direct canonical link or note the redirection.

Can Be Approved Without Fix:
YES

---

## Gap 2

Type:
IMPLEMENTATION_GAP

Severity:
MEDIUM

Location:
`internal/circuitbreaker/circuit_breaker.go:118-124`, `README.md:136`

Problem:
The documentation and research warn readers against tripping circuit breakers on 4xx client errors, but the Go implementation does not provide an error classifier or predicate to distinguish retryable/system faults (5xx, timeouts) from client faults (4xx). Any error returned by `fn()` increments `failureCount`.

Required Revision:
Explicitly document in `README.md` that the lab's Go implementation is an educational reference focusing on the state machine core, and that production circuit breakers require an error filter function (e.g. `IsFailure(error) bool`).

Can Be Approved Without Fix:
YES

---

## Gap 3

Type:
IMPLEMENTATION_GAP

Severity:
LOW

Location:
`research/08-observability.md`, `README.md:77-84`

Problem:
The research outlines metrics like `circuit_state`, `circuit_open_count`, `rejected_call_count`, and `failure_count`, but the codebase exposes only `State() State`. The counters are internal struct fields and are not exposed or exported via any metrics registry.

Required Revision:
Add a note in the Observability section clarifying that these metrics are architectural specifications for production systems rather than an implemented telemetry export in this minimal pedagogical code.

Can Be Approved Without Fix:
YES

---

## Gap 4

Type:
MISSING_TEST

Severity:
LOW

Location:
`internal/circuitbreaker/circuit_breaker_test.go`

Problem:
While test 11 covers general concurrency in CLOSED and OPEN states, there is no explicit concurrent test verifying that when multiple goroutines call `Execute()` while in `HALF_OPEN`, only `HalfOpenMaxCalls` are permitted to execute the downstream probe, while the rest are immediately rejected with `ErrCircuitOpen`.

Required Revision:
Add a specific unit test for probe rate limiting during HALF_OPEN under concurrency.

Can Be Approved Without Fix:
YES
