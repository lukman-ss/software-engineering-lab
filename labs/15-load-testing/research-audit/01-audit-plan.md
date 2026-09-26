# Audit Plan: Research Audit for Lab 15 (Load Testing)

## Target Lab
`labs/15-load-testing`

## Audit Scope
Research artifact audit only (per pipeline override). Implementation/code execution excluded.

## Files Reviewed
- `labs/15-load-testing/research/01-plan.md`
- `labs/15-load-testing/research/02-sources.md`
- `labs/15-load-testing/research/03-evidence.md`
- `labs/15-load-testing/research/04-contradictions.md`
- `labs/15-load-testing/research/05-report.md`
- `labs/15-load-testing/research/06-open-questions.md`

## Claims To Verify
1. Six industry test types: smoke, average-load, stress, soak/endurance, spike, breakpoint.
2. Percentile monitoring (P50/P95/P99) + throughput + error rate + resource utilization.
3. Bottleneck isolation methodology (TTFB vs connecting time vs external dependencies).
4. Load calculation method for VUs ($Sessions/sec \times Duration$).
5. Tool selection characteristics (k6 JS vs Locust Python vs JMeter GUI vs Gatling Scala).
6. Common pitfalls (health-only testing, synthetic/low data, unmonitored infrastructure).
7. SDLC timing for performance testing.

## Code To Execute
None (Research-only audit per instruction).

## Primary Risks
- Inaccurate citation or stray references from previous labs (e.g. Spring DI mentions).
- Tool-specific claims overgeneralized as universal truths.
- Arbitrary numeric recommendations (e.g. ramp-up 5-15%, specific soak durations).
- Source availability/tier accuracy.

## Audit Strategy
1. Examine all 21 sources listed in `02-sources.md` for relevance, reachability, tier classification, and scope match.
2. Cross-check claims in `03-evidence.md` and `05-report.md` against cited sources.
3. Identify unsupported, overgeneralized, or erroneous claims.
4. Evaluate contradictions recorded in `04-contradictions.md`.
5. Check research completeness and gaps in `06-open-questions.md`.
6. Deliver final quality verdict.
