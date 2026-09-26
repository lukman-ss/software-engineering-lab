# Audit Plan — Research Audit: Load Testing (Lab 15)

Target Lab: `labs/15-load-testing`  
Audit Scope: Research artifacts only (`labs/15-load-testing/research/`)  
Auditor Role: Independent Technical Research Auditor  
Audit Date: 2026-09-26  

## Files Reviewed
- `labs/15-load-testing/research/01-plan.md`
- `labs/15-load-testing/research/02-sources.md`
- `labs/15-load-testing/research/03-evidence.md`
- `labs/15-load-testing/research/04-contradictions.md`
- `labs/15-load-testing/research/05-report.md`
- `labs/15-load-testing/research/06-open-questions.md`

## Claims To Verify
1. Standard definitions & distinct purposes of 4 core performance test types (Average-Load, Stress, Spike, Soak/Endurance).
2. RED Method (Rate, Errors, Duration) and Four Golden Signals (Latency, Traffic, Errors, Saturation) applicability in load testing.
3. Tail latency / Percentiles (P95, P99) superiority over arithmetic averages.
4. SLO codification via threshold assertions (`p(95)<200`, `rate<0.01`).
5. Tool characteristics & suitability: k6 (JS/Go, API/HTTP), Locust (Python/gevent, distributed), JMeter (Java/XML, multi-protocol), Gatling (Scala/JVM).
6. Multi-layer bottleneck triage strategy (Application vs Database vs External API/Downstream).
7. External dependency testing tradeoffs: real sandbox calls vs latency-injected mocking.
8. Common pitfalls: `/health` only testing, inadequate dataset scaling, server unmonitored during tests, unquantified targets.
9. Progressive testing workflow: Smoke → Load → Stress → Soak.

## Code Execution
- NOT APPLICABLE in this research-only audit phase per pipeline override.

## Primary Risks
- Inaccessible or paywalled sources cited without full text access (e.g. ISO/IEC 25010, Gatling documentation HTTP 403, JMeter timeout).
- Overgeneralization of vendor-specific capabilities or guidelines as universal standards.
- Synthesized troubleshooting trees presented without formal citation backing.

## Audit Strategy
1. Cross-examine all 25 entries in `02-sources.md` against claims in `03-evidence.md` and `05-report.md`.
2. Verify primary vs secondary source integrity and reachability disclaimers.
3. Evaluate whether evidence claims strictly match the scope and text of cited sources.
4. Assess resolution of apparent contradictions in `04-contradictions.md`.
5. Identify unverified assertions, gaps, and open questions in `06-gaps.md`.
6. Formulate formal audit verdict in `07-verdict.md`.
