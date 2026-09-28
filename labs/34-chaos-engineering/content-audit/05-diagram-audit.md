# Diagram Audit

Source: `content/04-diagrams.md`
Reference: actual implementation logic and demo output.

## Diagram 1 — Circuit Breaker State Transition
**Status: ACCURATE**

Source logic (`internal/circuitbreaker/circuitbreaker.go`):
- CLOSED → failures >= threshold (or failure in Half-Open) → OPEN (lines 82-86)
- OPEN → cooldown elapsed → HALF-OPEN (lines 57-61, `checkStateLocked`)
- HALF-Open → failure → back to OPEN (line 84: `|| cb.state == StateHalfOpen`)
- HALF-Open → success → CLOSED (lines 95-98)
- CLOSED → success resets failure count (lines 99-101)

Diagram captures all transitions. No missing or incorrect transitions.

## Diagram 2 — Chaos Experiment Lifecycle & Auto-Abort
**Status: ACCURATE**

Source logic (`internal/experiment/runner.go`):
- PENDING → Run() → SetFault → RUNNING (lines 61-66)
- Loop with select on: ctx.Done() → ABORTED (75-77), timeout → COMPLETED (78-80), ticker.C → IsHealthy() → if unhealthy → terminate(ABORTED) + injector.Clear() (82-86)
- terminate() calls injector.Clear() synchronously (lines 91-97)

Diagram correctly shows: init, fault active, dual wait (duration vs monitor ticker), health check branch, abort path with Clear(), completion path with Clear(). No inaccuracies.

## Diagram 3 — Resilient Request Flow with Graceful Degradation
**Status: ACCURATE**

Traces `cmd/demo/main.go` `clientRequest` + `circuitbreaker.Execute`:
- CB.Execute: if state==Open & fallback exists → fallback() → success recorded (lines 68-71, 48)
- CB.Execute: else call downstream → if error → failures++ → if threshold hit → OPEN (82-86) → if fallback exists → fallback() → success recorded (88-89, 48); else return error → failure recorded (91, 45)
- CB.Execute: if success → if HalfOpen → CLOSED (95-98) else failures=0 (99-100)

Diagram flow matches. Minor: the "Failures >= Threshold? Yes → State = OPEN" decision box appears before "Fallback Provided?" in the diagram, but in code both happen in the same critical section — this is a presentation simplification, not an accuracy issue.

## Summary
All three diagrams accurately reflect the implementation and demo behavior. No hallucinated flows, no incorrect state transitions. Diagrams are pedagogical simplifications of correct logic.