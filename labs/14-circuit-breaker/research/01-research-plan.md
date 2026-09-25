# 01 Research Plan

## Research Topic
Circuit Breaker Pattern — Preventing a single degraded dependency from taking down the entire system.

## Objective
Design, verify, and document an in-memory Circuit Breaker in Go. Compare with Timeouts, Retries, Bulkheads, Fallbacks, and Load Shedding.

## Scope
1. State machine: CLOSED, OPEN, HALF_OPEN.
2. Cascade failure prevention via fail-fast execution.
3. Observability and metrics.
4. Real-world architecture case studies (CMMS WhatsApp, PPOB Payment/Pulsa).
