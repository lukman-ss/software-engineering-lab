# 10 Final Research Summary

## Synthesis
1. Circuit Breakers isolate network blast radius.
2. Caller fails in microseconds instead of multi-second HTTP timeout hangs.
3. Asynchronous non-critical flows (e.g. notifications) must be decoupled using queues + idempotency so third-party downtime never aborts core transaction persistence.
