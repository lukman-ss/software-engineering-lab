# Research Gap Analysis

## Gap 1

Type: WEAK_SOURCE
Severity: LOW
Location: `research/02-sources.md` (Source 3)
Problem: Martin Fowler bliki entry at `https://martinfowler.com/bliki/MutationTesting.html` returned HTTP 404 during live audit check.
Required Revision: Keep as auxiliary corroboration or replace with active mirror/reference if available.
Can Be Approved Without Fix: YES (core claims are independently corroborated by PIT and Wikipedia sources).

---

## Gap 2

Type: UNVERIFIED_CLAIM
Severity: LOW
Location: `research/02-sources.md` (Source 1, Source 9)
Problem: Foundational papers (DeMillo et al. 1978, Jia & Harman 2009) were referenced bibliographically via secondary citations rather than direct full-text reading.
Required Revision: None required beyond the clear disclaimers already present in `research/02-sources.md`.
Can Be Approved Without Fix: YES (proper attribution and explicit disclaimers are maintained).

---

## Gap 3

Type: IMPLEMENTATION_GAP
Severity: LOW
Location: `research/05-report.md` (Finding 7, Limitations)
Problem: Meta ACH trial figures reflect single-organization results on Android Kotlin repositories focusing on privacy invariants.
Required Revision: Maintain context boundary in future lab exercises so students recognize ACH metrics are context-specific rather than universal baseline guarantees.
Can Be Approved Without Fix: YES (already scoped in report).
