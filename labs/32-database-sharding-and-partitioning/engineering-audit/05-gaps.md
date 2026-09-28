# Gap Analysis

Target Lab: `labs/32-database-sharding-and-partitioning`

## Identified Gaps

No blocking gaps or critical vulnerabilities found.

### Non-Blocking Observations / Minor Notes

1. `ExtractTimeFromUUIDv7` helper function in `internal/idgen/idgen.go`:
   - Severity: LOW
   - Gap Type: MISSING_TEST
   - Description: `ExtractTimeFromUUIDv7` is present in `idgen.go` as a utility helper but is not explicitly asserted in `TestIDGenerators`. The primary `NewUUIDv7()` time-ordering property is thoroughly tested, however.
   - Impact: Non-blocking.

2. Fixed Virtual Node Count in Default Tests vs Demo:
   - Severity: LOW
   - Gap Type: IMPLEMENTATION_OVERCLAIM (Minor)
   - Description: `TestRoutingAndConsistentHashRelocation` uses 150 vnodes, while `cmd/demo` uses 100 vnodes. Both achieve valid relocation ratios within expected bounds (12% to 16% on 4->5 node resize vs minimal 20%).
   - Impact: Non-blocking.

## Summary

- Critical Gaps: 0
- High Gaps: 0
- Medium Gaps: 0
- Low Gaps: 2 (Non-blocking minor observation / test coverage edge)
