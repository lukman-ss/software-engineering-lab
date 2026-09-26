# Docs vs Code

## Comparison Matrix

| Source | Claim | Code Reality | Match |
|--------|-------|--------------|-------|
| README structure section | lists `cmd/demo/main.go`, `internal/{consumer,contract,model,provider}`, `tests/contract_test.go`, `engineering/{01-design,02-implementation-notes,03-execution-result}`, `go.mod` | Directory listing confirms those files exist plus `research/`, `research-audit/`, `engineering-revision/` (not mentioned but present). Core lab paths match. | PASS |
| README Running section | `go test -v ./...`, `go test -race ./...`, `go run ./cmd/demo` | All three commands execute successfully. | PASS |
| README Overview point 1 | Consumer declares minimal required schema | `GenerateMobileContract()` builds minimal interaction (id,status,customer.name,total). Verifier does subset match (ignores extra fields like notes). Implementation confirms. | PASS |
| README Overview point 2 | Provider CI gate verification before deployment | `Verifier.Verify()` runs contract against running server; demo shows ALLOWED/BLOCKED gates. Test proves gate logic. | PASS |
| README Overview point 3 | Breaking change detection: enum casing, field rename, primitive type mutation | Breaking provider implements all three; verifier catches all three (demo output shows 3 errors). | PASS |
| README Overview point 4 | Safe API evolution: V1 compat while exposing V2 | ProviderDual serves /v1 (compliant) + /v2 (evolved). Contract test passes on V1 path. | PASS (with note: no test asserts V2 endpoint shape directly; V2 existence proven by code inspection + handler logic, not by automated test) |

## Issues
- None blocking. README minimal but accurate.
- engineering notes and research claims out of scope per pipeline override (implementation+tests only).