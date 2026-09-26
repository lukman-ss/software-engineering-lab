# 07 Fallback & Bulkhead

## Fallback
When circuit is OPEN, provide degraded response instead of hard failure.

### Safe Fallbacks
- **Cached/stale read response**: return previously cached data with expiry indicator (e.g. "data may be up to N minutes stale").
- **Degraded static response**: default values, empty list, placeholder banner.
- **Queued processing**: defer non-critical write to durable queue; return "request received, will process async" (PPOB WhatsApp case).

### Unsafe Fallbacks (AVOID)
- Defaulting a financial transaction amount to 0 or a cached balance to succeed a payment write.
- Silencing an error in a critical path that masks permanent data inconsistency.

**Source**: Azure Circuit Breaker — "rather than returning a failure and raising an exception, the Open state can return a default value that's meaningful to the application." Microsoft explicitly warns fallback must not mask errors that should propagate. — Confidence HIGH.

## Bulkhead
- **Definition**: Isolate resources so failure in one dependency cannot consume resources needed by others.
- **Mechanism**: Separate thread pools / connection pools / memory quotas per dependency.
- **Example**: ThreadPool A for Service X; ThreadPool B for Service Y. X exhausts its pool → only X's calls fail; Y unaffected.

**Sources**: Azure Bulkhead pattern (2026-03-19) — "isolate resources for specific dependencies so that a disruption in one service doesn't affect the entire application"; Hystrix Wiki Isolation — "isolating dependencies from each other and limiting concurrent access to any one of them."

## Bulkhead vs Circuit Breaker — Distinct
| | Bulkhead | Circuit Breaker |
|---|---|---|
| **Controls** | Resource consumption concurrency | Request pass/fail decision |
| **Goal** | Prevent pool exhaustion propagation | Prevent repeated known-fail calls |
| **Triggers on** | Resource quota depletion | Error-rate threshold breach |
| **State** | Count of available slots | State machine (3 states) |

**Corroborated**: Azure Bulkhead — "Circuit breakers, throttling — combine [bulkhead] with retry, circuit breaker, and throttling patterns"; Azure Circuit Breaker — separate pattern pages.

## Load Shedding vs Circuit Breaker
- **Circuit Breaker**: Caller-side. Stops *outbound calls* to a known-failing dependency.
- **Load Shedding**: Target/Sever-side. Drops *inbound requests* when CPU/memory saturated (503/429).

**Source**: Google SRE Managing Load (Ch 11, 2018) — load shedding is server-side drop; Dressy case study where load shedding + load balancing misconfigured in isolation caused wrong routing. — https://sre.google/workbook/managing-load/ — Tier 1.

## Queue-Based Load Leveling (Async Fallback)
- Decouples caller from synchronous dependency via durable queue.
- Enables degraded async processing when sync dependency is down.
- Requires: idempotency, dead-letter for poison, durability.

**Source**: Azure Queue-Based Load Leveling (2026-06-09) — "application can continue to post messages to the queue even when the service isn't available." — NOT a circuit breaker; complements by converting sync-fail into async-buffer.

## PPOB / CMMS Guidance
- **Critical write path** (Order → Provider): must persist Order locally before attempting provider call. If provider down, queue the outbound integration message (idempotent) — never silently drop or fake success.
- **Non-critical notification** (WhatsApp send): can degrade to queue + retry without blocking user flow. Circuit breaker around integration client allows fail-fast on persistent failures.
- **Unsafe**: falling back a payment charge to "cached payment success."

**Sources**: Azure Queue-Based Load Leveling idempotency guidance (2026-06-09); Azure Retry idempotency (2024-07-18).