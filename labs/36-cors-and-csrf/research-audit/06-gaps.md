# 06 — Research Gaps

## Gap 1
Type: UNVERIFIED_CLAIM — Chronological Accuracy
Severity: LOW
Location: `research/02-sources.md:91` (PortSwigger), `research/06-open-questions.md:21`
Problem: PortSwigger states "Since 2021" while Chromium / web.dev documents Chrome 80 (February 2020). The exact date of full Lax-by-default enforcement is not precisely pinned to a single authoritative source.
Required Revision: Cross-check Chromium release notes (chromiumdash.appspot.com or Blink-dev mailing list) for Chrome 80 release date and `SameSite=Lax-by-default` activation timeline.
Can Be Approved Without Fix: YES

---

## Gap 2
Type: MISSING_SOURCE — Normative Backup for Fetch Metadata
Severity: LOW
Location: `research/03-evidence.md:201-207`, `research/05-report.md:123-131`
Problem: The W3C Fetch Metadata specification is cited in `03-evidence.md:213` ("Corroborated By: W3C Fetch Metadata specification") but no explicit URL or reference to that W3C spec is listed in `02-sources.md`. Only OWASP and MDN are listed as sources for Evidence 12.
Required Revision: Add W3C Fetch Metadata specification (`https://www.w3.org/TR/fetch-metadata/`) to `02-sources.md` as an explicit source entry.
Can Be Approved Without Fix: YES

---

## Gap 3
Type: MISSING_CASE — Pre-Chrome 80 Legacy Browser Behavior
Severity: LOW
Location: `research/05-report.md:82-93`, `research/06-open-questions.md:1-5`
Problem: The research only documents Chrome 80+ (2020+) as Lax-by-default. No data or analysis for Safari, Firefox, or pre-80 Chromium legacy behavior exists in the research.
Required Revision: Footnote or appendix documenting approximate adoption rates of Lax-by-default across the browser ecosystem.
Can Be Approved Without Fix: YES (noted in open questions)

---

## Gap 4
Type: MISSING_CASE — SameSite Behavior in Mobile Webviews
Severity: LOW
Location: `research/06-open-questions.md:8-9`
Problem: Lab research does not address WebView behavior (Android WebView, WKWebView on iOS). These are relevant deployment environments.
Required Revision: Explicitly state that lab scope is limited to standard browser contexts.
Can Be Approved Without Fix: YES

---

## Gap 5
Type: OVERGENERALIZATION — "Custom Header Always Prevents CSRF"
Severity: MEDIUM
Location: `research/03-evidence.md:185-189`
Problem: Evidence 11 claims that requiring a custom request header (e.g. `X-CSRF-Token`) prevents CSRF. However, this is only effective for requests initiated via `fetch()`/XHR APIs where CORS preflight applies. An HTML `<form>` with a custom `Content-Type` that happens to match safelisted values (or browsers with relaxed enforcement) may not trigger preflight. The research correctly states "only effective for API accessed via fetch/XHR, not for traditional form submissions" in the Notes section — but the claim headline does not make this restriction visible.
Required Revision: Ensure the qualifying statement ("not applicable to HTML form submissions") is front-loaded in any derived educational content.
Can Be Approved Without Fix: YES (caveat already documented internally)

---

## Gap 6
Type: MISSING_SOURCE — RFC 6265 / 6265bis for SameSite Normative Reference
Severity: LOW
Location: `research/01-plan.md:44` (RFC 6265bis listed as expected source)
Problem: RFC 6265bis is listed in the search plan as an expected primary source. However, `02-sources.md` does not include any RFC entry. MDN Set-Cookie and web.dev are used as the primary SameSite references, which are secondary to the RFC itself.
Required Revision: Add RFC 6265bis (or Incrementally Better Cookies IETF draft) to the source list, or explicitly note that MDN and web.dev are accepted proxies for the RFC in educational context.
Can Be Approved Without Fix: YES

---

## Gap 7
Type: UNVERIFIED_CLAIM — "Two-Minute POST Grace Period" Security Impact
Severity: LOW
Location: `research/03-evidence.md:107`, `research/06-open-questions.md:22`
Problem: Research notes Chrome's 2-minute POST grace period for default Lax cookies, but no analysis of the specific security consequence of this window (e.g., can CSRF in POST occur within 2 minutes of login?) is performed. Open question acknowledges this but provides no deeper analysis.
Required Revision: Acceptable to leave as-is for an educational lab context; annotate clearly as not analyzed.
Can Be Approved Without Fix: YES
