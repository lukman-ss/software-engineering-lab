## Section: Problem

**Research**: `research/05-report.md` — "Finding 1: The Dual-Write Problem and Failed Transactions"

**Implementation**: `internal/outbox/service.go:55-90` — `CreateOrderDualWriteNaive`

**Tests**: `tests/outbox_test.go:117-139` — `TestDualWriteProblem_Failure`

---

## Section: Why This Matters

**Research**: `research/05-report.md` — Executive Summary, Conclusion

---

## Section: Mental Model

**Research**: `research/05-report.md` — "Finding 2: Outbox Pattern Guarantees Atomicity"

**Implementation**: `internal/outbox/db.go:99-116` — `Commit()`

---

## Section: Core Concept

### Atomicity in One Transaction

**Research**: `research/05-report.md` — "Finding 2"

**Implementation**: `internal/outbox/service.go:18-53` — `CreateOrderWithOutbox`, `internal/outbox/db.go:118-127` — `Rollback()`

### Message Relay (Polling Publisher)

**Research**: `research/05-report.md` — "Finding 3: Message Relay Implementation Alternatives"

**Implementation**: `internal/outbox/relay.go:24-59` — `Start()`, `PollAndDispatch()`

### Idempotent Consumer

**Research**: `research/05-report.md` — "Finding 4: Idempotent Consumer Requirement"

**Implementation**: `internal/outbox/consumer.go:19-31` — `Handle()`

---

## Section: Architecture

**Engineering**: `engineering/01-design.md` — Architecture diagram, Components table

---

## Section: Implementation

**Engineering**: `engineering/02-implementation-notes.md` — Files Added, Core Design Decisions

---

## Section: Code Walkthrough

**Engineering**: `engineering/03-execution-result.md` — Demo execution output

**Implementation**: `cmd/demo/main.go` — Full demo executable

---

## Section: What the Tests Prove

**Research**: `research/05-report.md` — Findings 1-5

**Engineering**: `engineering-audit/03-test-audit.md` — Test coverage analysis

**Tests**: `tests/outbox_test.go` — 5 test functions

---

## Section: Recovery / Rollback

**Research**: `research/05-report.md` — "Finding 6: Operational Requirements"

**Implementation**: `internal/outbox/relay.go:43-59` — Retry on publish failure

---

## Section: Common Mistakes

**Research**: `research/06-open-questions.md` — Open questions about CDC, schema evolution

---

## Section: Case Study: Demo End-to-End

**Engineering**: `engineering/03-execution-result.md` — Demo command and output

**Implementation**: `cmd/demo/main.go` — Demo executable

---

## Section: Checklist

**Engineering**: `engineering-audit/06-verdict.md` — Quality gates (APPROVED)

---

## Section: Key Takeaways

**Research**: `research/05-report.md` — Conclusion, Areas of Agreement

---

## Section: Sources

**Research**: `research/01-plan.md`, `research/02-sources.md`, `research/03-evidence.md`, `research/04-contradictions.md`, `research/05-report.md`, `research/06-open-questions.md`

**Research Audit**: `research-audit/01-audit-plan.md`, `research-audit/02-source-audit.md`, `research-audit/03-claim-audit.md`, `research-audit/04-contradictions.md`, `research-audit/05-code-audit.md`, `research-audit/06-gaps.md`, `research-audit/07-verdict.md`

**Engineering**: `engineering/01-design.md`, `engineering/02-implementation-notes.md`, `engineering/03-execution-result.md`

**Engineering Audit**: `engineering-audit/01-audit-plan.md`, `engineering-audit/02-code-audit.md`, `engineering-audit/03-test-audit.md`, `engineering-audit/04-docs-vs-code.md`, `engineering-audit/05-gaps.md`, `engineering-audit/06-verdict.md`

**Source Files**:
- `internal/outbox/model.go`
- `internal/outbox/db.go`
- `internal/outbox/broker.go`
- `internal/outbox/service.go`
- `internal/outbox/relay.go`
- `internal/outbox/consumer.go`

**Tests**:
- `tests/outbox_test.go`

**Documentation**:
- `README.md`
