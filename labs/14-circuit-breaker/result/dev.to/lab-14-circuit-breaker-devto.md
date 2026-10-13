---
title: "Circuit Breakers Protect Callers from Failing Dependencies"
published: false
description: "Explore a Go checkout breaker: consecutive failures, fail-fast rejection, bounded recovery probes, generation guards, and timeout limits."
tags: go, architecture, testing, softwareengineering
---

# Circuit Breakers Protect Callers from Failing Dependencies

A payment service slows down. Checkout continues sending requests and waiting for timeouts. Lab 14 places a circuit breaker before the payment client so repeated failures can stop subsequent attempts.

> A circuit breaker protects the caller; it does not repair the dependency.

The lab uses real local HTTP through an `httptest` fake server. It does not perform real payments, balance mutations, or production gateway integration.

## The implemented path

Checkout calls `Breaker.Execute`, whose callback calls `payment.Client.ProcessPayment`. The client sends `POST /pay`. The fake server supports HEALTHY (200), DOWN (500), and SLOW (sleep, then 200).

The README motivates the pattern with cascading resource exhaustion. The tests demonstrate narrower behavior: state transitions, callback rejection, HTTP request counts, timeout errors, and recovery probes. They do not measure thread exhaustion, memory pressure, or a platform-wide outage.

## CLOSED counts consecutive failures

The breaker starts CLOSED and admits callbacks. A successful completion resets its failure count. A nonnil error increments it; reaching the threshold opens the circuit:

```go
// circuit_breaker.go: onFailureLocked excerpt
b.failures++
if b.failures >= b.cfg.FailureThreshold {
    b.state = Open
    b.openedAt = now
    b.generation++
}
```

This is a consecutive-error counter, not a sliding-window failure rate. With concurrent calls, valid completion order determines the counter updates.

Defaults are three failures, a 300 ms cooldown, and one HALF_OPEN call. `New` substitutes defaults for nonpositive values. These are lab settings, not universal production recommendations.

Callbacks execute outside the mutex. CLOSED has no general in-flight limit, so a concurrent burst can dispatch more requests than the threshold before the first failures complete. A circuit breaker is not automatically a bulkhead.

## OPEN rejects before invoking the callback

```go
// circuit_breaker.go: Execute state branch
case Open:
    b.mu.Unlock()
    return ErrCircuitOpen
```

During OPEN before cooldown expires, no callback runs. When the callback wraps the payment client, that rejected call produces no new network request. `TestOpenDoesNotCallDownstream` checks one initial failed callback followed by five rejections, with the callback count remaining one.

Already admitted calls are not canceled. “No downstream calls while rejecting” is different from claiming no older request is still running.

## Cooldown advancement is lazy

There is no background timer that sleeps and sends a probe. Both `State()` and `Execute()` call `advanceLocked`. Once enough time has elapsed since `openedAt`, the state becomes HALF_OPEN, the probe counter resets, and the generation increases.

Calling `State()` can therefore mutate state. An expired cooldown does not itself dispatch a recovery request; a later admitted `Execute` does that.

## HALF_OPEN limits admission, not the number of successes required

```go
// circuit_breaker.go: HALF_OPEN branch
case HalfOpen:
    if b.halfOpenIn >= b.cfg.HalfOpenMaxCalls {
        b.mu.Unlock()
        return ErrCircuitOpen
    }
    b.halfOpenIn++
```

Extra calls receive the same `ErrCircuitOpen` used for OPEN. The error alone does not distinguish an open circuit from a full probe allowance.

One successful probe closes the breaker. One failed probe reopens it and restarts cooldown. When the configured allowance exceeds one, the implementation does not wait for every admitted probe to succeed: the first state-changing result advances the generation, making remaining old results stale.

That is a specific recovery policy, not a requirement for N consecutive successful probes.

## Generation guards separate old work from new state

`Execute` captures a generation before releasing the mutex to run the callback. Result bookkeeping begins with this guard:

```go
// circuit_breaker.go: result bookkeeping guard
if b.generation != gen {
    return
}
```

A trailing request from an earlier generation cannot reset or reopen a newer state. The trailing-in-flight test starts a slow CLOSED request, opens the circuit using another request, advances to HALF_OPEN, and checks that the old failure does not corrupt HALF_OPEN.

