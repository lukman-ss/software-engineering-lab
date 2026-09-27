# Audit Plan

## Target Lab
`labs/15-load-testing`

## Scope
Pipeline Override: Audit research only. Implementation and code audit are deferred to engineering stages. Research files will remain unmodified.

## Files Reviewed
- `research/01-plan.md`
- `research/02-sources.md`
- `research/03-evidence.md`
- `research/04-contradictions.md`
- `research/05-report.md`
- `research/06-open-questions.md`
- `research-revision/01-revision-plan.md`
- `research-revision/02-changes-made.md`
- `research-revision/03-revision-result.md`

## Claims To Verify
1. Standard performance test classifications cover six primary types: Smoke, Average-Load, Stress, Soak/Endurance, Spike, Breakpoint.
2. Core metrics to monitor: Percentile latency (P50, P95, P99), Error Rate, Throughput (RPS), and Resource Utilization (CPU, Memory, Disk, Network).
3. Tool characteristics: k6 (JS/API focus), Locust (Python/gevent greenlets), JMeter (GUI/multi-protocol/XML), Gatling (Scala DSL/JVM high-throughput).
4. Bottleneck isolation strategy: Component metrics breakdown (e.g. `http_req_waiting` / TTFB vs `http_req_connecting`), hypothesis-driven experimentation, and layer budgeting.
5. Common pitfalls: Testing `/health` only, unrealistic data volume, lack of server-side monitoring, undefined thresholds, testing on non-representative local hardware.
6. Concurrent user sizing methodology: Formula `Concurrent Users = Hourly Sessions * Average Session Duration (s) / 3600`.
7. SDLC timing: Pre-go-live, pre-promotions, post-architectural/database changes, continuous CI/CD automated gates.

## Code To Execute
None (Pipeline override: Audit research only).

## Primary Risks
- Duplicate sources listed as separate entries (Source 11 vs Source 19, Source 7 vs Source 17).
- Apache JMeter site connectivity timeout from current environment.
- Numeric recommendations presented without application domain context.
- General formulas applied to stateful multi-step booking flows without accounting for session think-times.

## Audit Strategy
1. Live network verification for all cited source URLs.
2. Claim-by-claim verification against cited source materials.
3. Verification of internal consistency and contradiction handling across research reports.
4. Categorization of gaps and severity mapping.
5. Final verdict formulation per audit quality gates.
