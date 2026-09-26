# Contradiction Audit

## Contradiction 1

Statement A:
"Eager loading (fetching relationships alongside the parent model in a single or batched operation) reduces the total query count from N+1 down to 1 or 2." (Presented as the primary recommended solution).

Location:
`05-report.md` (Finding 3)

Statement B:
"Using FetchType.EAGER either implicitly or explicitly for your JPA associations is a bad idea because you are going to fetch way more data that you need." / "Eager fetching is a code smell" (Presented as an anti-pattern).

Location:
`05-report.md` (Finding 4), `02-sources.md` (Source 4)

Type:
SOURCE_CONFLICT

Impact:
Conceptual ambiguity. Readers unfamiliar with ORM distinctions may be confused when "eager loading" is identified as the canonical solution to N+1 in one section, yet declared a dangerous "code smell" in the next.

Assessment:
The conflict arises from overlapping terminology:
1. **Mapping-Level Eager Fetching (`FetchType.EAGER`)**: Hardcoded association configuration in JPA entities that forces the framework to load associations globally, even when unnecessary. This is universally regarded as a code smell.
2. **Query-Level Eager Loading (`JOIN FETCH` or Eloquent `with()`)**: Dynamic, query-specific instructions to fetch associations eagerly only when required for a specific business case. This is the canonical solution.
The research report did not explicitly decouple these two concepts, leaving a slight conceptual contradiction in the narrative.
