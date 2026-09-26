# Research Audit Plan: SLO, SLI & Error Budget

## Target Lab
`labs/24-slo-sli-error-budget`

## Audit Scope
Research artifacts only (PIPELINE OVERRIDE):
- `research/01-plan.md`
- `research/02-sources.md`
- `research/03-evidence.md`
- `research/04-contradictions.md`
- `research/05-report.md`
- `research/06-open-questions.md`

*(Note: Code and implementation files are skipped during this research audit stage per pipeline instructions).*

## Files Reviewed
1. `research/01-plan.md`
2. `research/02-sources.md`
3. `research/03-evidence.md`
4. `research/04-contradictions.md`
5. `research/05-report.md`
6. `research/06-open-questions.md`

## Claims To Verify
1. Definitions of SLI, SLO, and SLA.
2. Common SLI categories (availability, latency, throughput, correctness, durability) across system types.
3. Time-based vs Aggregate availability calculation formulas.
4. Standard availability tables (99%, 99.9%, 99.99%, 99.999% downtime per period).
5. Error budget definition (100% - SLO) and governance mechanism for release velocity.
6. Target selection principles (avoiding 100%, avoiding setting targets solely on current performance, keeping simple).
7. Percentiles vs averages for latency SLIs.
8. Client-side vs Server-side SLI collection.
9. Burn rate alerting principles and specific threshold conventions.
10. Error budget remaining formula `100 * (current - target) / (100 - target)`.
11. Alerting on symptoms vs causes.
12. Cost scaling non-linearity (~100x per nine).
13. Safety margins and planned outages (anti-overachievement / Chubby case).

## Code To Execute
- None (research-only audit phase per pipeline override).

## Primary Risks
- Over-generalization of vendor-specific implementations (e.g. Datadog burn rate thresholds or status corrections) as universal SRE standards.
- Calculation discrepancies between specification examples and standardized availability tables.
- Lack of empirical backing for heuristics such as "100x cost per nine" or user variance preferences.

## Audit Strategy
1. Audit all 6 listed sources for reachability, authority, relevance, and tier accuracy.
2. Cross-examine all 21 evidence entries and 12 report findings against primary source texts.
3. Validate mathematical consistency of availability and error budget formulas.
4. Verify explicit identification of research gaps and vendor-specific limitations.
5. Provide evidence-based final verdict.
