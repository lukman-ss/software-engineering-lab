# Docs vs Code

Sources:
- README.md (lab root)
- engineering/01-design.md, 02-implementation-notes.md, 03-execution-result.md

Claims in README.md:
1. Implemented Constraints list:
   1. NOT NULL (23502) — ✓ code/test match.
   2. CHECK (23514) — ✓.
   3. UNIQUE (23505) — ✓.
   4. FOREIGN KEY (23503) — ✓.
   5. PARTIAL UNIQUE INDEX (conditional uniqueness) — ✓.
2. Running Tests: go test -v ./... → works.
3. Running with race: go test -race ./... → works.
4. Demo: go run ./cmd/demo → output matches claimed behavior (see engineering/03-execution-result.md).

Design doc (01-design.md):
- Concept To Prove #1–4 matches implemented constraints.
- Expected Behavior lines 13–18 match code/test.
- Failure Scenario lines 21–22 matches UnsafeStore race test claim.
- Success Criteria lines 25–29: test coverage, concurrency test, error classification, test pass, demo run → all verified.
- Architecture lines 31–37: package map:
   - internal/db → does not exist (actual: internal/engine, internal/model, internal/store, internal/dberr). MISMATCH but non-behavioral.
   - internal/errors → actual internal/dberr → naming difference only.
   - internal/service → does not exist (not implemented). No claim in code uses it.
   - cmd/demo → matches.
   Note: These are internal package organization mismatches; they do not affect external behavior claims.
- Implementation Decisions lines 66–68: in-memory pure Go engine matches code (no external DB).

Execution result (03-execution-result.md):
- Build: claims PASS → verified.
- Tests: claims PASS + listed test names → verified.
- Race Detector: claims PASS → verified.
- Demo: verbatim block matches actual demo run (observed earlier).
- Final status: READY_FOR_ENGINEERING_AUDIT — supported.

No demo/fake benchmark detected.

Primary doc/code mismatches (non-blocking):
1. README mentions "schema invariants" — vague but satisfied.
2. engineering/01-design.md mentions packages internal/db and internal/service that are not present; but no behavior is tied to them.
3. No claim about performance or latency.

Thus: DOC_CODE_MISMATCH only on internal package names (low severity); all behavioral claims match.