# 02 Source Audit: Circuit Breaker Lab

## Source 1

Claimed Title:
Circuit Breaker

Claimed Publisher:
Martin Fowler

URL:
https://martinfowler.com/bliki/CircuitBreaker.html

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Martin Fowler's canonical article introduces the software Circuit Breaker pattern (crediting Michael Nygard's *Release It!*), defining CLOSED, OPEN, and HALF-OPEN states, threshold counters, reset intervals, and the rationale for preventing resource exhaustion.

Assessment:
PASS

---

## Source 2

Claimed Title:
Circuit Breaker Pattern - Azure Architecture Center

Claimed Publisher:
Microsoft Learn

URL:
https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker

Reachable:
YES

Source Type:
SECONDARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Authoritative reference covering three-state lifecycle, state transition counters, interaction with the Retry pattern, canary probes in HALF-OPEN state, handling different exception types (distinguishing 4xx from 5xx/timeouts), and observability requirements.

Assessment:
PASS

---

## Source 3

Claimed Title:
Timeouts, retries, and backoff with jitter

Claimed Publisher:
AWS Builders Library

URL:
https://aws.amazon.com/builders-library/timeouts-retries-and-backoff-with-jitter/

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
PARTIAL

Problems:
- The cited URL permanently redirects (HTTP 301) to AWS Builder Center canonical URL: `https://builder.aws.com/content/3EumjoZascWd1oZiEgL8ORlv3qE/timeouts-retries-and-backoff-with-jitter`.
- The article focuses directly on client-side timeouts, retry amplification, and exponential backoff with full jitter (authored by Marc Brooker). It supports the lab's discussion of why timeouts and retries alone are insufficient without a circuit breaker to halt persistent retry storms, but it does not detail circuit breaker state machine mechanics.

Assessment:
WARNING
