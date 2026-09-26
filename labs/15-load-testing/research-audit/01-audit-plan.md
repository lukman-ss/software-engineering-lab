# Audit Plan

## Target Lab
`labs/15-load-testing`

## Files Reviewed
- `research/01-plan.md`
- `research/02-sources.md`
- `research/03-evidence.md`
- `research/04-contradictions.md`
- `research/05-report.md`
- `research/06-open-questions.md`

## Claims To Verify
1. Six primary performance test types: Smoke, Average-Load, Stress, Soak/Endurance, Spike, Breakpoint.
2. Key metrics to monitor: Percentiles (P50, P95, P99), Error Rate, Throughput (RPS), Resource Utilization (CPU, Memory, Disk I/O, Network).
3. Tool comparison and positioning: k6 (JS/APIs), JMeter (GUI/Multi-protocol), Locust (Python/Greenlets), Gatling (Scala/JVM).
4. Bottleneck identification methodology across Application, Database, and External dependencies via metric correlation and request phase breakdown.
5. Common load testing pitfalls: testing `/health` only, unrealistic test data, unmonitored infrastructure, lack of predefined SLAs/thresholds, non-representative test environments.
6. Timing of load testing across the SDLC (pre-release, post-major changes, architecture shifts, migrations).

## Code To Execute
None (Pipeline override: Audit research only).

## Primary Risks
- Overgeneralized recommendations or numeric SLA targets presented as universal facts.
- Outdated or unreachable tool documentation URLs.
- Paywalled standards (e.g., ISO/IEC 25010) cited without disclaiming verification status.
- Tool-specific nuances (e.g., k6 cloud evaluation latency vs local, elastic cloud limits in breakpoint tests) mischaracterized.

## Audit Strategy
1. Verify source URLs, tiers, publishers, and relevance.
2. Cross-reference claims against citations in evidence and report files.
3. Check for internal contradictions and unaddressed open questions.
4. Document research gaps and assign appropriate severity.
5. Emit verdict based strictly on research validity and evidentiary rigor.
