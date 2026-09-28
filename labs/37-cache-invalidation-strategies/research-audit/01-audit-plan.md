# 01 — Audit Plan

Target Lab: `labs/37-cache-invalidation-strategies`
Audit Date: 2026-09-28
Audit Scope: Research Audit (PIPELINE OVERRIDE: research files only; code/implementation omitted)

## Files Reviewed

- `research/01-plan.md`
- `research/02-sources.md`
- `research/03-evidence.md`
- `research/04-contradictions.md`
- `research/05-report.md`
- `research/06-open-questions.md`

## Claims To Verify

1. **Cache-Aside (lazy loading)**: Mechanics, read/write order, consistency trade-offs (Sources 01, 04, 10, Evidence 1-2).
2. **Write-Through & Write-Behind**: Synchronous dual-write vs deferred store write; read-after-write guarantees vs crash risk (Sources 01, 03, 06, Evidence 3-4).
3. **Cache Stampede / Thundering Herd**: Definition, mechanics, failure mode, zero-hit-rate congestion collapse (Sources 04, 05, Evidence 5-7).
4. **Single-Flight / Lock Coalescing**: Request coalescing, Go `singleflight` package mechanics, distributed lock extra write cost (Sources 03, 04, 09, Evidence 8-9).
5. **Probabilistic Early Expiration (XFetch)**: Algorithm, exponential draw formula, beta parameter, optimality claim (Sources 04, 08, 09, Evidence 10-11).
6. **Stale-While-Revalidate**: RFC 5861 semantics, stale window, request-triggered revalidation requirement (Sources 02, Evidence 12-13).
7. **TTL & Jitter**: Function of desynchronization vs hot key stampede prevention (Sources 01, 05, 07, Evidence 15-16).
8. **Quantitative / Benchmark Claims**: 10,000 RPS, 500 concurrent goroutines, P99 metrics (Evidence 18).

## Primary Risks

- Unverified claims around XFetch optimality due to failed primary PDF parsing.
- Mathematical sign error/confusion in XFetch formula between lab spec and paper/Wikipedia.
- Synthetic/unverified quantitative benchmarks presented without explicit disclaimers.
- Unreachable URLs for Redis official documentation (Sources 10).
- Conflation of IETF Informational RFC 5861 with formal Internet Standards.

## Audit Strategy

- Audit all 11 sources listed in `02-sources.md` for reachability, relevance, and tier appropriateness.
- Verify whether extracted evidence items match source capabilities and accurately represent source content.
- Verify all major technical claims across report findings and contradiction matrices.
- Identify research gaps, unverified claims, and mathematical/formal inaccuracies.
- Render audit verdict (`07-verdict.md`) in accordance with audit guidelines.
