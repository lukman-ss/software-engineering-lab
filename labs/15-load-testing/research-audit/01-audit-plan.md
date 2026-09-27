# Audit Plan: Research Audit for Lab 15 (Load Testing)

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
1. Six primary load test types (smoke, average-load, stress, soak, spike, breakpoint) are industry-standard definitions.
2. Metrics (P50, P95, P99, RPS, error rate, CPU, memory, network/disk I/O) are backed by authoritative sources (k6, ISO/IEC 25010, Azure, Google SRE).
3. Bottleneck identification methodology (`http_req_waiting`, `http_req_connecting`, third-party latency correlation) is accurate and supported.
4. Concurrent user calculation formula `(Hourly sessions * session duration) / 3600` is accurate per source.
5. Common pitfalls (health check testing only, small dummy data, laptop testing, lack of targets) are verified in sources.
6. SDLC timing recommendations for load testing are backed by engineering frameworks.
7. Tool comparison claims (k6, JMeter, Locust, Gatling) accurately reflect official tool documentation.

## Code To Execute
- PIPELINE OVERRIDE: Code execution and code/implementation auditing are skipped per pipeline override instructions. Only research artifacts will be audited.

## Primary Risks
- Inaccessible or dead URLs in source list.
- Misrepresentation of source content or Tier misclassification.
- Claiming standard compliance (e.g. ISO/IEC 25010) based only on Wikipedia without noting source limitations.
- Overgeneralized claims presented as universal truths (e.g., specific percentage increase in stress testing, fixed threshold metrics).

## Audit Strategy
1. **Source Audit**: Verify 21 sources listed in `02-sources.md` for URL validity, publisher accuracy, tier classification, and relevance.
2. **Claim Audit**: Audit findings in `05-report.md` and evidence entries in `03-evidence.md` against severity model and factual accuracy.
3. **Contradiction Analysis**: Review `04-contradictions.md` for completeness and identify any unrecorded internal or external contradictions.
4. **Gap Analysis**: Synthesize missing evidence, weak sources, open questions, and overgeneralizations into `06-gaps.md`.
5. **Verdict Generation**: Render final evidence-based verdict in `07-verdict.md`.
