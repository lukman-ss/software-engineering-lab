# Gap Analysis

## GAP-01: BROKEN_IMPLEMENTATION
Location: internal/idgen/idgen.go:90-113 ExtractTimeFromUUIDv7
Severity: HIGH
Evidence: Success branch `return time.Time{}, nil` discards parsed timestamp; function never returns extracted time. Unused by tests/demo but public API broken.
Type: BROKEN_IMPLEMENTATION

## GAP-02: MISSING_TEST
Location: internal/idgen/idgen.go:90-113
Severity: MEDIUM
Evidence: TestIDGenerators covers NewUUIDv7 ordering + SequenceBlockAllocator but never calls ExtractTimeFromUUIDv7, so GAP-01 undetected.
Type: MISSING_TEST

## GAP-03: MISSING_EDGE_CASE
Location: internal/sharding/sharding.go GetShard empty-router path; internal/partitioning/table.go Insert out-of-range path
Severity: LOW
Evidence: Code returns ErrShardNotFound / ErrNoMatchingPartition; no negative test asserts these errors.
Type: MISSING_EDGE_CASE

## GAP-04: UNHANDLED_ERROR
Location: internal/sharding/sharding.go:433 RebalanceData `targetShard, _ := c.router.GetShard`
Severity: LOW
Evidence: Router error discarded; safe only when router non-empty. No test for empty-router rebalance.
Type: UNHANDLED_ERROR

No FAKE_DEMO / FAKE_BENCHMARK / RACE_CONDITION found. Demo output reproduced live; race detector clean.
