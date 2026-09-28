# Changes Made

## Revision 1

Audit Issue:
MEDIUM — Martin Fowler bliki citation marked DRAFT used without inline caveats (audit/04-contradictions.md, audit/06-gaps.md: Gap 6)

Files Changed:
- `research/02-sources.md`
- `research/03-evidence.md`
- `research/05-report.md`

Action:
- Explicitly labeled all inline citations of Martin Fowler's bliki as pre-publication draft ("Martin Fowler (pre-publication draft; carries 'This is a draft entry' notice)").
- Treated as corroborating secondary insight rather than primary authoritative evidence.

Verification:
- In-line text matches source disclaimer.
- Draft status documented in limitations and contradiction analysis.

Status:
RESOLVED

---

## Revision 2

Audit Issue:
MEDIUM — Go tooling gap for mutation testing in lab research (audit/06-gaps.md: Gap 7, audit/07-verdict.md: Non-Blocking Issue 4)

Files Changed:
- `research/02-sources.md`
- `research/03-evidence.md`
- `research/05-report.md`

Action:
- Added Source 10 (`go-mutesting` - https://github.com/zimmski/go-mutesting) and Source 11 (`gremlins` - https://github.com/go-gremlins/gremlins, https://gremlins.dev) with full source metadata and operator capabilities.
- Added Evidence 13 detailing Go mutation testing landscape, AST mutation mechanisms, and comparison with JVM/JS maturity.
- Updated Finding 8 and Limitations in `05-report.md` to document Go tooling specifics.

Verification:
- Repositories and tools verified.
- Clear distinction established between mature JVM/JS ecosystems and emerging Go tooling.

Status:
RESOLVED

---

## Revision 3

Audit Issue:
MEDIUM — ACH arXiv preprint unverified text / single industry source (audit/06-gaps.md: Gap 3, audit/07-verdict.md: Non-Blocking Issue 3)

Files Changed:
- `research/02-sources.md`
- `research/03-evidence.md`
- `research/05-report.md`

Action:
- Added Source 12 (`arXiv:2501.12862`) abstract metadata detailing trial statistics (9,095 mutants, 571 privacy-hardening test cases across 10,795 Android Kotlin classes, 7 platforms, 0.79/0.47 precision/recall rising to 0.95/0.96 with preprocessing).
- Maintained honest notation that the abstract was verified while raw binary PDF text extraction was limited.
- Attributed trial numbers explicitly to Meta's internal context (Kotlin, privacy domain) to prevent overgeneralization.

Verification:
- Claims corroborated against arXiv abstract.
- Scope explicitly defined as implementation-specific trial.

Status:
RESOLVED

---

## Revision 4

Audit Issue:
MEDIUM — Overgeneralization of "Five barriers to mutation testing" (audit/06-gaps.md: Gap 4)

Files Changed:
- `research/03-evidence.md`
- `research/05-report.md`

Action:
- Added explicit attribution: "According to Harman (Meta, 2025)..." and qualified that the five barriers are Meta/Harman's formulation, not a universal academic standard.

Verification:
- Claim classification updated to INTERPRETATION / implementation-specific framing.

Status:
RESOLVED

---

## Revision 5

Audit Issue:
LOW — Missing subsumed mutants case (audit/06-gaps.md: Gap 8)

Files Changed:
- `research/02-sources.md`
- `research/03-evidence.md`
- `research/05-report.md`

Action:
- Added Source 13, Evidence 14, and Finding 10 defining subsumed mutants and explaining their role in mutation score interpretation.

Verification:
- Definition aligns with standard literature.

Status:
RESOLVED

---

## Revision 6

Audit Issue:
LOW — Mutation score threshold ambiguity (audit/06-gaps.md: Gap 5)

Files Changed:
- `research/05-report.md`

Action:
- Added Finding 11 explicitly documenting that there is no universal industry-standard mutation score threshold (e.g. 80% vs 90%), explaining project-dependent contextual selection.

Verification:
- Documented clearly without fabricating standard numbers.

Status:
RESOLVED
