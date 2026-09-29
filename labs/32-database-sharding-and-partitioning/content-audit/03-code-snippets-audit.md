# Audit: 03-code-snippets.md

Status: NEEDS_REVISION
Issues: 1 blocking

## Blocking Issue: Snippet 4 Line 370 — Wrong Bit-Shift Value

**Location**: `content/03-code-snippets.md` line 370
**Content claim**: `uuid[3] = byte(ms >> 36) // Note: bit shift byte assembly`
**Actual code** (`internal/idgen/idgen.go:19`): `uuid[3] = byte(ms >> 16)`

The content snippet incorrectly reports the bit-shift value as `36` instead of `16`. This is a factual inaccuracy — the snippet must match the actual source code.

**Secondary note**: Neither the content nor the code implements the correct RFC 9562 bit layout. The RFC-compliant value should be `byte(ms >> 24)`. The engineering team should fix the code; the content should match the fixed code.

## Other Snippets — All Accurate

| Snippet | Source File | Status |
|---|---|---|
| Snippet 1 (Partition Pruning) | `internal/partitioning/table.go` | ACCURATE |
| Snippet 2 (Modulo + Consistent Hash) | `internal/sharding/sharding.go` | ACCURATE |
| Snippet 3 (Scatter-Gather + GSI) | `internal/sharding/sharding.go` | ACCURATE |
| Snippet 4 (UUIDv7 + Sequence) — lines 363-369, 371-389 | `internal/idgen/idgen.go` | PARTIALLY ACCURATE (see blocking issue) |

## Required Action
Fix line 370 in `content/03-code-snippets.md` to match `idgen.go:19` (`byte(ms >> 16)`), and ideally also fix the underlying code to use `byte(ms >> 24)` for RFC 9562 compliance.
