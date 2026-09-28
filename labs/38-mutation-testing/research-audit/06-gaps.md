# Research Gap Analysis

## Gap 1

Type: WEAK_SOURCE

Severity: MEDIUM

Location: `research/02-sources.md:3-12`, `research/03-evidence.md:5-11`

Problem:
DeMillo, Lipton, Sayward (1978) — the foundational primary source — was NOT directly accessed. All claims attributed to this paper come through Wikipedia's bibliography. The original IEEE Computer paper text, quotations, and formulations were not inspected.

Required Revision:
Obtain the 1978 paper text (via ACM Digital Library, IEEE Xplore, or academic archive) and verify direct quotations. Clearly mark all claims derived from this paper as "Wikipedia-attributed, not directly verified."

Can Be Approved Without Fix:
YES — research document itself declares this limitation explicitly and does not fabricate direct quotations from the paper. No verbatim quotes from the original are presented; bibliographic attribution is honest.

---

## Gap 2

Type: WEAK_SOURCE

Severity: MEDIUM

Location: `research/02-sources.md:85-94`

Problem:
Jia & Harman (2009) survey — a Tier 1 peer-reviewed academic source — was not opened directly. Metadata via Wikipedia only.

Required Revision:
Access the paper via SemanticScholar, arXiv, or Wayback Machine PDF link included in Wikipedia reference list. Verify scope and claims.

Can Be Approved Without Fix:
YES — research explicitly annotates this as "bibliographic metadata only" and makes no specific quantitative claims attributable only to this paper.

---

## Gap 3

Type: UNVERIFIED_CLAIM

Severity: MEDIUM

Location: `research/03-evidence.md:113-119`, `research/05-report.md:87-93`

Problem:
ACH arXiv preprint (arXiv:2501.12862) referenced as corroboration for Meta statistics (73% acceptance, 36% privacy relevance). Research notes "PDF was fetched but content was binary." The arXiv PDF was not successfully read and verified.

Required Revision:
Access the preprint text directly (https://arxiv.org/pdf/2501.12862) and verify exact trial statistics, scope, language, and conditions.

Can Be Approved Without Fix:
YES — primary claims are sourced from the verified Meta Engineering Blog post, which is accessible and text-matches the cited statistics. The arXiv preprint is supplementary corroboration; its inaccessibility reduces depth but does not fabricate data.

---

## Gap 4

Type: OVERGENERALIZATION

Severity: MEDIUM

Location: `research/05-report.md:87-89`

Problem:
"Five barriers to mutation testing at scale" presented as a universal enumeration. These five categories are specifically Harman/Meta's framing as articulated in the blog post. Other academic surveys may frame barriers differently (or count different challenges).

Required Revision:
Add attribution qualifier: "According to Harman (Meta, 2025), five major barriers have historically prevented industrial mutation testing adoption."

Can Be Approved Without Fix:
YES — the research correctly attributes the framing to Meta/Harman and does not claim the five categories as a universal standard.

---

## Gap 5

Type: MISSING_SOURCE

Severity: LOW

Location: `research/06-open-questions.md:17-19`

Problem:
No consensus found on mutation score thresholds ("80%", "85%", "90% variously proposed elsewhere but NOT VERIFIED here"). Readers may misinterpret the absence of a threshold as one being unimportant. PIT documentation does not state a recommended minimum score.

Required Revision:
Explicitly note that "there is no industry-standard mutation score threshold" and reference any academic literature on threshold selection if available.

Can Be Approved Without Fix:
YES — gap is explicitly flagged in open questions.

---

## Gap 6

Type: WEAK_SOURCE

Severity: LOW

Location: `research/02-sources.md:27-33`

Problem:
Martin Fowler's bliki entry is explicitly marked as DRAFT by the author ("Please do not share or link to this URL until I remove this notice"). Using this source in supporting evidence contradicts Fowler's explicit request about the page's publication status.

Required Revision:
Avoid citing this source as supporting evidence until the draft notice is removed, or clearly label all references as "pre-publication draft, per author notice."

Can Be Approved Without Fix:
YES with caveat — research acknowledges this limitation in `02-sources.md` and `06-open-questions.md`. The information in the draft is consistent with published source material. Remove formal citations to this source if the research moves to publication stage.

---

## Gap 7

Type: IMPLEMENTATION_GAP

Severity: MEDIUM

Location: `research/05-report.md:99-106`, `research/06-open-questions.md:5`

Problem:
The lab's own specification targets Go for implementation, but mutation testing tooling for Go is acknowledged as limited. No Go-specific tooling source was inspected (e.g., GoMutator, gremlins, go-mutesting). The research does not provide guidance or references for Go mutation testing implementation.

Required Revision:
Research Go-specific mutation testing tools (e.g., https://github.com/zimmski/go-mutesting, https://github.com/avito-tech/gremlins) and document their capability, operator support, and maturity.

Can Be Approved Without Fix:
YES — gap is identified in `06-open-questions.md`. Research-only phase; implementation tooling research can proceed in implementation phase.

---

## Gap 8

Type: MISSING_CASE

Severity: LOW

Location: `research/05-report.md:121-131`

Problem:
Subsumed mutants are mentioned in Wikipedia but not included in the research evidence or findings. This concept is relevant for understanding mutation score accuracy.

Required Revision:
Add a brief note on subsumed mutants as a secondary concern that affects mutation score interpretation.

Can Be Approved Without Fix:
YES — subsumed mutants are an advanced refinement topic. Core concept is correct without it.
