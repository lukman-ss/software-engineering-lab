# Revision Plan

Target Lab: labs/25-rate-limiting-and-backpressure

Previous Audit Status: APPROVED_WITH_WARNINGS

## Non-Blocking Issues (from research-audit/07-verdict.md)

1. **Source 11 Title Mismatch** (LOW): Source 11 title in `02-sources.md` mentions "Cloudflare's Rate Limiting Documentation (Redis rate limiter)" while URL points to `redis.io/docs/latest/develop/use-cases/rate-limiter/`.

2. **Generic RabbitMQ URL** (LOW): Source 13 links to generic `rabbitmq.com/tutorials` index rather than specific QoS/consumer prefetch page.

3. **Circular Specification Evidence** (MEDIUM): Evidences 14 & 15 in `03-evidence.md` cite the prompt/topic specification as their sole evidentiary source instead of external industry literature.

4. **AWS SDK Constants Overgeneralization** (LOW): Report highlights AWS SDK specific constants (50ms/1000ms base delays, 20s max cap) without explicit contextual note that these are AWS SDK defaults, not universal distributed systems constants.

## Files To Modify

- `research/02-sources.md` — Fix Source 11 title and Source 13 URL
- `research/03-evidence.md` — Re-attribute Evidence 14 & 15 to external primary sources
- `research/05-report.md` — Add contextual note for AWS SDK constants

## Verification Plan

- [ ] Source verification: Source 11 title matches publisher/URL
- [ ] Source verification: Source 13 URL points to specific consumer prefetch documentation
- [ ] Source verification: Evidence 14 & 15 cite external primary sources (Stripe, AWS Well-Architected, Kubernetes)
- [ ] Content verification: AWS SDK constants contextualized in report
- [ ] Documentation consistency: All files internally consistent after changes
- [ ] Source list updated in 02-sources.md
- [ ] Validation: All sources verified reachable
- [ ] Result: READY_FOR_RESEARCH_REAUDIT