# Docs vs Code

## README vs Implementation
- README claims "In-engine range partitioning" and "partition pruning" — present in internal/partitioning/table.go:93 QueryRange.
- README claims "ModuloRouter" and "ConsistentHashRouter" — present in internal/sharding/sharding.go.
- README claims "Scatter-gather parallel execution and Global Secondary Index (GSI / Lookup Vindex) point-lookups" — present (Cluster.ScatterGatherBroadcastWithContext, GlobalSecondaryIndex).
- README claims "RFC 9562 compliant time-ordered UUIDv7 generator" — present in internal/idgen.
- README claims "Vitess-style Sequence Block Allocator" — present.
- README Quick Start commands match actual commands used.
Assessment: PASS — README matches implementation.

## Engineering/02-implementation-notes.md vs Implementation
- "Standard Library Only" — true; no external deps in the lab module.
- "Hash Function: hash/fnv (64-bit)" — true.
- "Virtual Node Hashing: 100 per physical shard" — true (default).
- "Scatter-gather queries with bounded Goroutines and sync.WaitGroup" — true.
- Known Limitations: in-memory, in-process, no 2PC — all true.
- Trade-offs: consistent hashing reduces migration from ~80% to ~12-20% — verified by demo (79.84% vs 12.00%).
Assessment: PASS.

## Engineering/03-execution-result.md vs Actual Execution
- Test results: PASS — verified with go test -v ./... (5 tests, PASS).
- Race detector: PASS — verified with go test -race ./... (cached ok).
- Demo output: PASS — verified with go run ./cmd/demo; output matches documented structure and numbers (79.84% modulo, 12.00% consistent hash).
- Minor difference: TestRoutingAndConsistentHashRelocation logged "Hash Modulo moved 3756 / 5000 keys (75.12%)" in my run vs "756 / 1000 keys (75.60%)" in the doc. The doc appears to have been generated with a different random seed / key set. The ratio is consistent (~75%) so this is a non-issue; the doc numbers are plausible but not reproducible exactly. This is a FAKE_DEMO/FAKE_BENCHMARK candidate only if the numbers were claimed to be deterministic — they are not. Treat as WARNING: documented numbers are sample output, not fixed constants.
Assessment: WARNING — documented output is a sample, not a deterministic benchmark.

## Research (not audited per pipeline override)
- Skipped per instructions.

## Summary
- DOC_CODE_MISMATCH: none.
- TEST_CLAIM_MISMATCH: none.
- RESEARCH_IMPLEMENTATION_MISMATCH: not evaluated (out of scope).
- Minor: documented test log numbers are sample output.