# Docs vs Code

## README vs Code
README sections: NOT NULL (23502), CHECK (23514), UNIQUE (23505), FOREIGN KEY (23503), PARTIAL UNIQUE INDEX (WHERE deleted_at IS NULL), run tests (`go test -v ./...`, `go test -race ./...`), run demo (`go run ./cmd/demo`).

Code reality:
- NOT NULL checks: `engine.go:51` email, `engine.go:54` username, `engine.go:128` orders.user_id. Matches README.
- CHECK: `engine.go:60` age >= 18, `engine.go:64` status enum, `engine.go:133` total_cents > 0. Matches README example predicates (age >= 18, status IN (...), total_cents > 0).
- UNIQUE: `engine.go:80` `emailIndex` / constraint name `users_email_key`. Matches README.
- FOREIGN KEY: `engine.go:138` `fk_orders_user`. Matches README.
- PARTIAL UNIQUE INDEX: `engine.go:73-83` `activeEmails` gated on `u.DeletedAt == nil`. Matches README.
- Commands: README lists `go test -v ./...`, `go test -race ./...`, `go run ./cmd/demo` — all executed, all consistent.

Assessment: PASS — README accurately reflects implementation surface.

## Engineering design doc vs Code
`engineering/01-design.md` Section "Architecture" lists packages: `internal/db`, `internal/store`, `internal/errors`, `internal/domain`, `internal/service`.

Actual packages: `internal/dberr`, `internal/model`, `internal/engine`, `internal/store`.

| Doc Package | Actual Package |
|---|---|
| `internal/db` | `internal/engine` |
| `internal/errors` | `internal/dberr` |
| `internal/domain` | `internal/model` |
| `internal/engine` | (matches) |
| `internal/service` | n/a (logic lives in `internal/store`) |
| `internal/store` | `internal/store` |

Assessment: WARNING
Severity: MEDIUM (DOC_CODE_MISMATCH)
Notes:
Design doc references package paths that do not exist (`internal/db`, `internal/errors`, `internal/domain`, `internal/service`) and omits naming that exists (`internal/dberr`, `internal/engine`, `internal/model`). This is documentation drift from the actual file layout. The conceptual architecture described (storage engine integrity, SQLSTATE standard, partial index semantics, zero-dependency in-memory simulator) is nonetheless consistent with code. Not a blocking code defect.

## Engineering implementation notes vs Code
`engineering/02-implementation-notes.md` Section "Files Added" lists:
- `internal/dberr/errors.go` ✓ exists
- `internal/model/model.go` ✓ exists
- `internal/engine/engine.go` ✓ exists
- `internal/store/store.go` ✓ exists
- `internal/store/store_test.go` ✓ exists
- `cmd/demo/main.go` ✓ exists

All file claims match actual files.
"Known Limitations" (no live PostgreSQL; EXCLUDE omitted) — honest, matches simulation approach.
"Trade-offs" (coarse-grained table locks) — matches `sync.RWMutex` over whole engine.

Assessment: PASS

## Execution result doc vs actual execution
engineering notes match actual execution exactly (build/vet pass, tests pass with/without race, demo output as recorded). No fabricated text.

## Demo vs Claims
Demo output (recorded live):
- NOT NULL 23502 → mapped domain message. ✓
- CHECK 23514 (age=15, invalid status). ✓
- FOREIGN KEY 23503. ✓
- PARTIAL UNIQUE (soft delete re-registration). ✓
- Concurrency: 50 goroutines, 1 success, 49 rejected as 23505. ✓ matches README claim.
- SQLSTATE taxonomy verification (`IsConstraintViolation` true for 23505). ✓

Assessment: PASS

## Claims verification matrix
| Claim | Verified | Note |
|---|---|---|
| NOT NULL 23502 rejects missing email/username/orders.user_id | YES | engine.go + tests |
| CHECK 23514 rejects age<18, bad status, total_cents<=0 | YES | engine.go + tests |
| UNIQUE 23505 rejects duplicate email | YES | engine.go + test + demo |
| FOREIGN KEY 23503 rejects orphan order | YES | engine.go + test + demo |
| PARTIAL UNIQUE allows re-use after soft delete | YES | engine.go + test + demo |
| Safe concurrency → exactly 1 success / rest 23505 | YES | demo 50-goroutine live |
| Unsafe concurrency → >1 duplicates | YES | test + demo intent |
| Error map to domain messages | YES | demo output |
| IsConstraintViolation classifies codes | YES | TestErrorClassification |
