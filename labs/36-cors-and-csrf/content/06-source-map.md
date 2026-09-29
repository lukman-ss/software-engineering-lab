# Source Map

Dokumen ini memetakan setiap seksi pada `02-master-draft.md` ke file riset, file implementasi, dan file pengujian yang relevan.

---

## Problem & Why This Matters

Research:
- `research/05-report.md` (Executive Summary, Finding 1, Finding 2)
- `research/03-evidence.md`
- `research-audit/07-verdict.md`

---

## Mental Model (SOP / CORS / CSRF Boundary)

Research:
- `research/05-report.md` (Finding 1 — Perbedaan Fundamental SOP, CORS, dan CSRF)
- `research/03-evidence.md`

Research Audit:
- `research-audit/03-claim-audit.md`

---

## Core Concept: Simple vs Preflight Request

Research:
- `research/05-report.md` (Finding 3 — Mekanisme Preflight OPTIONS dan Simple Requests)

Implementation:
- `internal/cors/middleware.go:44-93` (Fungsi `Handler`, penanganan `http.MethodOptions`)

Tests:
- `internal/cors/middleware_test.go:44-69` (`TestCORS_Preflight_Success`)
- `internal/cors/middleware_test.go:28-42` (`TestCORS_DisallowedOrigin_Preflight`)

---

## Core Concept: CORS Tidak Mencegah CSRF

Research:
- `research/05-report.md` (Finding 2 — CORS Bukan Pelindung Sisi Server dan ACAO: * Tidak Mencegah CSRF)

Implementation:
- `internal/bank/app.go:101-146` (`HandleTransferVulnerable`)
- `internal/cors/middleware.go:52-58` (Logika simple POST lewat ke `next.ServeHTTP`)

Tests:
- `tests/integration_test.go:54-87` (`TestIntegration_CORS_Does_Not_Prevent_CSRF_Execution`)

Demo:
- `cmd/demo/main.go:44-61` (Seksi `[2] DEMO: Cross-Origin Attack on Vulnerable Endpoint`)

---

## Core Concept: Wildcard `*` & Credentialed Requests

Research:
- `research/05-report.md` (Finding 2, Evidence: MDN CORS credentialed requests)

Implementation:
- `internal/cors/middleware.go:62-72` (Cek `AllowCredentials` dengan logika spesifik origin)

Tests:
- `internal/cors/middleware_test.go:71-95` (`TestCORS_Credentials_With_Wildcard_DisallowedInSpec`)

---

## Core Concept: SameSite Cookie Attribute

Research:
- `research/05-report.md` (Finding 4 — SameSite Cookie Attribute)

Implementation:
- `internal/bank/app.go` (BankServer tidak menyetel atribut cookie; konfigurasi dikelola di level infrastruktur — caveat terdokumentasi di engineering audit)

---

## Implementation: Signed Double-Submit Token

Research:
- `research/05-report.md` (Finding 5 — Pola Anti-CSRF Token: Synchronizer Token vs Double-Submit Cookie)

Implementation:
- `internal/csrf/token.go:27-59` (`NewTokenManager`, `GenerateToken`)
- `internal/csrf/token.go:62-110` (`ValidateToken`)

Tests:
- `internal/csrf/token_test.go:13-35` (`TestTokenManager_GenerateAndValidate`)
- `internal/csrf/token_test.go:37-49` (`TestTokenManager_ExpiredToken`)
- `tests/integration_test.go:289-311` (`TestIntegration_CrossSession_Token_Reuse_Rejected`)

---

## Implementation: Anti-CSRF Middleware

Research:
- `research/05-report.md` (Finding 5)

Implementation:
- `internal/csrf/middleware.go:24-53` (`Middleware.Handler`)

Tests:
- `internal/csrf/token_test.go:51-92` (`TestCSRFMiddleware_FormPost`)
- `tests/integration_test.go:89-115` (`TestIntegration_CSRF_Token_Prevents_Attack`)
- `tests/integration_test.go:117-161` (`TestIntegration_Legitimate_Flow_With_CSRF_Token`)
- `tests/integration_test.go:253-287` (`TestIntegration_CSRF_Token_In_Header`)

---

## Implementation: Fetch Metadata Middleware

Research:
- `research/05-report.md` (Finding 6 — Fetch Metadata)

Implementation:
- `internal/csrf/middleware.go:56-70` (`FetchMetadataMiddleware`)

Tests:
- `internal/csrf/token_test.go:94-118` (`TestFetchMetadataMiddleware`)
- `tests/integration_test.go:163-187` (`TestIntegration_SecFetchSite_Protection`)
- `tests/integration_test.go:313-336` (`TestIntegration_SecFetchSite_SameOrigin_Allowed`)

---

## Implementation: Custom Header Middleware

Research:
- `research/05-report.md` (Finding 3 — OWASP CSRF Prevention custom header caveat)

Implementation:
- `internal/csrf/middleware.go:73-86` (`RequireCustomHeaderMiddleware`)

Tests:
- `tests/integration_test.go:210-251` (`TestIntegration_CustomHeader_Protection`)

---

## Concurrency Safety

Implementation:
- `internal/bank/app.go:22` (`sync.RWMutex mu`)
- `internal/csrf/token.go` (Stateless token; `crypto/rand` aman untuk concurrent access)

Tests:
- `tests/integration_test.go:189-208` (`TestIntegration_Concurrency_RaceCondition` — `go test -race`)

Engineering Audit:
- `engineering-audit/06-verdict.md` (Race Detector: PASS)
- `engineering-audit-opensource/06-verdict.md` (Race Detector: PASS)

---

## Demo & Case Study

Demo:
- `cmd/demo/main.go:38-113` (Simulasi penuh: initial state, serangan endpoint rentan, serangan endpoint terproteksi, alur klien legitimate)

Engineering:
- `engineering/03-execution-result.md`
- `engineering-revision/03-revision-result.md`

---

## Caveats / Production Considerations

Research:
- `research/05-report.md` (Limitations — Lingkup Browser-Only & Ketergantungan XSS)

Research Audit:
- `research-audit/07-verdict.md` (Non-Blocking: RFC Reference Gap, Custom Header caveat)

Engineering Audit:
- `engineering-audit/06-verdict.md` (GAP-02: delimiter `:` token, GAP-03: negative branch testing)
