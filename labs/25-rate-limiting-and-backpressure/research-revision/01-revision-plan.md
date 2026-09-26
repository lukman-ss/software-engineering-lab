# Revision Plan

Target Lab: labs/25-rate-limiting-and-backpressure

Previous Audit Status: NEEDS_REVISION

## Blocking Issues

1. **Missing Research Findings**: Only `research/01-plan.md` exists. The core research findings document (`02-findings.md`) containing evidence extraction and source synthesis has not been generated.

2. **Formula Misattribution Risk**: Risk of confusing Little's Law ($L = \lambda W$) with deterministic queue accumulation ($(\text{arrival} - \text{processing}) \times \text{time}$).

## Non-Blocking Issues

1. **Unsubstantiated Operational Metric Claim**: The assertion that queue age is superior to queue depth lacks explicit broker/observability citations (AWS SQS `ApproximateAgeOfOldestMessage`, Kafka consumer lag / timestamp metrics).

2. **Generalization of Autoscaling Failure Modes**: Needs specific architectural framing on downstream database/dependency capacity constraints.

## Files To Modify

- `research/02-findings.md` — **CREATE** new findings document
- `research/03-sources.md` — **CREATE** updated source index
- `research/04-contradictions.md` — **UPDATE** with resolved contradictions
- `research/05-report.md` — **UPDATE** with verified claims
- `research/06-open-questions.md` — **UPDATE** with resolved/remaining questions

## Verification Plan

- [x] Source verification: AWS Exponential Backoff & Jitter blog (Marc Brooker, 2015) - confirmed
- [x] Source verification: Google SRE Book - Cascading Failures chapter - confirmed
- [x] Source verification: Little's Law original paper (Little, 1961) - DOI verified
- [x] Source verification: RFC 6585 (429), RFC 9110 (HTTP semantics), RFC 2697 (srTCM), RFC 6598 (CGNAT) - all reachable
- [ ] SQS CloudWatch metrics: `ApproximateAgeOfOldestMessage` - need to cite
- [ ] Kafka documentation: consumer lag metrics - need to cite
- [ ] Tests: N/A (Pipeline Override - research only)
- [ ] Build: N/A (Pipeline Override - research only)
- [ ] Demo: N/A (Pipeline Override - research only)
- [ ] Documentation consistency: Plan → Findings → Report alignment