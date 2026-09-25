# 07 Fallback & Bulkhead

## Fallback
When circuit is OPEN, providing a degraded response (cached values, empty lists, queued tasks) rather than hard failing to user.

## Bulkhead
Isolating resource pools per dependency (e.g., ThreadPool A for Service 1; ThreadPool B for Service 2). If Service 1 hangs, ThreadPool A is exhausted, but Service 2 operates fine.

## Load Shedding
Server-side protection. Dropping incoming requests immediately with 503/429 when CPU/memory runs hot. Circuit breaker is caller-side protection; Load shedding is target-side.
