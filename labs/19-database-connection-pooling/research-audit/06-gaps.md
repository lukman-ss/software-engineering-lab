# Research Gap Analysis

## Gap 1

Type:
MISSING_CASE

Severity:
LOW

Location:
research/06-open-questions.md (Dynamic vs Fixed Size Pools)

Problem:
The research focuses primarily on fixed-size connection pools (matching standard sizing formulas) and does not detail concrete latency trade-offs of dynamic pool scaling (min-idle scaling up to max-pool-size) during sudden microsecond traffic spikes.

Required Revision:
Include brief note or future benchmark testing dynamic pool scaling latency in production implementation.

Can Be Approved Without Fix:
YES

---

## Gap 2

Type:
SCOPE_ERROR

Severity:
LOW

Location:
research/06-open-questions.md (PgBouncer Prepared Statement Trade-offs)

Problem:
Transaction pooling historically prevented protocol-level prepared statements. While modern PgBouncer (v1.21+) supports `max_prepared_statements`, memory overhead across distinct application schemas warrants explicit lab testing if prepared statements are utilized.

Required Revision:
Verify whether the upcoming practical lab exercises session-level vs transaction-level pooling features.

Can Be Approved Without Fix:
YES
