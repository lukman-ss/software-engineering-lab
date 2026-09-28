# Content Audit Plan

Target Lab: `labs/24-slo-sli-error-budget`

## Audit Scope

Review generated publication content against:
- Approved research: `research/05-report.md`, `research/03-evidence.md`, `research/04-contradictions.md`
- Approved engineering: `engineering/01-design.md`, `engineering/02-implementation-notes.md`, `engineering/03-execution-result.md`
- Approved audits: `research-audit/07-verdict.md` (APPROVED), `engineering-audit/06-verdict.md` (APPROVED), `engineering-audit-opensource/06-verdict.md` (APPROVED with 2 non-blocking findings)
- Content revision record: `content-revision/01-changes-made.md`

## Content Files to Audit

1. `content/01-content-brief.md`
2. `content/02-master-draft.md`
3. `content/03-code-snippets.md`
4. `content/04-diagrams.md`
5. `content/05-key-takeaways.md`
6. `content/06-source-map.md`

## Quality Gates

1. **Accuracy**: All claims traceable to approved research or verified implementation
2. **Code Fidelity**: Snippets verbatim from approved source files
3. **Test Alignment**: Test claims match actual test implementations and results
4. **Completeness**: All core concepts covered (SLI, SLO, Error Budget, Burn Rate, Multi-Window, Criticality)
5. **Transparency**: Engineering audit findings (LatencyThreshold unused, per-rule window fields unimplemented) disclosed
6. **No Hallucination**: No fabricated benchmarks, incidents, or platform-specific bias
7. **Formatting**: Consistent structure, clear diagrams, correct cross-references

## Audit Process

1. Read all reference materials
2. Verify each content file against sources
3. Document issues (blocking/non-blocking)
4. Output verdict