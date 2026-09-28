# Engineering Design

Target Lab: `labs/36-cors-and-csrf`
Research Status: APPROVED

## Concept To Prove

1. **CORS is browser-side response-read relaxation under SOP, NOT backend write protection**: Cross-origin requests (e.g. form POST or cross-site fetch) are received and executed by the backend regardless of CORS headers.
2. **Wildcard `Access-Control-Allow-Origin: *` does not protect against CSRF**: State mutation still occurs; browsers only block JavaScript from reading the response when credentials are used.
3. **Simple vs Non-Simple Preflight**: Simple requests (`GET`, `POST` with `application/x-www-form-urlencoded`, `multipart/form-data`, `text/plain`) do not trigger preflight; non-simple requests (custom headers, `application/json`, `PUT`/`DELETE`) trigger `OPTIONS` preflight.
4. **Defense in Depth**:
   - `SameSite` cookies (`Strict`, `Lax`, `None`)
   - Synchronizer Token Pattern / Signed Double-Submit Cookie (HMAC-SHA256)
   - Fetch Metadata (`Sec-Fetch-Site`)
   - Custom Header requirement for JSON APIs (`X-Requested-With` / `Content-Type: application/json` enforcement)

## Expected Behavior

- **Vulnerable Server**: Processes cross-origin state-changing POST requests carrying session cookies regardless of `Origin` header.
- **CORS Middleware**: Correctly handles preflight `OPTIONS`, validates allowed origins, reflects origin safely (no wildcard when credentials allowed), handles headers/methods.
- **CSRF Protected Server**:
  - Validates Anti-CSRF token (Synchronizer / Signed Double-Submit).
  - Validates Fetch Metadata (`Sec-Fetch-Site: cross-site` rejected on mutating methods).
  - Validates custom header / content-type for API endpoints.

## Failure Scenario

- Cross-origin form submission executes money transfer / balance mutation without CSRF token on vulnerable endpoint.
- Attacker site triggers state change despite strict CORS policy on backend if endpoint accepts simple requests without CSRF protection.

## Success Criteria

- All unit and integration tests pass cleanly (`go test ./...`).
- Zero race conditions detected under `go test -race ./...`.
- Demo program executes and showcases the exploit vs defenses sequentially with exact HTTP flows.

## Architecture

```
labs/36-cors-and-csrf/
├── cmd/
│   └── demo/
│       └── main.go
├── internal/
│   ├── cors/
│   │   ├── middleware.go
│   │   └── middleware_test.go
│   ├── csrf/
│   │   ├── token.go
│   │   ├── token_test.go
│   │   ├── middleware.go
│   │   └── middleware_test.go
│   ├── bank/
│   │   ├── app.go
│   │   └── app_test.go
├── tests/
│   └── integration_test.go
├── engineering/
│   ├── 01-design.md
│   ├── 02-implementation-notes.md
│   └── 03-execution-result.md
├── go.mod
└── README.md
```

## Components

1. **`internal/cors`**:
   - `Config`: AllowedOrigins, AllowedMethods, AllowedHeaders, AllowCredentials, MaxAge.
   - `Handler`: Preflight OPTIONS handling, response header decoration, credential validation.
2. **`internal/csrf`**:
   - `TokenManager`: Generates and verifies HMAC-signed double-submit tokens / synchronizer tokens.
   - `FetchMetadataMiddleware`: Evaluates `Sec-Fetch-Site` header.
   - `CSRFMiddleware`: Enforces token verification on state-changing methods (`POST`, `PUT`, `DELETE`, `PATCH`).
3. **`internal/bank`**:
   - In-memory bank account service (balance check, transfer).
   - Session store for cookie-based authentication.
   - HTTP routes simulating vulnerable vs protected endpoints.
4. **`cmd/demo`**:
   - Standalone CLI demonstration simulating Victim, Attacker, and Bank Server interactions.

## Test Strategy

- Unit tests for CORS header generation, preflight response, and credential constraints.
- Unit tests for CSRF token creation, HMAC verification, expiration, and tampering detection.
- Integration tests simulating:
  - Cross-Origin read attempt blocked by SOP/CORS vs allowed with proper CORS.
  - CSRF attack success on vulnerable endpoint.
  - CSRF attack failure on protected endpoint (Token, SameSite, Sec-Fetch-Site).
- Race detector validation for concurrent token verification and bank transactions.

## Execution Plan

1. Setup `go.mod` (Go 1.22+ standard library only).
2. Implement CORS package and unit tests.
3. Implement CSRF package and unit tests.
4. Implement Bank application domain and endpoints.
5. Implement Integration test suite.
6. Build and run demo in `cmd/demo`.
7. Generate execution records and implementation notes.

## Implementation Decisions

1. **Pure Go Standard Library**: No external routing/HTTP dependencies (utilizes `net/http` ServeMux enhancements in Go 1.22+ and `crypto/hmac`, `crypto/sha256`, `crypto/rand`).
2. **Signed Double-Submit Cookie**: Uses HMAC-SHA256 bound to user session ID and timestamp to protect against naive cookie injection while remaining testable and robust.
