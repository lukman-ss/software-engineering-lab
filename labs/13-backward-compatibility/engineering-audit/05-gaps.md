# Engineering Audit Gaps

Target Lab: labs/13-backward-compatibility

## Gaps Identified

### Gap 1: In-Memory Storage Simplification
- **Type**: IMPLEMENTATION_OVERCLAIM (Minor / Scoped)
- **Severity**: LOW
- **Description**: The implementation uses in-memory Go maps with `sync.RWMutex` to represent relational tables (`users` and `user_phones`). While this accurately models primary-key access, foreign-key relationships, and column dropping, it does not demonstrate database-level concerns such as transaction deadlocks, lock timeouts, or network partition failures during dual-writes.
- **Remediation**: Engineering notes (`engineering/02-implementation-notes.md`) already explicitly document this limitation and scope it with a ponytail comment. No code changes needed for this lab scope.

### Gap 2: HTTP Endpoint End-to-End Tests
- **Type**: MISSING_TEST
- **Severity**: LOW
- **Description**: While `internal/compat/service_test.go` and `tests/migration_test.go` thoroughly test the core domain and migration logic, and `internal/compat/handler.go` implements HTTP deprecation and sunset headers, there is no direct HTTP `net/http/httptest` test verifying header emission at the HTTP layer. The behavior is however demonstrated in manual code inspection and handlers are thin wrappers around tested domain functions.
- **Remediation**: Add a unit test using `net/http/httptest` to assert `Deprecation: true` and `Sunset` response headers in future iterations. Non-blocking.