The old callback still runs and its error still returns to its caller. Generation filtering changes its authority to update breaker state; it is not cancellation.

## Panic cleanup does not swallow panic

Deferred bookkeeping treats a callback panic as a failure, including reopening HALF_OPEN. The panic continues to the caller rather than becoming an ordinary returned error. The test recovers outside the breaker and checks that a later recovery probe can run.

This prevents a probe panic from leaving the breaker permanently stuck with its allowance consumed, while preserving the panic's normal propagation.

## Timeouts remain necessary

The HTTP client has a configured timeout. `Checkout(ctx)` passes the caller's context into the payment request; the two demo helpers use `context.Background()`.

`Breaker.Execute` itself has no context argument or execution deadline. An admitted callback can still wait indefinitely without its own timeout, including a HALF_OPEN probe occupying a slot.

Retry and exponential backoff are discussed in the README but not implemented in this client or checkout service. The breaker does not schedule retries or automatically control an external retry loop. It only admits or rejects each callback evaluation.

The fake SLOW handler sleeps without checking request cancellation. A client timeout is therefore not proof that the server has stopped its work.

## Failure classification is intentionally broad

The client accepts only HTTP 200. Every other status becomes `ErrPaymentFailed`; transport errors, timeouts, and cancellation also return errors. The breaker counts all nonnil errors without classification, so a 4xx or canceled request can affect the threshold if it reaches this path.

A 200 response is not decoded to validate payment business semantics. A successful probe means this client returned nil, not that a real transaction was completed correctly.

## Read the demo as four scenarios, not one benchmark

The first scenario uses SLOW with a 200 ms server delay and a 100 ms client timeout, making three sequential calls without a breaker. The second uses DOWN responses, a threshold of three, and six sequential calls with a breaker.

Those are different failure modes. They should not be presented as a controlled latency benchmark of an identical workload. The N × timeout approximation describes sequential waiting, not general concurrent request latency.

The third scenario makes the server HEALTHY before probing and demonstrates HALF_OPEN → CLOSED. The fourth keeps it DOWN and demonstrates HALF_OPEN → OPEN.

Nanosecond or microsecond durations shown in the README are example output, not a latency guarantee or a newly measured result here.

## Fallback must preserve business meaning

The implemented checkout returns an error on payment failure or rejection. It has no cached success, queue, or automatic alternative payment flow.

The README's CMMS and PPOB examples discuss critical versus noncritical dependencies. Invoice creation, WhatsApp delivery, provider fulfillment, queue retention, and debit protection are architectural discussion, not implemented workflows in this folder.

A fast payment rejection still means checkout has failed. The breaker reduces futile attempts; it does not manufacture a successful payment during an outage.

## What the tests and instrumentation establish

Unit tests cover thresholds, success resets, rejected callbacks, cooldown, both probe outcomes, probe limits, panic cleanup, stale generations, and concurrent access. Integration tests use the HTTP fake server to check 500 responses, recovery, request counts, and slow-call timeouts.

The general concurrency test checks a valid final state. The interleaved test waits for goroutines without additional semantic assertions. A race-detector run can provide evidence about executed paths, but does not independently prove every interleaving or policy.

Production gauges, open/rejected counters, and dependency latency percentiles are recommendations in the README. This implementation exposes state and the fake server's request counter, not a production metrics exporter.

```sh
cd labs/14-circuit-breaker
go test -v ./...
go test -race ./...
go run ./cmd/demo
```

The module declares Go 1.22. Go is unavailable in this authoring environment, so I inspected code and assertions rather than rerunning the tests. No new passing-test result or benchmark is claimed.

The reusable lesson is to bound permitted calls and interpret recovery precisely: timeout limits waiting, the breaker controls subsequent admission, and probes decide whether that path should reopen.

## Source

[Software Engineering Lab 14](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/14-circuit-breaker), snapshot commit `82af988b89439fa92b1200406c5afbd59ce72a56`: README, circuit breaker implementation/tests, payment client/fake server, checkout service, integration tests, demo, and go.mod. Excerpts belong in their existing package context.
