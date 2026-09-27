## Finding 1

Location: internal/contract/verifier.go lines 87-111 (response body validation)
Claimed Behavior: Verifier validates response headers, status code, and body schema against contract expectations.
Observed Implementation: Verifier validates status code and body schema (subset matching), but does not validate response headers (e.g., Content-Type) despite contract defining ResponseDefinition.Headers.
Assessment: WARNING
Severity: MEDIUM
Notes: Missing header validation could allow providers to return incorrect headers (e.g., wrong content type) while still passing contract verification. This deviates from claimed behavior of validating interactions fully.

## Finding 2

Location: internal/contract/verifier.go lines 15-19, 46-55 (Interaction and Verifier structs)
Claimed Behavior: ProviderState field in Interaction is used to set up server state before verification.
Observed Implementation: ProviderState is stored in the contract and transmitted but never used during verification; verification runs against live endpoints without state setup.
Assessment: WARNING
Severity: LOW
Notes: The lab uses deterministic endpoints (e.g., /v1/orders/ORD-123) that return fixed data, so provider state is not needed. However, the contract structure includes the field implying stateful verification capability which is not implemented.

## Finding 3

Location: internal/consumer/client.go lines 65-71 (contract validation in FetchOrder)
Claimed Behavior: Consumer client validates that fetched order matches expected contract constraints (non-empty customer.name, valid status enum).
Observed Implementation: Client validates CustomerName not empty and Status is either IN_PROGRESS or COMPLETED. This matches the contract expectations.
Assessment: PASS
Severity: LOW
Notes: Validation is correct and aligns with generated contract.

## Finding 4

Location: internal/provider/server.go (ProviderBreaking and ProviderDual)
Claimed Behavior: ProviderBreaking introduces three breaking changes: status casing, field rename, primitive type mutation. ProviderDual maintains V1 compatibility while exposing V2.
Observed Implementation: ProviderBreaking returns status "in_progress", customer.FullName instead of Name, and Total as string. ProviderDual returns V1-compatible response for /v1/orders/ path. Matches claims.
Assessment: PASS
Severity: LOW
Notes: Implementation correctly demonstrates breaking changes and backward compatibility.