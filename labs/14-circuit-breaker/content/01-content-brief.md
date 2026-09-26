# Content Brief

**Topic:** Circuit Breaker Pattern — Implementation and Demonstration in Go

**Target Reader:** Backend engineers, SREs, and software architects implementing resilience patterns in distributed systems.

**Problem:** Remote calls across networks fail or hang. Without protection, callers block waiting for timeouts, holding threads, sockets, and memory until system resources deplete and cascading failures propagate through the system.

**Core Mental Model:** Circuit Breaker monitors downstream failures, trips open to block calls when a threshold is exceeded, fails fast in microseconds, then safely probes downstream before restoring traffic.

**Approved Research Status:** APPROVED

**Approved Engineering Status:** APPROVED

**Main Concepts:**
- Three-state machine: CLOSED, OPEN, HALF_OPEN
- Fail-fast behavior when OPEN
- Cooldown period with configurable timeout
- Probe calls during HALF_OPEN to test recovery
- Thread-safe state transitions using mutex
- Consecutive failure counting

**Verified Behaviors:**
- CLOSED: All requests route to downstream; failures increment counter; successes reset counter; transitions to OPEN when failures reach threshold
- OPEN: All requests immediately fail fast with `ErrCircuitOpen` without network calls; starts cooldown timer
- HALF_OPEN: Allows limited probe calls; success transitions to CLOSED; failure transitions back to OPEN
- Zero downstream requests executed when circuit is OPEN
- Concurrency-safe under race detector with no data races
- Panics handled without corrupting internal state
- Default configuration values applied when not specified

**Available Case Studies:**
- Scenario 1: Without Circuit Breaker — slow dependency causes blocking timeouts
- Scenario 2: With Circuit Breaker — fail-fast behavior after threshold exceeded
- Scenario 3: Recovery — HALF_OPEN to CLOSED transition on successful probe
- Scenario 4: Failed Recovery — HALF_OPEN to OPEN transition on failed probe

**Warnings:**
- Demo timing values (100ms HTTP timeout, 300ms cooldown) are illustrative for fast testing; production must tune to actual SLA and recovery profiles
- Implementation uses consecutive failure counting only, not sliding window or error rate calculation
- Circuit breaker does not heal a broken dependency