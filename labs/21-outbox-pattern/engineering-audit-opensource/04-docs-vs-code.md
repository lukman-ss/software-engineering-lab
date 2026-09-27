# Documentation vs Code Audit

## README.md
Matches implementation accurately:
- internal/outbox/db.go = in-memory transactional db (README says "in-memory transactional database" OK)
- broker.go = thread-safe mock broker OK
- service.go = dual-write vs atomic outbox OK
- relay.go = async polling relay OK
- consumer.go = idempotent consumer OK
- Run commands: `go test ./...`, `go test -race ./...`, `go run ./cmd/demo` — all executed OK

## engineering/01-design.md (RESEARCH MISMATCH)
Design doc CLAIMS SQLite implementation:
- Line 10: single SQL transaction inserts to `orders` & `outbox_events` tables
- Line 50: SQLite DB components
- Line 69: `modernc.org/sqlite` / `database/sql`
ACTUAL implementation: in-memory map-based Tx in db.go — no SQL, no SQLite, no CGO.
Status: RESEARCH_IMPLEMENTATION_MISMATCH (Design doc / research describes SQLite, code is in-memory mock). Not fatal: implementation notes explicitly choose in-memory for portability. But design doc is stale vs chosen implementation.

## engineering/02-implementation-notes.md
Accurately describes the chosen in-memory approach and known limitations (no CDC/log-tailing, polling publisher). Matches code. OK.

## engineering/03-execution-result.md
Matches re-executed results (tests pass, race passes, demo output matches). OK.

## README vs DESIGN
README: in-memory DB. Design doc: SQLite. Conflict between documentation layers. README/code agree; design doc out of sync.

## Summary of mismatches
1. DOC_CODE_MISMATCH: design.md describes SQLite tables, actual code is map-based in-memory DB.
2. README_CODE: aligned.
3. TEST_CLAIM_MISMATCH: Test strategy (design.md) claims "Broker failure retry mechanism" tested; no test covers relay retry after transient broker failure (SetFailNext not used in relay test path).
4. ARCH_MISMATCH: design.md architecture implies Broker -> Consumer delivery; code has no Broker->Consumer wire (consumer reads published list manually).
