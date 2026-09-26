# Changes Made

## Revision 1

Audit Issue:
HIGH (per audit MEDIUM) — unsupported claim: "connection pool exhaustion" in Finding 2

Files Changed:
- research/runs/2026-09-25-n-plus-one-query-problem/05-report.md

Action:
- Removed "connection pool exhaustion" from Finding 2 since no cited source supports it as a stated mechanism.
- Finding now scopes impact to aggregate latency and repeated round-trips, which are covered by Source 1's quoted passage.

Verification:
- No source or evidence entry supports connection pool exhaustion claim. Removed.

Status:
RESOLVED

## Revision 2

Audit Issue:
LOW — Source 4 title mismatch

Files Changed:
- research/runs/2026-09-25-n-plus-one-query-problem/02-sources.md

Action:
- Changed Source 4 Title from "Eager fetching is a code smell" to "JPA and Hibernate FetchType EAGER is a code smell" to match the article's actual H1 (verified via source fetch).

Verification:
- Webfetch of https://vladmihalcea.com/eager-fetching-is-a-code-smell/ confirmed H1.

Status:
RESOLVED

## Revision 3

Audit Issue:
LOW — Source 4 publication date recorded as "Unknown"

Files Changed:
- research/runs/2026-09-25-n-plus-one-query-problem/02-sources.md

Action:
- Changed Source 4 Published from "Unknown" to "December 15, 2014" (visible on article page).

Verification:
- Webfetch showed "Posted on December 15, 2014 by vladmihalcea".

Status:
RESOLVED

## Revision 4

Audit Issue:
LOW/MEDIUM — Source 4 relevance description overstated; not used by any evidence entry

Files Changed:
- research/runs/2026-09-25-n-plus-one-query-problem/02-sources.md
- research/runs/2026-09-25-n-plus-one-query-problem/03-evidence.md

Action:
- Revised Source 4 relevance description to accurately reflect the article's focus on EAGER fetching strategy inconsistency across Hibernate query methods (Persistence Context vs JPQL vs Criteria API), not primarily memory bloat.
- Removed the "Corroborated By: Source 4" attribution from Evidence 4, which was incorrect (Evidence 4's quote comes from Source 1, and Source 4 does not actually serve as a memory-bloat evidence source).

Verification:
- Article content confirmed via fetch; memory-bloat evidence correctly sourced to Source 1 only.

Status:
RESOLVED

## Revision 5

Audit Issue:
MEDIUM — APM tracing "strictly necessary" without supporting source

Files Changed:
- research/runs/2026-09-25-n-plus-one-query-problem/05-report.md

Action:
- Softened "total request query counting or full APM tracing is strictly necessary" to "total request query counting is recommended, and APM tracing can provide additional visibility."
- No cited source supports APM as strictly necessary.

Verification:
- Areas of Agreement now reflects source support level.

Status:
RESOLVED

## Revision 6

Audit Issue:
LOW — vague corroboration labels ("General APM documentation", "REST architecture constraints")

Files Changed:
- research/runs/2026-09-25-n-plus-one-query-problem/03-evidence.md
- research/runs/2026-09-25-n-plus-one-query-problem/05-report.md

Action:
- Removed the "General APM documentation" corroboration from Evidence 2 (replaced with a note stating no additional source was cited).
- Removed the "REST architecture constraints" corroboration from Evidence 5 (replaced with scope note on GraphQL specificity).

Verification:
- No specific verifiable sources were available for these labels.

Status:
RESOLVED

## Revision 7

Audit Issue:
LOW — "severe memory bloat" / "OOM errors" language stronger than source supports

Files Changed:
- research/runs/2026-09-25-n-plus-one-query-problem/03-evidence.md
- research/runs/2026-09-25-n-plus-one-query-problem/05-report.md

Action:
- Changed Evidence 4 claim from "severe memory bloat" to "memory bloat by fetching excessive amounts of unneeded data".
- Changed Finding 4 claim from "severe memory bloat" to "excessive memory consumption".
- Removed "leading to memory bloat/OOM errors" language in Areas of Agreement (replaced with "unnecessary memory and resource consumption").
- Updated Conclusion to remove "strictly limiting payload sizes... to prevent memory bloat" phrasing and "necessarily tracking".

Verification:
- Source 1 supports "fetch way more data than needed" and labels EAGER a "bad idea"; does not use "severe" or mention OOM. Language now stays within source support.

Status:
RESOLVED

## Revision 8

Audit Issue:
LOW — "1 or 2" numeric claim presented as guarantee; "universally accepted" / "massive latency overhead" overstatements

Files Changed:
- research/runs/2026-09-25-n-plus-one-query-problem/05-report.md

Action:
- Changed Finding 3 claim from "reduces... down to 1 or 2" to "typically reduces the total query count; for a JOIN approach this can be one query, for a batched approach typically two", presenting the numeral as inference/example.
- Changed Areas of Agreement "universally accepted primary solutions" to "widely recommended approaches to mitigate" as demonstrated by frameworks.
- Scoped Finding 5 from "REST/GraphQL" to GraphQL-specific as demonstrated by Source 3, removing the "massive latency overhead" generalization.

Verification:
- Claim audit (Claim 3, 5, 6, 8) noted these as technically accurate but overstated/inferential.

Status:
RESOLVED
