# Engineering Gaps Analysis

Target Lab: `labs/37-cache-invalidation-strategies`

## Identified Gaps

### None (Zero Blocking or High Severity Gaps)

All requirements and claims from the approved research and engineering designs have been implemented, tested, and verified.

## Observations / Minor Limitations (Non-Blocking)

1. **In-Process SingleFlight Scope**
   - **Type**: Scope Limitation (Documented)
   - **Severity**: LOW
   - **Notes**: `singleflight.Group` operates strictly within a single Go process. Multi-node distributed caching would require external locks (e.g., Redis Redlock). This scope boundary is clearly noted in `engineering/02-implementation-notes.md`.

2. **Write-Behind Queue Overflow Dropping**
   - **Type**: Scope Limitation (Documented)
   - **Severity**: LOW
   - **Notes**: `WriteBehindService.Update` drops items when channel buffer is full (`default: ` branch). Clearly documented as a demonstration trade-off.

3. **In-Memory Volatility**
   - **Type**: Scope Limitation (Documented)
   - **Severity**: LOW
   - **Notes**: `MemoryCache` does not persist state across process restarts. Appropriate for standalone executable lab environment.
