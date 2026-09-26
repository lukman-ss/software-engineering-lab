# Research Gaps

## Gap 1

Type:
WEAK_SOURCE

Severity:
MEDIUM

Location:
`02-sources.md` (Source 8, 9, 10), `03-evidence.md` (Evidence 1, 13)

Problem:
The research relies on Wikipedia summaries for foundational academic concepts: Coffman Conditions (1971) and Two-Phase Locking (2PL). While Wikipedia accurately reflects the concepts, an authoritative engineering reference should cite the primary texts (Coffman et al., Bernstein) or a recognized industry standard textbook (e.g., Silberschatz) directly rather than routing through a tertiary crowdsourced platform.

Required Revision:
Verify and cite the primary academic papers or standard operating systems textbooks directly.

Can Be Approved Without Fix:
YES

---

## Gap 2

Type:
MISSING_SOURCE / SCOPE_ERROR

Severity:
LOW

Location:
`02-sources.md` (Source 12), `05-report.md` (Limitations)

Problem:
Oracle Database concepts were originally targeted but excluded due to a URL failure (ORA-00060). If Oracle compatibility is a requirement for the final lab, this leaves a coverage gap.

Required Revision:
If the lab strictly requires Oracle context, a valid Oracle 19c or 23c documentation source must be located. If the lab only requires PostgreSQL and MySQL, the Oracle references can remain cleanly excluded.

Can Be Approved Without Fix:
YES

---

## Gap 3

Type:
IMPLEMENTATION_GAP / UNVERIFIED_CLAIM

Severity:
LOW

Location:
`05-report.md` (Finding 5, Limitation 3)

Problem:
The research applies general deadlock principles to a specific "Sistem PPOB" (Payment Point Online Bank) context. As noted in the limitations, there is no primary source explicitly linking PPOB architectural patterns to deadlock occurrences. It is purely an illustrative extrapolation.

Required Revision:
Ensure the final lab documentation clarifies that PPOB is merely a pedagogical framing mechanism, not a specialized architectural pattern distinct from standard concurrent transactional processing.

Can Be Approved Without Fix:
YES
