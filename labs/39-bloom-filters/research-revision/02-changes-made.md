# Changes Made

Target Lab: `labs/39-bloom-filters`

## Revision 1

Audit Issue:
NON-BLOCKING / MEDIUM — Broken URL on Source 2 and Source 5 in `research/02-sources.md`

Files Changed:
- `labs/39-bloom-filters/research/02-sources.md`

Action:
- Updated Source 2 URL to `https://www.eecs.harvard.edu/~michaelm/postscripts/tr-02-05.pdf` (canonical Mitzenmacher technical report, HTTP 200).
- Updated Source 5 URL to `https://cassandra.apache.org/doc/latest/cassandra/architecture/storage_engine.html` (active Apache Cassandra storage engine documentation permalink).

Verification:
- Confirmed endpoints resolve to relevant authoritative content.

Status:
RESOLVED

---

## Revision 2

Audit Issue:
NON-BLOCKING / MEDIUM — Missing registered citations for corroborating literature in `research/02-sources.md`

Files Changed:
- `labs/39-bloom-filters/research/02-sources.md`

Action:
- Registered Source 7: Martin Kleppmann, *Designing Data-Intensive Applications* (O'Reilly Media, Chapter 3).
- Registered Source 8: Andrei Broder & Michael Mitzenmacher, *Network Applications of Bloom Filters: A Survey* (Internet Mathematics, 2004).

Verification:
- Formal source registry now completely covers references made in `03-evidence.md` and `05-report.md`.

Status:
RESOLVED

---

## Revision 3

Audit Issue:
NON-BLOCKING / MEDIUM — Heuristic Hash Set memory footprint comparison lacked explicit breakdown in `research/05-report.md`

Files Changed:
- `labs/39-bloom-filters/research/05-report.md`

Action:
- Clarified baseline calculation in Executive Summary: conventional in-memory hash set (Go map / Java HashSet) incurs ~48-96 bytes per entry due to pointer, bucket table, and string key object overhead, leading to ~50-100 MB for 1,000,000 items vs ~1.2 MB in Bloom filter.

Verification:
- Technical rationale is grounded in standard managed runtime memory layout benchmarks.

Status:
RESOLVED
