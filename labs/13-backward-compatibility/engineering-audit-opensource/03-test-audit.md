# Test Audit

Coverage Matrix (claims vs tests):

| Claim / Behavior | Covered By Test | Result |
|---|---|---|
| Expand: dual-write to legacy + modern, additive payload | TestSerializationBackwardCompatibility, migration_test.TestFullExpandMigrateContractLifecycle | PASS |
| V1 legacy consumer parses `phone` string (backward compat) | TestSerializationBackwardCompatibility | PASS |
| V2 modern consumer parses `phones` array w/ is_primary | TestSerializationBackwardCompatibility | PASS |
| Migrate: backfill resumable + idempotent | TestBackfillIdempotentAndResumable | PASS |
| Migrate: fallback read + lazy backfill | TestFallbackRead | PASS |
| Migrate: drift reconciliation before/after | TestDataReconciliationAndDrift, TestFullExpandMigrateContractLifecycle | PASS |
| Read switch (ReadNewOnly) reads historical from modern | TestFullExpandMigrateContractLifecycle | PASS |
| Rollback safety during DualWrite (no data loss) | TestRollbackScenarios Scenario A | PASS |
| Premature rollback after NewOnly data loss (negative case) | TestRollbackScenarios Scenario B | PASS |
| Contract: deprecation headers on V1 | TestDeprecationHeadersAndContractEnforcement | PASS |
| Contract: guard blocks forced migration while traffic active | TestDeprecationHeadersAndContractEnforcement | PASS |
| Contract: 410 Gone after apply + V2 still works | TestDeprecationHeadersAndContractEnforcement, TestFullExpandMigrateContractLifecycle | PASS |
| Concurrency: parallel writes/reads/backfill/drift | TestConcurrency (under -race) | PASS |

Happy Path: PASS — Expand, Migrate, Contract lifecycle fully covered in migration_test.

Failure Path: PASS — drift detection, premature contract rejection (ErrContractViolation), post-contract 410, legacy-unavailable error propagation. Scenario B proves data loss when rollback occurs after NewOnly writes (negative case demonstrated, not merely asserted).

Edge Cases:
- Batch boundaries: covered (batchSize=3 over 10 records; asserts 3+3+4 distribution + completion flag).
- Empty final batch: covered (RunBatch sets IsComplete when ids empty).
- Idempotent rerun: covered (second RunAll == 0).
- Un-backfilled drift: covered (ReconcileData detects drift=1 before backfill, 0 after).

Concurrency: PASS — `go test -race ./...` exits cleanly across `TestConcurrency` (10 writers x 20, 10 readers x 50, backfill, drift) and all other tests. No data races, deadlocks, or `DualWriteErrors` spikes.

Negative Cases: PASS — contract guard rejection, empty legacy phone assertion, unmigrated-record fallback.

Gaps (missing / weakly asserted):
1. HTTP handler input validation: no test for invalid/missing `id` (Atoi error path) -> 404 vs 400. MISSING_TEST (LOW).
2. No unit test asserting `SavePhoneEntry` returns the pre-existing entry on duplicate (idempotency at storage method level). Covered only indirectly via backfill rerun. MISSING_TEST (LOW).
3. No test driving `WriteNewOnly` CreateUser end-to-end through the service; only `WriteDual`/`WriteLegacyOnly`/`CreateModern` direct path. MISSING_TEST (LOW).
4. `GetUserIDs` returns empty-slice + IsComplete behavior with zero users: covered indirectly via TestBackfillIdempotentAndResumable second-run. Acceptable.

Test Strength: PASS
A passing suite here is not a façade — assertions check exact counts (migrated==2, drifts==0/1, 2 phones parsed, is_primary flagged), not just "no error".
