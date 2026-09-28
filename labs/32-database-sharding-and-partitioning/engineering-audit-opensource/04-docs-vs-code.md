# Docs vs Code Analysis

## DOC_CODE_MISMATCH: ExtractTimeFromUUIDv7
- README & engineering notes reference UUIDv7 timestamp extraction. Code function exists but returns zero time on success.

## RESEARCH_IMPLEMENTATION_MISMATCH: Consistent Hash Migration Percent
- Design expects ~1/N movement. Observed live demo ~12–16%; tests assert 5–40%. Within bounds; no mismatch.

## TEST_CLAIM_MISMATCH: Edge Cases
- README/tests do not cover empty router or out-of-range partition behavior claimed by error types.

## DOC_CODE_MISMATCH: GSI Mapping
- Docs describe GSI mapping non-shard-key → shardKey; code matches. OK.

## DOC_CODE_MISMATCH: Concurrency Model
- Docs state thread-safe shards and scatter-gather goroutines; code implements RWMutex/goroutines consistently. OK.

## Documentation Accuracy: WARNING (ExtractTimeFromUUIDv7 contradiction)
