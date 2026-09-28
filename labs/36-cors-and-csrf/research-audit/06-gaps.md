# 06 — Research Gap Analysis

## Gap 1
Type: WEAK_SOURCE
Severity: LOW
Location: `research/02-sources.md:91`, `research/06-open-questions.md:21`
Problem: PortSwigger article states Chrome Lax-by-default was enforced "Since 2021", while Chromium release timeline dates Chrome 80 to Feb 2020 (rollout was paused during COVID-19 and resumed in 2020-2021).
Required Revision: None for conceptual correctness. Ensure downstream lab documentation specifies Chrome 80+ / 2020-2021 rollout.
Can Be Approved Without Fix: YES

---

## Gap 2
Type: MISSING_CASE
Severity: LOW
Location: `research/06-open-questions.md:9`
Problem: Behavior of SameSite and CORS in mobile native WebViews (Android WebView / iOS WKWebView) is not deeply analyzed in evidence files.
Required Revision: Add brief note in implementation guidelines indicating that default cookie policies in embedded WebViews may differ from standalone desktop browsers.
Can Be Approved Without Fix: YES

---

## Gap 3
Type: SCOPE_ERROR
Severity: LOW
Location: `research/05-report.md:159-162`
Problem: Clarification that CORS and CSRF mechanisms are strictly browser-enforced concepts, and do not apply to direct API consumers (cURL, microservice-to-microservice, mobile apps not using WebViews).
Required Revision: Already documented under `Limitations` in `05-report.md`. No further fix needed.
Can Be Approved Without Fix: YES

---

## Gap 4
Type: UNVERIFIED_CLAIM
Severity: LOW
Location: `research/06-open-questions.md:17`
Problem: Frequency of intermediaries (proxies, CDNs) stripping `Sec-Fetch-*` headers is unmeasured.
Required Revision: Rely on standard token-based defenses as primary or fallback whenever Fetch Metadata is used.
Can Be Approved Without Fix: YES
