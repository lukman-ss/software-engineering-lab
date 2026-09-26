# Audit: 03-code-snippets.md

## Scope
346 lines, 12 Go code snippets referencing source files with line ranges.

## Accuracy Verification (vs source code)

| Snippet | Claimed Source | Actual Source | Match |
|---|---|---|---|
| 1 CreateOrderWithOutbox | service.go:18-53 | service.go 18-53 | YES |
| 2 Commit | db.go:99-116 | db.go 99-116 | YES |
| 3 Rollback | db.go:118-127 | db.go 118-127 | YES |
| 4 PollAndDispatch | relay.go:43-59 | relay.go 43-59 | YES |
| 5 Handle | consumer.go:19-31 | consumer.go 19-31 | YES |
| 6 Rollback Test | tests 61-90 | tests 61-90 | YES |
| 7 Idempotency Test | tests 92-115 | tests 92-115 | YES |
| 8 DualWrite Failure Test | tests 117-139 | tests 117-139 | YES |
| 9 CreateOrderDualWriteNaive | service.go:55-90 | service.go 55-90 | YES |
| 10 Relay Start | relay.go:24-37 | relay.go 24-37 | YES |
| 11 OutboxMessage | model.go:26-32 | model.go 26-32 | YES |
| 12 Order | model.go:12-17 | model.go 12-17 | YES |

## Code Content Verification
- All 12 snippets copy source code exactly byte-for-byte (verified against file reads)
- Comments preserved identically
- Struct field types preserved: float64, string, time.Time, MessageStatus, OrderStatus
- Function signatures match source

## Explanations
- Purpose fields all accurate
- Explanation paragraphs correct for each snippet's behavior
- No hallucinated code or invented APIs

## Hallucination / Bias Check
- Purely sourced from implementation, no platform bias
- No invented metrics or thresholds

## Formatting & Clarity
- Consistent two-part format: Source File + Purpose + code + Explanation
- Clear labeling throughout
- No formatting errors

## Issues
- None
