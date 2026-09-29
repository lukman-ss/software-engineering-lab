# Changes Made

## Revision 1

Audit Issue: MEDIUM — The 3-way transaction taxonomy (compensable, pivot, retryable) rests on a single source (Microsoft Azure Architecture Center).

Files Changed:
- research/05-report.md (Finding 4 — added source qualification sentence)
- research/03-evidence.md (Evidence 7 — noted single-source limitation in Source field and Notes)

Action:
- Added explicit note that this taxonomy is documented in detail only in Microsoft's Architecture Center, and that Richardson's public page references countermeasures without explicitly enumerating the 3-type step taxonomy.
- Maintained MEDIUM confidence for this claim; clarified source scope in evidence.

Verification:
- Claim remains supported (Microsoft source verified reachable and authoritative); now explicitly attributed rather than presented as universal literature consensus.

Status: RESOLVED

## Revision 2

Audit Issue: MEDIUM — The specific 6-item list of isolation countermeasures rests on Microsoft Azure Architecture Center without external cross-enumeration.

Files Changed:
- research/05-report.md (Finding 8 — expanded evidence and confidence notes)
- research/03-evidence.md (Evidence 11, Evidence 12 — added source-limitation notes)
- research/04-contradictions.md (Contradiction 2 — added clarification sentence)

Action:
- Explicitly labeled the 6-countermeasure enumeration as "Microsoft's formalization" rather than a universal standard.
- Noted that Microservices.io confirms the concept of countermeasures without providing an independent enumeration.
- Maintained HIGH confidence for anomaly existence; MEDIUM confidence for the specific 6-item list.

Verification:
- Claim correctly scoped to source; no overgeneralization introduced.

Status: RESOLVED

## Revision 3

Audit Issue: LOW — No automated protocol exists for failure of compensating transactions after retries; operational intervention is required.

Files Changed:
- research/05-report.md (Finding 9 — added note that no standard recovery protocol exists; production systems must add retry + backoff, DLQ, and manual intervention)

Action:
- Clarified that absence of a universal compensation-of-compensation protocol is an inherent pattern limitation, not a documentation gap.
- Listed expected operational patterns (retry with backoff, dead-letter queue, operator alerts) without inventing specific "best practice" numbers.

Verification:
- Claim aligns with Microsoft, AWS, and Temporal sources cited in Evidence 14; no new unsupported claims added.

Status: RESOLVED

## Revision 4

Audit Issue: LOW — Garcia-Molina & Salem (1987) text was verified via citation chain / ACM references rather than direct OCR/text parsing of the scanned PDF.

Files Changed:
- No changes made to research files; this limitation was already correctly documented in research/05-report.md (Limitations section) and research/06-open-questions.md (Unanswered Question 3).

Action:
- Verified that existing documentation accurately reflects the verification method (citation chain via ACM DOI, Cornell mirror link, Temporal footnote).
- No claim requires modification since the historical origin statement remains well-supported by multiple secondary sources and ACM indexing.

Verification:
- Existing documentation consistent with audit finding; no action needed.

Status: NO CHANGE NEEDED — Already documented
