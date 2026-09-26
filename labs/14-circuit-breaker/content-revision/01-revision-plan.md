# Content Revision Plan — Lab 14 Circuit Breaker

**Date:** 2026-09-26  
**Based on:** content-audit/03-full-audit.md, content-audit/02-discrepancies.md

## Files Requiring Revision

| File | Issues to Fix |
|------|---------------|
| content/02-master-draft.md | HIGH: Struct name `CircuitBreaker` → `Breaker`; State type `string` → `int` with `iota`; Field names (`failureCount`→`failures`, `halfOpenCalls`→`halfOpenIn`, `lastStateChange`→`openedAt`+`generation`); Method `checkStateTransitionLocked` → `advanceLocked`; HALF_OPEN logic (remove `consecutiveSuccesses`); Constructor `checkout.NewService` → `checkout.New`; Remove non-existent research file references |
| content/03-code-snippets.md | HIGH: Snippet 1 state type and constants; Snippet 2 default OpenTimeout 5s → 300ms, return `CircuitBreaker{}` → `Breaker{}`; Snippet 3 method name and field names; Snippet 4 full method rewrite to match actual implementation |
| content/05-key-takeaways.md | HIGH: OpenTimeout default 5s → 300ms; Duplicate numbering (two #6s) |
| content/06-source-map.md | LOW: Remove references to non-existent `research/03-evidence.md`, `research/06-open-questions.md`; Fix `checkStateTransitionLocked` → `advanceLocked`; Fix `CircuitBreaker` → `Breaker` |

## Revision Strategy
Direct content edits to match actual implementation in `internal/circuitbreaker/circuit_breaker.go` and `internal/checkout/service.go`.

---

# Content Revision Changes Made — Lab 14 Circuit Breaker

**Date:** 2026-09-26  
**Revisor:** Technical Content Reviser

## Changes Applied

### content/02-master-draft.md
- **Struct dump (lines 22-33)**: Replaced `CircuitBreaker` struct with actual `Breaker` struct matching code (`mu`, `cfg`, `state`, `failures`, `openedAt`, `halfOpenIn`, `generation`); removed `now func() time.Time`, `lastStateChange`, `consecutiveSuccesses`
- **State type (lines 35-37)**: Added correct `type State int` with `iota` constants (`Closed`, `Open`, `HalfOpen`) and note about exported aliases
- **State transitions (lines 39-43)**: Updated to use `failures`, `openedAt`, `advanceLocked(now)`, and HALF_OPEN→CLOSED on first successful probe
- **Execution flow (lines 45-71)**: Rewrote to match actual `Execute` method with `b *Breaker`, `advanceLocked`, `b.halfOpenIn`, `b.cfg.HalfOpenMaxCalls`, generation tracking
- **Checkout constructor (line 84)**: Changed `checkout.NewService` → `checkout.New` (primary constructor)
- **Sources section**: Removed references to non-existent `research/03-evidence.md`

### content/03-code-snippets.md
- **Snippet 1 (lines 7-22)**: Changed `type State string` → `type State int` with `iota`; constants `Closed`, `Open`, `HalfOpen`; added exported aliases block (`StateClosed = Closed`, etc.)
- **Snippet 2 (lines 30-58)**: Fixed `OpenTimeout` default from `5 * time.Second` → `300 * time.Millisecond`; return type `*CircuitBreaker` → `*Breaker`; struct init `state: StateClosed, config: cfg, lastStateChange: time.Now(), now: time.Now` → `cfg: cfg, state: Closed`; added `DefaultConfig()` function
- **Snippet 3 (lines 60-74)**: Renamed `checkStateTransitionLocked()` → `advanceLocked(now time.Time)`; receiver `*CircuitBreaker` → `*Breaker`; uses `b.openedAt`, `b.cfg.OpenTimeout`, `b.halfOpenIn`, `b.generation`; removed `consecutiveSuccesses` reset
- **Snippet 4 (lines 76-118)**: Complete rewrite to match actual `Execute` method: single deferred panic handler with `panicked` flag, generation tracking, `onSuccessLocked`/`onFailureLocked` calls; uses `b.halfOpenIn`, `b.cfg`, `b.failures`; HALF_OPEN→CLOSED on first success
- **Snippet 5 (lines 120-132)**: Simplified to actual `Checkout` method using `s.breaker` and `s.payment`

### content/05-key-takeaways.md
- **Takeaway 6**: Changed "Consecutive successes required" → "HALF_OPEN transitions to CLOSED on first successful probe"
- **Takeaway 7**: Fixed numbering (was duplicate #6)
- **Takeaway 8**: Changed `OpenTimeout=5s` → `OpenTimeout=300ms`
- **Takeaway 9-11**: Renumbered to 9, 10, 11

### content/06-source-map.md
- **State Machine section**: `CircuitBreaker struct` → `Breaker struct`; `checkStateTransitionLocked` → `advanceLocked`
- **Fail-Fast section**: Removed `research/03-evidence.md` reference
- **Recovery section**: `checkStateTransitionLocked` → `advanceLocked` (line numbers updated); `CircuitBreaker` → `Breaker`
- **Concurrency Safety**: `CircuitBreaker struct` → `Breaker struct`; `cb.mu.Lock()` → `b.mu.Lock()`
- **Full Source Inventory**: Removed `research/03-evidence.md` and `research/06-open-questions.md` from Research Files list

## Verification
All changes align with:
- `internal/circuitbreaker/circuit_breaker.go` (actual implementation)
- `internal/checkout/service.go` (actual checkout service)
- `cmd/demo/main.go` (demo config uses 300ms)
- `engineering/03-execution-result.md` (verified demo output)
- Approved research sources

---

# Content Revision Result — Lab 14 Circuit Breaker

**Date:** 2026-09-26  
**Status:** All audit issues addressed

## Summary
Content files now accurately reflect the actual implementation. The major structural mismatches in `02-master-draft.md` and `03-code-snippets.md` (incorrect type names, wrong defaults, non-existent fields/methods) have been corrected. `05-key-takeaways.md` has correct default value and numbering. `06-source-map.md` no longer references non-existent research files.

Expected audit outcome after re-audit: **APPROVED**