# 02 - Source Audit

## Source 1

Claimed Title: Parallel Change
Claimed Publisher: Martin Fowler / Danilo Sato (Thoughtworks)
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
- None. Foundational article for the Expand-Migrate-Contract pattern.

Assessment:
PASS

---

## Source 2

Claimed Title: API Versioning
Claimed Publisher: Stripe
URL: https://docs.stripe.com/api/versioning

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- The cited URL details Stripe's client SDK versioning and `Stripe-Version` header behavior. The research cites Stripe alongside database dual-write/read strategies, but this specific page covers HTTP header pinning and release cadence, not database migration mechanics.

Assessment:
PASS

---

## Source 3

Claimed Title: Avoiding downtime in migrations
Claimed Publisher: GitLab
URL: https://docs.gitlab.com/ee/development/database/avoiding_downtime_in_migrations.html

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Document is strongly tied to GitLab's specific architecture (Ruby on Rails, ActiveRecord schema cache, and PostgreSQL). The 3-release column dropping rule (`ignore_column` in M, post-deployment drop in M+1, remove ignore in M+2) stems from ActiveRecord's boot-time schema cache, not a generic constraint of all relational databases.

Assessment:
PASS

---

## Source 4

Claimed Title: AIP-180: Backwards compatibility
Claimed Publisher: Google Cloud API Design
URL: https://cloud.google.com/apis/design/compatibility

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Authoritative guideline defining wire, source, and semantic compatibility, additive changes, and breaking change classifications.

Assessment:
PASS
