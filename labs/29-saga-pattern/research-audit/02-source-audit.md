# Source Audit: Saga Pattern Research

## Source 1

Claimed Title: Saga Design Pattern
Claimed Publisher: Microsoft Azure Architecture Center
URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/saga

Reachable:
YES (HTTP/2 200 OK)

Source Type:
PRIMARY (Official Cloud Provider Architecture Documentation)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Authoritative reference detailing context, choreography vs orchestration, transaction types (compensable, pivot, retryable), and 6 countermeasures.

Assessment:
PASS

---

## Source 2

Claimed Title: Pattern: Saga
Claimed Publisher: Chris Richardson / Microservices.io
URL: https://microservices.io/patterns/data/saga.html

Reachable:
YES (HTTP/2 200 OK)

Source Type:
PRIMARY (Authoritative reference on microservices architecture patterns)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Explicitly covers database-per-service context, choreography vs orchestration, lack of isolation, and transactional outbox tie-in.

Assessment:
PASS

---

## Source 3

Claimed Title: Microservices Patterns (Book Reference)
Claimed Publisher: Chris Richardson / Manning Publications
URL: https://livebook.manning.com/book/microservices-patterns/chapter-4/143

Reachable:
YES (HTTP/1.1 200 OK)

Source Type:
SECONDARY (Published professional textbook reference)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Full text requires Manning subscription/login beyond free preview chapters, but canonical book citation and chapter 4 outline are verifiable.

Assessment:
PASS

---

## Source 4

Claimed Title: Garcia-Molina and Salem - "Sagas" Paper (Original Reference)
Claimed Publisher: ACM SIGMOD Record
URL: https://dl.acm.org/doi/10.1145/62224.62226

Reachable:
PARTIAL (ACM returns HTTP 403 / Cloudflare challenge to automated crawlers; canonical paper DOI is authentic)

Source Type:
PRIMARY (Foundational academic paper, 1987)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Paper originally formulated Sagas for long-lived transactions in single centralized database environments (LLTs), not distributed microservices. Research acknowledges this limitation in `06-open-questions.md`.

Assessment:
PASS

---

## Source 5

Claimed Title: Cloud-Native Patterns - Saga Pattern (Microsoft Learn)
Claimed Publisher: Microsoft Learn
URL: https://learn.microsoft.com/en-us/dotnet/architecture/cloud-native/saga-pattern

Reachable:
NO (HTTP/2 404 Not Found)

Source Type:
UNKNOWN / DEAD LINK

Relevant:
NO

Supports Claimed Topic:
NO

Problems:
- The cited URL returns 404. Microsoft documentation path has either been moved, renamed, or retired.

Assessment:
FAIL
