# API Backward Compatibility

## 1. Modifying Payloads Safely
When a requirement changes the shape of a resource, breaking changes must be avoided.

**Requirement**: `phone` (string) -> `phones` (array).

### Strategy 1: Additive Fields (Expand in Place)
Add the new field alongside the old field.

```json
{
  "id": 1,
  "name": "Budi",
  "phone": "+628111",            // Legacy field
  "phones": ["+628111", "+628222"] // New field
}
```
Consumers utilizing the legacy `phone` field continue working unmodified. New consumers utilize the `phones` array.

### Strategy 2: URL Versioning
Launch a completely new endpoint (`/v2/customers`). Old traffic remains on `/v1/customers`.
- **Drawback**: Maintaining multiple API controllers/code paths indefinitely causes severe technical debt.

### Strategy 3: Internal Transformation Pipelines (Stripe Model)
- Expose the API through a single controller.
- Requests pass through a middleware pipeline of "Version Changes".
- The system processes data using the absolute latest schema natively in the backend.
- As data exits the API, transformation middleware mutates the response backward through time to match the API version pinned by the consumer.
- **Benefit**: Core application logic only thinks about the latest internal domain model. Old schema mapping is cleanly encapsulated in boundary transformation layers.

## 2. API Lifecycle and Deprecation

### The Deprecation Window
1. **ACTIVE**: The current standard.
2. **DEPRECATED**: The endpoint/field functions normally, but clients are notified (via headers, documentation, or email) to migrate.
3. **MIGRATION WINDOW**: Time allowed for consumers to upgrade (often 6-12 months for external clients).
4. **SUNSET**: The endpoint is formally turned off (returns 410 Gone or 404 Not Found).
5. **REMOVED**: Code is deleted from the repository.

## 3. Detecting Active Consumers (Observability)
Never guess if a field is still used. Ensure visibility before deletion.
- Log metrics (e.g., `legacy_field_used=true`) when clients access deprecated fields.
- Analyze load balancer / API gateway logs for old endpoint access.
- Proceed to Contract/Sunset only when legacy usage drops to exactly 0.
