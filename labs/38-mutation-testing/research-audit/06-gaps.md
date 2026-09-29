# Research Gaps Audit

Target Lab: `/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/38-mutation-testing`
Audit Date: 2026-09-29

---

## Gap 1
Type: WEAK_SOURCE
Severity: LOW
Location: `research/02-sources.md` Source 3; `research/05-report.md` Finding 2
Problem: Martin Fowler's draft bliki URL returned HTTP 404. It was marked REMOVED in `02-sources.md`, but a residual reference remains in the Sources line of Finding 2.
Required Revision: None strictly required; research explicitly notes the claim rests on reachable primary sources (PIT, Wikipedia). If future revisions occur, removing the residual draft URL string from Finding 2 is recommended for hygiene.
Can Be Approved Without Fix: YES

---

## Gap 2
Type: SECONDARY_SOURCE_ATTRIBUTION
Severity: LOW
Location: `research/02-sources.md` Source 1 & Source 9; `research/03-evidence.md` Evidence 1
Problem: Foundational academic papers (DeMillo, Lipton, Sayward 1978; Jia & Harman 2009) are cited via Wikipedia references rather than opened directly via DOI.
Required Revision: Disclosed accurately in source notes as "NOT VERIFIED for direct quotations." Acceptable for foundational historical context when Wikipedia and official tool documentation agree.
Can Be Approved Without Fix: YES

---

## Gap 3
Type: SCOPE_NUANCE
Severity: LOW
Location: `research/02-sources.md` Source 11; `research/05-report.md` Finding 8
Problem: The review of `gremlins` Go tool omitted the upstream README's explicit statement that it is in pre-1.0 (0.x.x) semver and designed primarily for smallish Go modules/microservices because large codebases can take hours.
Required Revision: Incorporating this detail strengthens the research's existing argument that Go tooling lacks enterprise maturity compared to PIT/Stryker.
Can Be Approved Without Fix: YES

---

## Gap 4
Type: SINGLE_ORGANIZATION_EVIDENCE
Severity: LOW
Location: `research/03-evidence.md` Evidence 12; `research/05-report.md` Finding 7
Problem: Meta ACH performance statistics (73% acceptance, 36% privacy relevance, 0.95/0.96 precision/recall) originate solely from Meta's research team (Mark Harman et al.) and have not been independently replicated by external industry or academic studies.
Required Revision: Correctly classified in research with MEDIUM confidence and explicitly noted in limitations (`05-report.md` Limitations 2 & 3). No further change required.
Can Be Approved Without Fix: YES
