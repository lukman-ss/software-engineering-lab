# 06 - Research Gap Analysis

## Gap 1

Type:
IMPLEMENTATION_GAP

Severity:
MEDIUM

Location:
`labs/13-backward-compatibility/` root

Problem:
The lab currently contains only markdown research documents and a README file. It does not provide runnable code, migration scripts, or verification tests demonstrating the Expand-Migrate-Contract cycle in practice.

Required Revision:
Develop the practical demonstration code, schema migrations, and automated tests validating the dual write/read and contract phases.

Can Be Approved Without Fix:
YES (as a research foundation for subsequent implementation)

---

## Gap 2

Type:
OVERGENERALIZATION

Severity:
MEDIUM

Location:
`research/04-database-migration.md:8-13`

Problem:
The 3-release column dropping rule (Release M: ignore, Release M+1: drop, Release M+2: clean up ignore) is presented as a universal zero-downtime database migration mandate. In reality, this specific lifecycle is necessitated by Ruby on Rails' ActiveRecord boot-time schema cache in combination with GitLab's release cadence, rather than a universal requirement across all database access libraries (e.g. Go sqlx, jOOQ, or explicit SQL queries).

Required Revision:
Clarify that the 3-release ignore pattern is an implementation-specific mitigation for ORMs that perform wildcard `SELECT *` queries based on cached schema definitions.

Can Be Approved Without Fix:
YES

---

## Gap 3

Type:
SCOPE_ERROR

Severity:
LOW

Location:
`research/02-sources.md:14-19`, `research/05-api-compatibility.md:4`

Problem:
Stripe's API documentation is cited to support sunset and deprecation header strategies. In practice, Stripe uses long-term version pinning via the `Stripe-Version` request header without enforcing sunsetting on client code.

Required Revision:
Note the difference between Stripe's perpetual backward compatibility (version routing) and sunset-driven deprecation policies.

Can Be Approved Without Fix:
YES
