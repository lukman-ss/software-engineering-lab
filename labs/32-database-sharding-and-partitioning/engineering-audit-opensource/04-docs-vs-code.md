# Docs vs Code Comparison

## README vs Implementation
- README claims logical partitioning with range pruning – code implements Table.QueryRange with overlap check; demo prints pruned partitions count – MATCH.
- README claims Modulo vs Consistent Hash routing comparison – code provides both routers, demo prints key remap percentages – MATCH.
- README claims Scatter‑Gather vs GSI query performance – code implements both, demo prints broadcast count and timings – MATCH.
- README claims UUIDv7 and Sequence Block Allocator ID generation – code provides both, tests verify ordering – MATCH.

## Test Claims vs Code
- Tests assert partition pruning, move ratios, GSI lookup, scatter‑gather behavior – all backed by implementations – MATCH.
- No test claims unsupported behavior.

## Research vs Code
- Pipeline override disables research verification for this stage.

## Mismatches
- None observed.
