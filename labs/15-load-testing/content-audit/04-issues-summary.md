# Issues & Discrepancies Summary

Target Lab: labs/15-load-testing
Audit Date: 2026-09-26

## Critical

| # | Issue | Location | Severity |
|---|---|---|---|
| 1 | Stress test metrics understated: Average 504.8ms vs actual 739ms; P95 981ms vs actual 1357ms; P99 1.175s vs actual 1.588s | `02-master-draft.md:199-204` | Critical |
| 2 | Latency multiplier claimed ~46x but actual ratio is ~63.5x (P95 Stress / P95 Smoke) | `02-master-draft.md:206` | Critical |

## Major

| # | Issue | Location | Severity |
|---|---|---|---|
| 3 | Code snippet omits `io.Copy(io.Discard, resp.Body)` before body close; misleading for replication | `03-code-snippets.md:87` | Major |

## Minor

| # | Issue | Location | Severity |
|---|---|---|---|
| 4 | Sources section lists only 2 references; 06-source-map.md references 3 additional sources not cited | `02-master-draft.md:248-250` | Minor |
| 5 | Execution variability caveat missing; single-run numbers presented as definitive | `02-master-draft.md:192-204` | Minor |
| 6 | HTTP client timeout (5s) not documented | `02-master-draft.md` (implementation section) | Minor |
| 7 | `activeReq` counter serves dual purpose but only delay-triggering role documented | `02-master-draft.md:117` | Minor |

## Verified Correct (No Issues)

- Smoke test metrics: accurate within rounding
- Server semaphore implementation: code matches documentation
- Percentile calculation: code and formula match
- Thread-safe per-VU result collection: correctly described
- Custom HTTP transport `MaxIdleConns: 1000`: correctly documented
- Race detector passes: verified
- Architecture diagram: accurate
- Key takeaways: accurate summary
- All conceptual content aligns with approved research
- Source map: correct file references

## Summary
- Total issues found: 7
- Critical: 2
- Major: 1
- Minor: 4
- Items requiring fix before content can be considered accurate: 3 (issues #1, #2, #3)