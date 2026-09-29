# Revision Plan

Target Lab: `labs/39-bloom-filters`

Previous Audit Status: `APPROVED_WITH_WARNINGS`

## Blocking Issues

None.

## Non-Blocking Issues

1. **Broken Source URLs**:
   - Source 2 (`esa2006.pdf`) URL broken (404). Updated to canonical author tech report mirror `tr-02-05.pdf`.
   - Source 5 (`operating/bloom_filters.html`) URL broken (404). Updated to active Apache Cassandra storage engine architecture documentation permalink.
2. **Missing Registered Citations**:
   - Martin Kleppmann (Designing Data-Intensive Applications, O'Reilly) and Broder & Mitzenmacher (2004, Internet Mathematics) cited in evidence/report but missing from formal source catalog (`02-sources.md`).
3. **Unclarified Heuristic Memory Comparison**:
   - Heuristic memory estimation "~50-100 MB untuk Hash Set biasa" in Executive Summary lacked explicit memory layout overhead assumptions (pointer/bucket/object overhead in Go map / Java HashSet).

## Files To Modify

- `labs/39-bloom-filters/research/02-sources.md`
- `labs/39-bloom-filters/research/05-report.md`
- `labs/39-bloom-filters/research-revision/01-revision-plan.md`
- `labs/39-bloom-filters/research-revision/02-changes-made.md`
- `labs/39-bloom-filters/research-revision/03-revision-result.md`

## Verification Plan

- Source URL and citation completeness validation.
- Consistency check across research documentation.
- Confirmation of zero pending critical/high audit defects.
