# Contradictions

No material contradictions discovered.

Minor note on scope: 
- Source 1 (Pact intro) states contract testing is "immediately applicable anywhere two services communicate" including mobile app + service.
- Source 3 (Pact FAQ - "When to use Pact") states Pact is "not good for public APIs" where consumers cannot be identified individually.

These are compatible: the intro describes the general technique; the FAQ describes Pact's specific applicability. No actual conflict.

Another nuance:
- Source 4 (Robinson, 2006) discusses contract testing in terms of XSD schemas, Schematron assertions, and provider contracts - the abstract, theoretical framing.
- Source 1/2/7 (Pact, 2011-2026) implements it via executable consumer-driven tests using mock providers and pact files.

These are complementary (theory vs implementation), not contradictory.

Regarding end-to-end tests:
- Lab text says "Jangan membuat 2.000 E2E test hanya untuk memastikan schema API tidak berubah" (don't build 2000 E2E tests just for schema change detection).
- Source 7 (Pact FAQ) says "It depends" - the amount of E2E depends on risk profile; recommends separating integration from functional aspects.

Not a disagreement: lab is cautioning against misuse; FAQ provides nuanced operational guidance. Both agree contract tests reduce (but do not eliminate) E2E need.
