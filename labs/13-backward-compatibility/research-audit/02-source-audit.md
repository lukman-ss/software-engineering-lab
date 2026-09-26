# Source Audit

## Source 1

Claimed Title: Parallel Change  
Claimed Publisher: martinfowler.com (ThoughtWorks)  
URL: https://martinfowler.com/bliki/ParallelChange.html  

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Canonical primary reference for the Expand-Migrate-Contract pattern.

Assessment:
PASS

---

## Source 2

Claimed Title: APIs as infrastructure: future-proofing Stripe with versioning  
Claimed Publisher: Stripe Engineering Blog  
URL: https://stripe.com/blog/api-versioning  

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Authoritative reference on date-based versioning, transformation pipelines, and backward-compatible evolution.

Assessment:
PASS

---

## Source 3

Claimed Title: Backwards compatibility (AIP-180)  
Claimed Publisher: Google API Improvement Proposals  
URL: https://google.aip.dev/180  

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Authoritative industry standard outlining source, wire, and semantic compatibility rules.

Assessment:
PASS

---

## Source 4

Claimed Title: Feature Toggle  
Claimed Publisher: martinfowler.com (ThoughtWorks)  
URL: https://martinfowler.com/bliki/FeatureFlag.html  

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- The bibliography lists `bliki/FeatureFlag.html` while `03-core-concepts.md` also refers to `articles/feature-toggles.html` (both by Martin Fowler / Thoughtworks). Both URLs are valid and reachable.

Assessment:
PASS

---

## Source 5

Claimed Title: ALTER TABLE  
Claimed Publisher: PostgreSQL Global Development Group  
URL: https://www.postgresql.org/docs/current/sql-altertable.html  

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Source is engine-specific to PostgreSQL. Research must avoid presenting `CONCURRENTLY` and non-rewriting defaults as universal SQL behavior (the research correctly includes caveats).

Assessment:
PASS

---

## Source 6

Claimed Title: API Versions  
Claimed Publisher: GitHub Developer Documentation  
URL: https://docs.github.com/en/rest/about-the-rest-api/api-versions  

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Documents breaking vs additive changes and standard deprecation / sunset header handling (RFC 8594).

Assessment:
PASS

---

## Source 7

Claimed Title: Schema Evolution and Compatibility for Schema Registry  
Claimed Publisher: Confluent Documentation  
URL: https://docs.confluent.io/platform/current/schema-registry/fundamentals/schema-evolution.html  

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Focuses on streaming serialization (Avro/Protobuf), but definitions of backward and forward compatibility are formal and applicable to distributed messaging.

Assessment:
PASS

---

## Source 8

Claimed Title: Blue Green Deployment  
Claimed Publisher: martinfowler.com (ThoughtWorks)  
URL: https://martinfowler.com/bliki/BlueGreenDeployment.html  

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Establishes the core rule of separating database schema evolution from application code deployments.

Assessment:
PASS

---

## Source 9

Claimed Title: Semantic Versioning 2.0.0  
Claimed Publisher: semver.org  
URL: https://semver.org/  

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Canonical definition for MAJOR, MINOR, and PATCH version increments.

Assessment:
PASS

---

## Source 10

Claimed Title: Pattern: Transactional outbox  
Claimed Publisher: Microservices.io  
URL: https://microservices.io/patterns/data/transactional-outbox.html  

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- The source describes outbox for messaging rather than schema migration specifically, but the underlying atomicity failure mode and dual-write drift are identical distributed systems problems.

Assessment:
PASS
