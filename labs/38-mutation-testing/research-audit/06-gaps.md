# Research Gap Analysis

Target Lab: labs/38-mutation-testing
Audit Scope: Research Files Only (PIPELINE OVERRIDE)

---

## Gap 1

Type: WEAK_SOURCE
Severity: MEDIUM
Location: research/02-sources.md (Source 1), research/03-evidence.md (Evidence 1)
Problem:
Foundational academic paper by DeMillo, Lipton, and Sayward (1978) — "Hints on test data selection: Help for the practicing programmer" — was not directly read or verified via DOI/IEEE Xplore; claims are secondary-attributed via Wikipedia citations.
Required Revision:
None blocking. The research honestly annotates this limitation ("Paper text NOT accessed directly... NOT VERIFIED for direct quotations"). Can be retained as secondary citation with existing disclaimer.
Can Be Approved Without Fix: YES

---

## Gap 2

Type: WEAK_SOURCE
Severity: MEDIUM
Location: research/02-sources.md (Source 9)
Problem:
Jia & Harman (2009) comprehensive survey paper on mutation testing was cited only by bibliographic metadata from Wikipedia; full text was not directly examined.
Required Revision:
None blocking. The research properly disclaims direct access ("Paper NOT opened directly; cite as secondary-attributed only").
Can Be Approved Without Fix: YES

---

## Gap 3

Type: WEAK_SOURCE
Severity: LOW
Location: research/02-sources.md (Source 3), research/03-evidence.md (Evidence 3)
Problem:
Martin Fowler's bliki entry carries an explicit "This is a draft entry" banner on the live web page.
Required Revision:
None needed. Resolved in revision phase by explicitly qualifying all in-line references as pre-publication drafts.
Can Be Approved Without Fix: YES

---

## Gap 4

Type: SCOPE_ERROR
Severity: LOW
Location: research/05-report.md (Finding 7)
Problem:
Meta ACH performance figures (73% acceptance, 36% privacy relevance, 0.95/0.96 precision/recall) stem from an internal Meta trial on Kotlin Android code, which may not generalize to other programming languages, domains, or smaller codebases.
Required Revision:
None needed. The research report explicitly notes these figures are Kotlin/Android/privacy trial specific and highlights that peer review is pending.
Can Be Approved Without Fix: YES

---

## Gap 5

Type: MISSING_SOURCE
Severity: LOW
Location: research/05-report.md (Finding 11), research/06-open-questions.md (Weak Evidence 2)
Problem:
Numeric targets often cited colloquially for mutation scores (e.g. 80%, 85%, 90%) lack an authoritative industry or academic consensus standard.
Required Revision:
None needed. The research report explicitly documents Finding 11: "No industry-standard mutation score threshold exists", preventing arbitrary recommendations.
Can Be Approved Without Fix: YES
