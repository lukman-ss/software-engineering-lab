# Evidence

## Evidence 1
Claim: Timeouts block concurrent requests, exhausting critical system resources (memory, threads, DB connections) resulting in cascading failures.
Evidence: "If a service is busy, failure in one part of the system might lead to cascading failures. [...] this strategy can block concurrent requests to the same operation until the time-out period expires. These blocked requests might hold critical system resources, such as memory, threads, and database connections."
Source: Microsoft Azure Architecture Center, Circuit Breaker Pattern
URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker
Confidence: HIGH
Corroborated By: Martin Fowler ("you can run out of critical resources leading to cascading failures across multiple systems")

## Evidence 2
Claim: Circuit Breakers implement a three-state machine: CLOSED, OPEN, HALF_OPEN.
Evidence: "You can implement the proxy as a state machine that includes the following states... Closed... Open... Half-Open"
Source: Microsoft Azure Architecture Center, Circuit Breaker Pattern
URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker
Confidence: HIGH
Corroborated By: Martin Fowler ("There is now a third state present - half open - meaning the circuit is ready to make a real call as trial")

## Evidence 3
Claim: Circuit Breakers differentiate from Retry Patterns by actively *preventing* an operation from occurring instead of blindly repeating it.
Evidence: "The Retry pattern enables an application to retry an operation with the expectation that it eventually succeeds. The Circuit Breaker pattern prevents an application from performing an operation that's likely to fail."
Source: Microsoft Azure Architecture Center, Circuit Breaker Pattern
URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker
Confidence: HIGH
Corroborated By: Martin Fowler ("The basic idea behind the circuit breaker is very simple. You wrap a protected function call in a circuit breaker object... Once the failures reach a certain threshold, the circuit breaker trips, and all further calls... return with an error, without the protected call being made at all.")

## Evidence 4
Claim: Half-Open limits traffic to probe whether the downstream service has recovered.
Evidence: "Half-Open: A limited number of requests from the application are allowed to pass through and invoke the operation. If these requests are successful... switches to the Closed state."
Source: Microsoft Azure Architecture Center, Circuit Breaker Pattern
URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker
Confidence: HIGH
Corroborated By: Martin Fowler

## Evidence 5
Claim: Retries without bounding/jitter cause amplified failures (Retry Storms) against struggling downstream systems.
Evidence: Retries increase total request volume onto a downstream dependency. If a dependency is overloaded, sudden retry spikes further exacerbate the outage. (AWS Builder's Library context)
Source: AWS Builder's Library - Timeouts, retries, and backoff with jitter
URL: https://aws.amazon.com/builders-library/timeouts-retries-and-backoff-with-jitter/
Confidence: HIGH
Corroborated By: General Software Engineering consensus.
