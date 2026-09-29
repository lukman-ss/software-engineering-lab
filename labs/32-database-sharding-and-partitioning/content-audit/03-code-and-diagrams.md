# Code & Diagram Consistency Audit

## 1. Code Snippets Consistency
- `content/03-code-snippets.md` contains 4 code snippets:
  1. `Snippet 1`: `internal/partitioning/table.go` (Partition Pruning, DropPartition)
  2. `Snippet 2`: `internal/sharding/sharding.go` (ModuloRouter & ConsistentHashRouter with virtual nodes)
  3. `Snippet 3`: `internal/sharding/sharding.go` (Scatter-Gather with Context & GlobalSecondaryIndex)
  4. `Snippet 4`: `internal/idgen/idgen.go` (RFC 9562 UUIDv7 & SequenceBlockAllocator)
- **Check**: All snippets accurately reproduce types, functions, and logic implemented in `internal/` packages. No mock methods or broken syntax.

## 2. Diagrams Consistency
- `content/04-diagrams.md` provides 4 ASCII diagrams:
  1. Architecture comparison: Single Instance Pruning vs Sharded Cluster Routing.
  2. Key migration visualization: Modulo 80% data shuffle vs Consistent Hash minimal relocation.
  3. Query resolution paths: Parallel scatter-gather vs GSI direct point-lookup.
  4. Bit layout of RFC 9562 UUIDv7 and Vitess sequence block flow.
- **Check**: All diagrams accurately mirror mathematical properties and operational flow described in the draft and codebase.
