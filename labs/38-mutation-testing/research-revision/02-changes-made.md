# Changes Made

## Revision 1

Audit Issue:
MEDIUM — Finding 11 parenthetical `(80%, 85%, 90%)` mutation score thresholds stated in `research/05-report.md` without any source citation. Flagged in `06-open-questions.md` as NOT VERIFIED.

Files Changed:
- `research/05-report.md`
- `research/06-open-questions.md`

Action:
- Removed specific parenthetical numbers `(80%, 85%, 90%)` from Finding 11 claim statement in `05-report.md`.
- Replaced with neutral language: "recommended targets vary by project risk profile, codebase maturity, and team disincentive risks."
- Updated `06-open-questions.md` Weak Evidence item 2 to remove NOT VERIFIED numbers and replace with a neutral description of the gap.

Verification:
- Finding 11 claim now states no universally accepted threshold, consistent with PIT FAQ and Stryker docs (both confirmed to not prescribe a threshold).
- No numbers are asserted without citation.

Status:
RESOLVED

---

## Revision 2

Audit Issue:
LOW — Gap 4 / Claim 8 overgeneralization: "All major mutation testing tools follow the same pattern: generate mutants → run tests → classify as killed/survived" omits generative tools like Meta ACH which use mutants to generate tests rather than evaluate existing ones.

Files Changed:
- `research/03-evidence.md` (Evidence 8)
- `research/05-report.md` (Executive Summary)

Action:
- Updated Evidence 8 claim to qualify: "All major evaluation-focused mutation testing tools follow the same pattern... (while generative approaches such as Meta ACH extend this by generating tests to target surviving mutants)."
- Updated Notes section in Evidence 8 to explicitly distinguish evaluation workflow from generative workflow.
- Updated Executive Summary paragraph to add clarifying sentence: "Evaluation-focused tools follow a standard execution workflow (generate mutants → run tests → classify outcome), whereas modern generative tools extend this model by using mutants as guides for test synthesis."

Verification:
- Claim is now appropriately scoped to evaluation-focused tools.
- Generative approaches (ACH) are distinguished without removing or contradicting any verified evidence.

Status:
RESOLVED

---

## Non-Change Documented: Source 3

Audit Issue:
HIGH — Martin Fowler bliki URL returned HTTP 404.

Action:
- Source 3 was already marked `REMOVED (UNREACHABLE)` in `research/02-sources.md` during the research phase, with explicit note that no archived copy was found at web.archive.org.
- Evidence 1 and Evidence 3 in `research/03-evidence.md` already carry notes that the Fowler citation was removed and claims rest on Wikipedia and PIT.
- Finding 2 in `research/05-report.md` still contains a citation to Martin Fowler's draft URL in its Sources field. Verified that claim is fully backed by Wikipedia and PIT independently — no fix needed to claim support, but the stale URL reference in Finding 2 is noted.

Verification:
- Finding 2 Sources field lists: "Wikipedia, Martin Fowler (pre-publication draft; unreachable), PIT"
- Claim support does not rely on Fowler; Wikipedia and PIT independently confirm the mutation score formula and definition.
- The audit verdict confirmed this as non-blocking; existing disclosure text in research is accurate.

Status:
NO_CHANGE_NEEDED (audit verdict confirmed non-blocking; research files already carry correct disclosure)

---

## Non-Change Documented: Academic Secondary Attribution (Gap 3)

Audit Issue:
LOW — DeMillo et al. (1978) and Jia & Harman (2009) cited via Wikipedia references; primary papers not directly accessed.

Action:
- Source 1 and Source 9 in `research/02-sources.md` already carry explicit `NOT VERIFIED for direct quotations` disclaimers.
- No false attribution detected.
- Research Agent correctly disclosed secondary citation status.

Status:
NO_CHANGE_NEEDED (properly disclosed; audit verdict confirmed this as LOW, approved without fix)
