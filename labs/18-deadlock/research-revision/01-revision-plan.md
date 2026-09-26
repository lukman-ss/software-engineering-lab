# Revision Plan

Target Lab: labs/18-deadlock

Previous Audit Status: PRELIMINARY_AUDIT_PASSED (Audit: 2026-09-26)

## Blocking Issues

None.

## Non-Blocking Issues

1. **MySQL docs 403 Forbidden** — MySQL/InnoDB claims marked NOT VERIFIED in research.
2. **Coffman conditions via Wikipedia** — Tier 2 source instead of original 1971 paper.
3. **Backoff pattern derivation** — Retry/backoff+jitter described as "standard practice" rather than from a specific authoritative source.

## Files To Modify

- None required — research correctly documents and marks limitations

## Verification Plan

- Source audit: PASS — PostgreSQL sources verified reachable
- Claim audit: PASS — all claims either supported by primary sources or marked MEDIUM
- Research limitations: Documented in 05-report.md Limitations section and 06-open-questions.md

## Resolution Rationale

### Issue 1: MySQL Documentation Inaccessible
**Action**: Research already correctly marks MySQL claims as "NOT VERIFIED" and excludes them from factual claims. Per audit: "Can Be Approved Without Fix: YES (the lab focuses primarily on PostgreSQL/general concurrency)."

**Evidence**: 05-report.md Line 147: "MySQL primary sources gagal diakses (403) — perbandingan PostgreSQL vs MySQL tidak dapat di-cross-check dua sumber independen."

### Issue 2: Wikipedia for Coffman Conditions
**Action**: The 4 Coffman conditions are universally accepted theoretical computer science fundamentals. Wikipedia correctly cites the original 1971 paper. The audit notes: "Can Be Approved Without Fix: YES (the 4 conditions are universally accepted)."

**Evidence**: 05-report.md Line 148: "Paper asli Coffman 1971 tidak dibuka; rujukan via Wikipedia (Tier 2) sehingga confidence dibatasi meski high-consensus."

### Issue 3: Transaction Retry Backoff
**Action**: Research correctly marks this as MEDIUM confidence with a note that additional source would be needed for HIGH. The claim relies on Go stdlib primitives which the audit notes provide "primitives for it."

**Evidence**: 05-report.md Line 150: "Rekomendasi backoff/jitter/idempotency key untuk retry PPOB berasal dari pola umum, bukan kalimat literal PG docs — perlu sumber tambahan (mis. Stripe/Retrying transactions RFC) untuk naik ke HIGH."

## Source Verification

| Source | Status |
|--------|--------|
| PostgreSQL Docs (10 sources) | VERIFIED |
| MySQL/InnoDB Docs | NOT VERIFIED (403) - Research correctly marks as such |
| Wikipedia Coffman | CITED as Tier 2, with note about original paper |
| Go time package | VERIFIED (provides primitives) |