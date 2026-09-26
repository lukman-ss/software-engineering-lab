# Research Gap Analysis — Lab 26: Contract Testing

## Gap 1: Date Inconsistency for Source 6 Reference

- **Type:** OUTDATED_SOURCE / TYPO
- **Severity:** LOW
- **Location:** `research/05-report.md:line 208` vs `research/02-sources.md:line 58`
- **Problem:** `05-report.md` cites `pactflow.io/blog/contract-testing-using-json-schemas-and-open-api-part-1 (May 30, 2023)` whereas `02-sources.md` cites `Updated 5 January 2023`. The live URL has "Updated 5 January 2023" as the article update date.
- **Required Revision:** Standardize the citation date string in the report footer to match the source file.
- **Can Be Approved Without Fix:** YES (Minor non-material typographical discrepancy; URL and content are fully authentic).

---

## Gap 2: Clarification of Enum Case Changes across Frameworks

- **Type:** UNVERIFIED_CLAIM / OVERGENERALIZATION
- **Severity:** LOW
- **Location:** `research/05-report.md:line 101`
- **Problem:** The statement that `IN_PROGRESS -> in_progress` is universally breaking is true for standard string equality and case-sensitive JSON unmarshalers, but can be non-breaking if clients normalize string enums or use case-insensitive deserializers.
- **Required Revision:** Kept accurately in `06-open-questions.md`. No further research revision needed since it was explicitly marked with confidence boundaries.
- **Can Be Approved Without Fix:** YES.

---

## Gap 3: Client Strict Parsing Configuration

- **Type:** SCOPE_ERROR / WEAK_SOURCE
- **Severity:** LOW
- **Location:** `research/05-report.md:line 86-89`
- **Problem:** Additive changes are described as backward-compatible with general appeal to industry consensus, but client ecosystems that default to strict schema validation (e.g. XML Schemas, strict JSON schema validators, or specific Java/Python configs) will fail.
- **Required Revision:** The research noted this limitation in `04-contradictions.md` and `06-open-questions.md`.
- **Can Be Approved Without Fix:** YES.
