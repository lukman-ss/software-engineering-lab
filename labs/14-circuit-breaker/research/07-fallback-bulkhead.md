# Fallback and Bulkhead

## Fallback

### Definition
A fallback is an alternative response strategy when a service dependency is unavailable (Circuit Breaker OPEN). It provides degraded but functional behavior rather than a hard error.

**Evidence**:
- Azure: "In some cases, rather than returning a failure and raising an exception, the Open state can return a default value that's meaningful to the application." — https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker
- Martin Fowler: "Breakers on their own are valuable, but clients using them need to react to breaker failures. As with any remote invocation you need to consider what to do in case of failure. Does it fail the operation you're carrying out, or are there workarounds you can do? A credit card authorization could be put on a queue to deal with later, failure to get some data may be mitigated by showing some stale data that's good enough to display." — https://martinfowler.com/bliki/CircuitBreaker.html

### Safe Fallback Strategies
| Strategy | Description | Safety Condition | Source |
|----------|-------------|------------------|--------|
| Cached response | Return previously cached data | Cache is reasonably fresh | Martin Fowler: "failure to get some data may be mitigated by showing some stale data that's good enough to display" |
| Default placeholder | Return static/default value (e.g., empty list, zero) | Placeholder is semantically safe for display | Azure: "the Open state can return a default value that's meaningful to the application" |
| Queued processing | Accept request, queue for async processing when dependency recovers | Request is not user-facing / not critical path | Martin Fowler: "A credit card authorization could be put on a queue to deal with later" |
| Degraded mode | Serve reduced-accuracy or limited-feature version | Core function preserved, non-critical data omitted | Google SRE: "serve degraded results when necessary" — https://sre.google/sre-book/handling-overload/ |

### Unsafe / Context-Dependent Fallbacks
| Strategy | Risk | When Prohibited |
|----------|------|-----------------|
| Silent success (200 OK with dummy data) | Data corruption / silent errors | Financial transactions, state mutations |
| Infinite retry queue | Resource exhaustion / delayed failure detection | Critical path user requests |
| Bypass validation | Security violations / invalid state | Any mutation requiring consistency checks |

### Financial Fallback Guidance (PPOB principle)
**Claim**: In financial systems (e.g., PPOB: Order → Payment Gateway), fallback must never silently deduct balance or alter account state without explicit confirmation.

**Evidence**:
- From topic spec: "Do not make unsafe assumptions about financial fallback behavior."
- Principle: Financial mutations must be synchronous and confirmed; they cannot rely on asynchronous fallback when the gateway is down.
- Corroborated by Martin Fowler's credit-card example implying queued retry is acceptable only if the operation is not finalized until successful processing.

### NOT VERIFIED
- Universal fallback strategy for all service types (depends on domain: financial vs content vs telemetry).
- Specific cache TTL recommendations (depends on data volatility).

---

## Bulkhead

### Definition
Bulkhead isolates resources (threads, connection pools, memory, CPU) per dependency or tenant so that exhaustion in one bulkhead does not affect others.

**Evidence**:
- Azure: "Isolate the elements of an application into pools so that if one fails, the others continue to function." — https://learn.microsoft.com/en-us/azure/architecture/patterns/bulkhead
- Azure analogy: "This pattern is named after the sectioned partitions (bulkheads) of a ship's hull. If the hull of a ship is compromised, only the damaged section fills with water, which prevents the ship from sinking."
- Azure: "A consumer can also partition resources to ensure that resources used to call one service don't affect the resources used to call another service."

### Bulkhead vs Circuit Breaker
| Aspect | Bulkhead | Circuit Breaker |
|--------|----------|-----------------|
| **Primary Goal** | Resource isolation | Failure detection + traffic gating |
| **Mechanism** | Separate pools (threads, connections) | State machine (CLOSED/OPEN/HALF_OPEN) |
| **When it helps** | Dependency A fails → its bulkhead exhausted, but bulkhead B for dependency C still has resources | Dependency fails → circuit OPEN → fail fast saves caller resources |
| **Do they conflict?** | No — complementary | No — complementary |
| **Can be combined?** | Yes | Yes |

**Evidence**:
- Azure Bulkhead: "To provide more sophisticated fault handling, consider combining bulkheads with retry, circuit breaker, and throttling patterns." — https://learn.microsoft.com/en-us/azure/architecture/patterns/bulkhead
- Azure Circuit Breaker: Notes on resource differentiation: "Be careful when you use a single circuit breaker for one type of resource if there might be multiple underlying independent providers."

### Implementation Examples
- Thread pools per dependency (Azure: "Processes, thread pools, and semaphores. Projects like resilience4j and Polly offer a framework for creating consumer bulkheads.")
- Connection pools per service (Azure diagram showing separate connection pools for Service A, B, C)
- Container/VM isolation (Azure: "When you partition services into bulkheads, consider deploying them into separate virtual machines, containers, or processes.")
- Queue-per-client for async services (Azure: "Services that communicate by using asynchronous messages can be isolated through different sets of queues.")

### NOT VERIFIED
- Specific numeric pool sizes (depends on traffic profile and SLAs — no universal recommendation).
- Whether bulkhead alone prevents cascade failure without circuit breaker (it contains resource exhaustion but does not fail-fast or provide recovery probing; circuit breaker adds these properties).