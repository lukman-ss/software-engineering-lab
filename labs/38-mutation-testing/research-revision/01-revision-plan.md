# Revision Plan

Target Lab: labs/38-mutation-testing

Previous Audit Status: APPROVED_WITH_WARNINGS

## Blocking Issues

None.

## Non-Blocking Issues (from audit/07-verdict.md)

1. **Academic Primary Sources Unopened Directly**: DeMillo et al. (1978) and Jia & Harman (2009) surveyed solely via Wikipedia citations. Research explicitly annotates this limitation honestly without fabricating citations or quotations.

2. **Pre-publication Draft Citation**: Martin Fowler bliki entry is marked DRAFT. The research acknowledges this in sources and open questions, though it is cited as corroborating evidence.

3. **Single Industry Source for LLM Mutation Testing**: Meta ACH statistics (73% acceptance, 36% privacy relevance, 0.95/0.96 precision/recall) stem from a single primary industry blog post (Mark Harman, Sep 2025). The arXiv preprint was not directly read.

4. **Go Tooling Gap**: The lab target domain is Go, but tooling research covers JVM (PIT) and JS/TS (Stryker) without examining Go mutation packages (e.g. `go-mutesting`, `gremlins`).

## Research Gaps (from audit/06-gaps.md)

| Gap | Type | Severity | Issue |
|-----|------|----------|-------|
| 1 | WEAK_SOURCE | MEDIUM | DeMillo et al. (1978) not directly accessed |
| 2 | WEAK_SOURCE | MEDIUM | Jia & Harman (2009) not directly accessed |
| 3 | UNVERIFIED_CLAIM | MEDIUM | ACH arXiv preprint PDF not successfully read |
| 4 | OVERGENERALIZATION | MEDIUM | "Five barriers" presented as universal; actually Harman/Meta's framing |
| 5 | MISSING_SOURCE | LOW | No consensus on mutation score thresholds; need explicit note |
| 6 | WEAK_SOURCE | LOW | Martin Fowler DRAFT used as supporting evidence |
| 7 | IMPLEMENTATION_GAP | MEDIUM | No Go-specific mutation testing tooling researched |
| 8 | MISSING_CASE | LOW | Subsumed mutants not documented |

## Contradictions (from audit/04-contradictions.md)

- Contradiction A (MEDIUM): Fowler DRAFT caveat not consistently embedded in-line in citations
- Contradiction B (LOW): Stryker evidence source URL inconsistency

## Files To Modify

- `research/02-sources.md` — Update source notes for Go tools, arXiv, Fowler DRAFT
- `research/03-evidence.md` — Fix citations, add Go tools evidence, add subsumed mutants, fix Fowler DRAFT inline
- `research/05-report.md` — Add attribution qualifier for "five barriers", add note on mutation score thresholds, add Go tools, add subsumed mutants, fix Fowler inline
- `research/06-open-questions.md` — Update with resolved items

## Verification Plan

- Source verification: Verify all new sources are accessible
- Internal consistency: No contradictions in modified files
- Documentation accuracy: README not applicable (research-only lab)
- Final status: READY_FOR_RESEARCH_REAUDIT