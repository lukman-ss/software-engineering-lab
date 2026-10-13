---
title: "Load Testing: Read the Queue Behind Successful Bookings"
published: false
description: "Explore a Go load test with bounded work slots, closed-loop VUs, RPS, tail percentiles, and the limits of success-only latency samples."
tags: go, testing, performance, softwareengineering
---

# Load Testing: Read the Queue Behind Successful Bookings

A booking endpoint can keep returning success while concurrent users wait much longer. Lab 15 makes that behavior visible with a local HTTP server, a constrained work pool, and a small Go load generator.

> Load testing reveals how latency changes when demand competes for limited resources.

## Fast bookings can queue under load

One successful booking does not explain concurrent behavior. Lab 15 sends POST /booking to a local HTTP server, comparing low and high concurrency to expose queueing and latency distributions. It models bounded work slots, not a production database.

## Two scenarios share one server

The demo uses five simulated database slots and a 20 ms base query. Smoke runs two VUs; stress runs 50. Both run for two seconds, sequentially, with the same endpoint and booking body. These are demonstration settings, not production capacity targets.

## Each VU waits before sending again

Each VU has one goroutine and sends its next request after the previous attempt finishes. There is no think time, ramp-up, or arrival-rate scheduler. This is a closed loop: slower responses reduce each VU's request rate.

## A buffered channel bounds simulated work

A buffered channel acts as a semaphore. Handlers acquire a token before simulated query work and release it with defer. Full capacity makes additional handlers wait; there is no explicit bounded waiting queue or overload rejection policy.

```go
select {
case s.semaphore <- struct{}{}:
case <-r.Context().Done():
    return
}
defer func() { <-s.semaphore }()
```

## The database pool is simulated

There is no SQL driver, real query, booking table, or database transaction. A timer simulates work. ActiveConnections returns len(semaphore): acquired slots, not waiting handlers, TCP connections, or real database telemetry.

## Random slow work amplifies queueing

activeReq includes waiting and working POST handlers. After admission, when activeReq exceeds capacity, a 10% random branch multiplies query duration by 25. The demo's 20 ms becomes 500 ms. Stress results therefore combine queueing with deliberately injected slow work.

```go
if atomic.LoadInt64(&s.activeReq) >
    int64(s.cfg.MaxDBConnections) {
    if rand.Float32() < 0.10 {
        dur = s.cfg.DBQueryDuration * 25
    }
}
```

## Context bounds the test lifetime

Run uses context.WithTimeout for the configured duration and passes that context to each HTTP request. A five-second client timeout also applies. The server checks cancellation while acquiring a slot and waiting for simulated work; it does not return a dedicated overload response.

## Canceled attempts are not all counted

Request-construction and client.Do errors are counted only while ctx.Err() is nil. End-of-test canceled attempts can be omitted. TotalRequests is recorded success samples plus counted errors, not every attempt started.

```go
if err != nil {
    if ctx.Err() == nil {
        errs++
    }
    continue
}
```

## Success follows an HTTP status rule

The runner drains and closes response bodies. Status >=400 increments errors; status <400 records success latency, which is broader than a 2xx-only rule. JSON business fields are not validated, and body-read errors are ignored.

## Latency includes reading the body

Timing starts before client.Do and ends after the successful response body is drained and closed. It is not query-only time. DNS, dialing, TLS, first-byte wait, pool wait, and body transfer are not separately measured.

## Percentiles use successful samples only

Failed responses and transport errors have no latency samples. Percentiles, average, min, and max use successful samples only. With no successes, latency fields remain zero; that does not mean the service responded instantly.

## RPS includes counted errors

Total is len(latencies)+errors. RPS divides that total by Run's actual elapsed duration after goroutines finish, so it is not success-only throughput. With no total requests, an empty Result is returned, including zero Duration.

```go
total := len(latencies) + errors
res.RPS = float64(total) /
    totalDuration.Seconds()
```

## Percentiles select a floored index

The calculator copies and sorts samples, then selects int((n-1)*p/100) without interpolation. For the test fixture 1–100 ms, P50=50 ms, P95=95 ms, and P99=99 ms. Other tools may use different percentile conventions.

```go
idx := int(float64(len(sorted)-1) *
    (pct / 100.0))
return sorted[idx]
```

## A fixture makes the calculation explicit

The 1–100 ms fixture has min 1 ms, max 100 ms, and average 50.5 ms. time.Duration retains nanosecond units, so division does not force whole-millisecond rounding. These are calculated fixture values, not HTTP performance measurements.

## Small samples make the tail coarse

With one 42 ms sample, every percentile is 42 ms. Small datasets can map several tail percentiles to the same sample and may exclude the maximum under this formula. Sample count and method matter when interpreting precision.

## Average and tail answer different questions

Average summarizes all successful samples; P95 selects a high-ranked sample. Unit tests verify ordered percentile invariants. The integration test separately checks stress P95 exceeds stress average; that relationship is not guaranteed for every possible distribution.

## Demo numbers are recorded examples

The repository records a sample smoke P95 of about 21.37 ms and stress P95 of about 1.357 s, with zero counted errors. These are historical example results, not a fresh measurement. Random spikes and scheduling vary; omitted canceled attempts also qualify the zero-error interpretation.

## Transport settings are not a full generator audit

MaxIdleConns and MaxIdleConnsPerHost are both 1000; these are idle-connection retention settings, not proof that every client bottleneck is eliminated. All successful latencies are retained per VU, merged, copied, and sorted, so memory grows with collected samples.

## A booking response is not a real transaction

Non-POST methods receive 405. Completed POST requests return 201, confirmed, and the fixed ID BK-1001. Request JSON is not decoded or validated, and there is no persistence, authentication, idempotency, or unique booking creation.

## Tests establish specific behaviors

Five metric tests cover fixtures, empty input, one sample, RPS, and ordering. Seven integration tests cover load comparison, HTTP 500, method rejection, dial errors, cancellation, slot limits, and total-count invariants. They do not prove queue fairness or maximum capacity.

## The lab bottleneck is known by design

The slot bottleneck is deliberately constructed. In production, rising P95 alone does not isolate application queues, database work, network, or external dependencies. The runtime has no exporter, tracing, or automated SLA threshold; the demo prints summaries.

## Rerun and interpret in context

The module declares Go 1.22 and uses the standard library. Go is unavailable in this authoring environment, so source and assertions were inspected without rerunning tests. Recorded repository results are prior-run evidence. Load testing reveals how latency changes when demand competes for limited resources.

```sh
cd labs/15-load-testing
go test -v ./...
go test -race ./...
go run ./cmd/demo
```

## Source

[Software Engineering Lab 15](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/15-load-testing), snapshot commit `82af988b89439fa92b1200406c5afbd59ce72a56`. Source: README, runner, metric calculator and tests, server, integration tests, demo, execution record, and go.mod. Code excerpts belong in their original package context.
