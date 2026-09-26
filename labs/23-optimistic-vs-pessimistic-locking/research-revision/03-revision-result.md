# Revision Result

Target Lab: labs/23-optimistic-vs-pessimistic-locking

Previous Audit Status: NEEDS_REVISION

## Issues

Critical:
- No research output existed — RESOLVED (4 research files created)

High:
- No verifiable sources (8 aspirational only) — RESOLVED (17 sources verified with URLs, publishers, content checks)
- No research findings answering 15 questions — RESOLVED (03-evidence.md + 05-report.md with classified findings)
- Isolation/MVCC version-specific behavior unaddressed — RESOLVED (version-specific claims in Findings 3, 5, 6, 12)

Medium:
- Universal framing risk (Q2, Q9) — RESOLVED (all claims scoped by vendor/version/isolation)
- Academic sources lack edition/year — RESOLVED (cited Silberschatz 7th ed 2019, Bernstein & Hadzilacos)
- Performance/decision thresholds without anchoring — RESOLVED (Evidence 11 cites Bernstein & Goodman 1981; Evidence 15 cites Microsoft retry pattern)
- Implementation patterns lack SQL/ORM doc citations — RESOLVED (Evidence 17, 18, 19, 20 with vendor doc URLs)

Low:
- Plan structure sound — MAINTAINED
- No README at lab root — PENDING (research-only stage, per pipeline override)

## Resolution

Resolved:
- Gap 1: Sources (17 verified with URLs and inspected excerpts)
- Gap 2: Findings (16 findings answering Q1-Q15, classified by evidence type)
- Gap 3: Version specificity (PostgreSQL 14+, 18; SQL Server 2022+; Hibernate 6.6+)
- Gap 4: Universal framing (all claims scoped by vendor/version/isolation)
- Gap 5: Academic sources (edition/year cited)
- Gap 6: Numeric guidance (anchored to Bernstein & Goodman, Microsoft patterns)
- Gap 7: Implementation patterns (concrete SQL with vendor doc URLs)

Partially Resolved:
- MySQL 8.0 documentation inaccessible during research — marked as inferred from cross-vendor patterns in Sources and Evidence

Unresolved:
- N/A — no blocking issues remain

## Validation

Build:
N/A (research-only, no build required)

Tests:
N/A (research-only, no tests to run)

Race Detector:
N/A (research-only)

Demo:
N/A (research-only, no demo code)

Documentation Accuracy:
N/A → IMPROVED (research files now present with verified sources vs. aspirational plan only)

Source Integrity:
FAIL → RESOLVED (17 verified sources replaces 8 aspirational entries)

Claim Support:
FAIL → RESOLVED (20 evidence items with per-claim citations; all claims classified)

## Remaining Risks

1. MySQL 8.0 official documentation was inaccessible — behavior inferred from PostgreSQL/SQL Server/Oracle patterns. Requires re-verification when MySQL docs become accessible.
2. Distributed database research (CockroachDB, Spanner) relies on cross-referenced sources, not direct vendor docs.
3. Performance benchmarks lack controlled methodology; crossover thresholds are theoretical/interpretive.

## Ready For Re-Audit
READY_FOR_REAUDIT