# Evidence

## Evidence 1
Claim: Timeouts block concurrent requests, exhausting critical system resources (memory, threads, DB connections) resulting in cascading failures.
Evidence: "If a service is busy, failure in one part of the system might lead to cascading failures. [...] this strategy can block concurrent requests to the same operation until the time-out period expires. These blocked requests might hold critical system resources, such as memory, threads, and database connections."
Source: Microsoft Azure Architecture Center, Circuit Breaker Pattern
URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker
Published: 2025-02-05
Accessed: 2026-09-25
Confidence: HIGH
Corroborated By: Martin Fowler ("you can run out of critical resources leading to cascading failures across multiple systems")
URL: https://martinfowler.com/bliki/CircuitBreaker.html
Published: 2014-03-06

## Evidence 2
Claim: Circuit Breakers implement a three-state machine: CLOSED, OPEN, HALF_OPEN.
Evidence: "You can implement the proxy as a state machine that includes the following states... Closed... Open... Half-Open"
Source: Microsoft Azure Architecture Center, Circuit Breaker Pattern
URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker
Published: 2025-02-05
Accessed: 2026-09-25
Confidence: HIGH
Corroborated By: Martin Fowler ("There is now a third state present - half open - meaning the circuit is ready to make a real call as trial")
URL: https://martinfowler.com/bliki/CircuitBreaker.html
Published: 2014-03-06
Corroborated By: Netflix Hystrix Wiki (How it Works)
"1. Assuming the volume across a circuit meets a certain threshold... 2. And assuming that the error percentage exceeds the threshold error percentage... 3. Then the circuit-breaker transitions from CLOSED to OPEN. 4. While it is open, it short-circuits all requests... 5. After some amount of time... the next single request is let through (this is the HALF-OPEN state)."
URL: https://github.com/Netflix/Hystrix/wiki/How-it-Works
Published: 2017-07-03

## Evidence 3
Claim: Circuit Breakers differentiate from Retry Patterns by actively *preventing* an operation from occurring instead of blindly repeating it.
Evidence: "The Retry pattern enables an application to retry an operation with the expectation that it eventually succeeds. The Circuit Breaker pattern prevents an application from performing an operation that's likely to fail."
Source: Microsoft Azure Architecture Center, Circuit Breaker Pattern
URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker
Published: 2025-02-05
Accessed: 2026-09-25
Confidence: HIGH
Corroborated By: Martin Fowler ("The basic idea behind the circuit breaker is very simple. You wrap a protected function call in a circuit breaker object... Once the failures reach a certain threshold, the circuit breaker trips, and all further calls... return with an error, without the protected call being made at all.")
URL: https://martinfowler.com/bliki/CircuitBreaker.html
Published: 2014-03-06

## Evidence 4
Claim: Half-Open limits traffic to probe whether the downstream service has recovered.
Evidence: "Half-Open: A limited number of requests from the application are allowed to pass through and invoke the operation. If these requests are successful... switches to the Closed state."
Source: Microsoft Azure Architecture Center, Circuit Breaker Pattern
URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker
Published: 2025-02-05
Accessed: 2026-09-25
Confidence: HIGH
Corroborated By: Netflix Hystrix Wiki (How it Works)
"After some amount of time (HystrixCommandProperties.circuitBreakerSleepWindowInMilliseconds()), the next single request is let through (this is the HALF-OPEN state). If the request fails, the circuit-breaker returns to the OPEN state for the duration of the sleep window. If the request succeeds, the circuit-breaker transitions to CLOSED and the logic in **1.** takes over again."
URL: https://github.com/Netflix/Hystrix/wiki/How-it-Works
Published: 2017-07-03

## Evidence 5
Claim: Retries without bounding/jitter cause amplified failures (Retry Storms) against struggling downstream systems.
Evidence: "An aggressive retry policy with minimal delay between attempts, and a large number of retries, could further degrade a busy service that's running close to or at capacity."
Source: Microsoft Azure Architecture Center, Retry pattern
URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/retry
Published: 2024-07-18
Accessed: 2026-09-25
Confidence: HIGH
Corroborated By: Google SRE Book, Handling Overload (Chapter 21)
"A large subset of backend tasks in the datacenter are overloaded. [...] It's much more typical that only a small portion of tasks become overloaded, in which case the preferred response is to retry the request immediately."
URL: https://sre.google/sre-book/handling-overload/
Published: 2016
Accessed: 2026-09-25
Corroborated By: AWS Architecture Blog on Exponential Backoff and Jitter
"Without jitter, the exponential backoff approach is the clear loser. It not only takes more work, but also takes more time than the jittered approaches."
URL: https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/
Published: 2015-03-04
Accessed: 2026-09-25

## Evidence 6
Claim: Bulkhead isolates resources per dependency; Circuit Breaker fails fast based on error thresholds; they solve different problems.
Evidence: "Bulkhead isolates consumers and services from cascading failures. A problem that affects a consumer or service can be isolated within its own bulkhead to prevent the entire solution from failing."
Source: Microsoft Azure Architecture Center, Bulkhead pattern
URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/bulkhead
Published: 2026-03-19
Accessed: 2026-09-25
Confidence: HIGH
Corroborated By: Netflix Hystrix Wiki (Isolation)
"Hystrix employs the bulkhead pattern to isolate dependencies from each other and to limit concurrent access to any one of them."
URL: https://github.com/Netflix/Hystrix/wiki/How-it-Works#isolation
Published: 2017-07-03

## Evidence 7
Claim: Observability metrics for circuit breakers include state, failure count, rejected calls, latency, and open count.
Evidence: "Circuit breakers should provide clear observability into both failed and successful requests so that operations teams can assess system health."
Source: Microsoft Azure Architecture Center, Circuit Breaker Pattern
URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker
Published: 2025-02-05
Accessed: 2026-09-25
Confidence: HIGH
Corroborated By: Netflix Hystrix Wiki (Metrics and Monitoring)
"Hystrix reports successes, failures, rejections, and timeouts to the circuit breaker, which maintains a rolling set of counters that calculate statistics."
URL: https://github.com/Netflix/Hystrix/wiki/How-it-Works#metrics-and-monitoring
Published: 2017-07-03