# Audit: 04-diagrams.md

## Scope
179 lines, 8 mermaid diagrams with source attributions.

## Diagram Verification

### Diagram 1 — Transactional Outbox Architecture
- Mermaid graph TB: correct flow
- Components: Service -> DB -> Relay -> Broker -> Consumer
- Steps 1-9 match actual sequence ✓
- Sources listed: internal/outbox/*.go files all valid ✓
- PASS

### Diagram 2 — Transaction Flow (Atomic Write)
- SequenceDiagram participants: App, Tx, DB
- Steps: BeginTx -> SaveOrder -> SaveOutbox -> Commit -> Lock -> Write stagedOrders -> Write stagedOutbox -> Unlock -> OK -> Success
- All steps map to actual code: Tx.Commit 99-116, Service 18-53 ✓
- Source attribution correct ✓
- PASS

### Diagram 3 — Rollback Flow
- Steps: BeginTx -> SaveOrder -> SaveOutbox -> Rollback -> closed=true -> Discard staged -> Error -> No writes
- Maps to db.Rollback 118-127, service 28-32 marshal error path ✓
- Source attribution correct ✓
- PASS

### Diagram 4 — Relay Polling Cycle
- Loop every pollInterval: GetPendingOutbox -> for each publish -> Success/OK/MarkProcessed OR Failure/Status remains PENDING
- Maps to relay.go 43-59, 24-37 ✓
- Correctly models retry-on-failure behavior ✓
- PASS

### Diagram 5 — Dual-Write Failure
- Steps: SaveOrder -> Commit -> OK -> Publish -> ERROR -> STATE INCONSISTENT
- Maps to service.go 55-90, tests 117-139 ✓
- Correctly shows order persisted, event lost ✓
- PASS

### Diagram 6 — Idempotent Consumer Handling Duplicates
- Steps: Handle evt-1 -> Check not found -> Add & return true -> (crash) -> Handle evt-1 retry -> Check found -> Return false
- Maps to consumer.go 19-31, tests 92-115 ✓
- Correctly models crash-before-mark scenario ✓
- PASS

### Diagram 7 — Outbox Message State Machine
- [*] -> PENDING -> PROCESSED
- PENDING -> PENDING on relay publish failure
- PROCESSED -> [*] cleanup (not implemented)
- States match model.go 21-24 constants ✓
- Transitions match relay.go 43-59 logic ✓
- PASS

### Diagram 8 — Demo Scenarios Flow
- Flowchart: Dual-Write Failure -> Outbox Solution -> Idempotency
- Maps to cmd/demo/main.go 21-62 ✓
- Correctly shows broker failure, atomic create, duplicate rejection ✓
- PASS

## Hallucination / Bias Check
- No invented components, states, or transitions
- All diagrams derivable from source code
- No platform-specific assumptions
- PASS

## Formatting Issues
- All mermaid fences use triple backticks with "mermaid" tag
- No syntax errors detected (renderable format)
- PASS

## Issues
- None
