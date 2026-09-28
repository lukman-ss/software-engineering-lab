# Docs vs Code

## README vs Code
Claimed component internal/saga/orchestrator.go exists. -> TRUE
Claimed internal/saga/choreography.go exists. -> TRUE
Claimed internal/services/services.go exists. -> TRUE
Claimed cmd/demo/main.go exists. -> TRUE
Claimed tests/saga_test.go exists. -> TRUE

## Engineering Design vs Code
Design states pkg/saga and pkg/services; actual code is in internal/. Doc mismatch (path).
Design states steps Order -> Payment -> Inventory -> Delivery. Code shows Order -> Payment -> Inventory -> ApproveOrder (no Delivery step). Mismatch.

## Implementation Notes vs Code
Notes claim LIFO compensation; confirmed in code (reverse iteration).
Notes claim idempotency payment operations check transaction keys; confirmed.
Notes claim semantic lock countermeasure applied on order creation; confirmed.
Known limitations (no persistent transaction logs); consistent.

## Execution Result vs Actual
All tests pass when run (PASS).
Race detector passes (ok labs/29-saga-pattern/tests).
Demo output matches recorded output.

## Gap Classification
DOC_CODE_MISMATCH: design docs reference pkg/ paths; repo uses internal/.
RESEARCH_IMPLEMENTATION_MISMATCH: design mentions Delivery step absent in code.