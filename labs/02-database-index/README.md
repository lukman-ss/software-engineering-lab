# Lab 02 — Database Index

> **Mental Model**: A well-structured index can reduce query time from seconds to milliseconds. But blindly adding indexes can hurt write performance and storage. Always investigate what the database is actually doing.

---

## Problem

You have a `service` table in an auto repair shop. The reports team complains that running the monthly branch report is unbearably slow:

```sql
SELECT *
FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;
```

This query:
- Selects all columns
- Filters by `branch_id = 2`
- Filters by `status = 'FINISHED'`
- Filters by `service_date` in January 2026
- Sorts by `service_date` descending

**Question**: Should we add an index? Which columns? What order?

---

## Dataset

The table is seeded with ~500,000 rows with realistic, skewed distributions via `generate_series()`:
- Branch 2 and 5 busiest (~20% each)
- Branch 1, 3, 4, 6 moderate (~20% each)
- FINISHED status: ~75% of rows (high-volume)
- CANCELLED status: ~24% of rows
- WAITING status: ~1% of rows (low-volume)
- PENDING_REFUND status: ~0.1% of rows (rare case for selectivity experiments)

---

## Learning Progression

### 1. Setup

Run the setup script to create the schema and seed data:

```bash
./scripts/setup_lab.sh
```

Or manually:
```bash
psql -d software_engineer_lab -f schema.sql
psql -d software_engineer_lab -f seed.sql
```

Verify data volume:
```sql
SELECT COUNT(*) FROM service;
-- Expected: ~500,000 rows
```

Check distributions:
```sql
SELECT status, COUNT(*) FROM service GROUP BY status;
SELECT branch_id, COUNT(*) FROM service GROUP BY branch_id;
```

---

### 2. Cardinality vs Selectivity

These are **NOT** the same thing:

| Concept | Definition | Example |
|---------|------------|---------|
| **Cardinality** | Number of distinct values | `status` has 5 distinct values |
| **Selectivity** | Fraction of rows filtered by predicate | `status = 'FINISHED'` → 75% selectivity = 0.75 |

**Key insight**: Low cardinality does NOT mean index is useless.
- `status = 'FINISHED'` (75% match) = low selectivity = Seq Scan may win
- `status = 'PENDING_REFUND'` (0.1% match) = high selectivity = Index wins!

Check PostgreSQL statistics:
```sql
SELECT attname, n_distinct, most_common_vals, most_common_freqs
FROM pg_stats
WHERE tablename = 'service' AND attname = 'status';
```

---

### 3. Measure Baseline Performance

Run the main query with timing:

```sql
\timing on
SELECT *
FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;
```

---

### 4. Inspect Execution Plan

Use `EXPLAIN (ANALYZE, BUFFERS)` to see the actual execution plan:

```sql
EXPLAIN (ANALYZE, BUFFERS)
SELECT *
FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;
```

---

## Understanding EXPLAIN Output

### Cost vs Time

| Field | Meaning |
|-------|---------|
| `cost` | Planner's estimate (abstract units, NOT milliseconds) |
| `actual time` | Microseconds from ANALYZE |
| `rows` (estimate) | Planner's guess |
| `Actual Rows` | Reality check |

### Buffers

| Type | Meaning |
|------|---------|
| `shared read` | Blocks had to be read from disk |
| `shared hit` | Blocks were in OS cache |

### Execution Plan Nodes

| Node | When You See It |
|------|-----------------|
| **Seq Scan** | No usable index, or index not helpful |
| **Index Scan** | Single index used |
| **Bitmap Heap Scan** | Multiple indexes combined via bitmap |
| **Sort** | ORDER BY couldn't use index order |
| **Limit** | `LIMIT` clause pushed down |

---

## Experiments

### Experiment 1: Baseline (queries/01-baseline.sql)
Run before any optimization. Observe Seq Scan, explicit Sort, total rows examined.

### Experiment 2: Single-Column Indexes (queries/02-single-column-index.sql)
Create three separate indexes. See how PostgreSQL combines them with BitmapAnd.

### Experiment 3: Composite Index (queries/03-composite-index.sql)
Test the optimal `(branch_id, status, service_date DESC)` index.

### Experiment 4: Column Order Demonstration (queries/04-column-order-experiment.sql)
Compare three different column orders. Learn why left-to-right ordering matters.

Key insight: A B-tree index can only be traversed when the **leftmost** columns are constrained.

### Experiment 5: Selectivity Analysis (queries/05-low-cardinality-selectivity.sql)
Test `status = 'FINISHED'` (75% match) vs `status = 'PENDING_REFUND'` (0.1% match).
See same index produce different plans based on predicate selectivity.

### Experiment 6: ORDER BY + LIMIT (queries/06-order-by-limit.sql)
Dashboard query pattern. See how proper index eliminates Sort and limits rows examined.

### Experiment 7: Covering Index (queries/07-covering-index.sql)
Test `INCLUDE` columns for index-only scans. Understand visibility map requirements.

### Experiment 8: Write Cost (queries/08-write-cost.sql)
Measure INSERT/UPDATE/DELETE performance with and without indexes. Understand write amplification and HOT updates.

### Experiment 9: Storage Cost (queries/09-storage-cost.sql)
Measure index sizes using `pg_relation_size`. Understand the disk footprint of optimization.

### Experiment 10: Usage Audit (queries/10-index-audit.sql)
Learn how to identify unused indexes using `pg_stat_user_indexes`. Understand why `idx_scan = 0` requires careful interpretation.

### Experiment 11: Redundant & Overlapping Indexes (queries/11-redundant-indexes.sql)
Detect exact duplicates and overlapping prefix indexes. Learn workload-driven cleanup rather than simplistic dropping.

### Experiment 12: Partial Indexes (queries/12-partial-index.sql)
Create partial indexes (`WHERE status = 'FINISHED'`). See how they reduce index size and improve performance for specific subsets.

### Experiment 13: Functions on Indexed Columns (queries/13-functions-on-indexes.sql)
Compare `DATE(service_date) = ...` (non-SARGable) with range predicates and expression indexes (`CREATE INDEX ON table((DATE(column)))`).

### Experiment 14: PostgreSQL Statistics & ANALYZE (queries/14-statistics-and-analyze.sql)
Inspect `pg_stats` (`n_distinct`, `most_common_vals`, `histogram_bounds`). See how running `ANALYZE` updates planner estimates.

### Experiment 15: When Seq Scan is Correct (queries/15-seqscan-is-correct.sql)
Destroy the misconception that Seq Scan = bad. See why reading 75% of the table sequentially beats index traversal + random heap access.

### Experiment 16: Production-Safe Index Creation (queries/16-production-safe-index.sql)
Learn how `CREATE INDEX CONCURRENTLY` avoids locking production tables during long builds.

### Experiment 17: Benchmark Harness (queries/17-benchmark.sql)
Repeatable comparison of baseline/no-index, single-column indexes, wrong composite order, recommended composite, and covering index. Runs each query 5× and captures Min/Max/Avg.

---

## Senior Engineer Exercise

Answer the following without looking at any solution:

```sql
SELECT *
FROM service
WHERE branch_id = 2
  AND status = 'FINISHED'
  AND service_date BETWEEN '2026-01-01' AND '2026-01-31'
ORDER BY service_date DESC;
```

1. What index would you create?
2. Why does column order matter?
3. How would you prove PostgreSQL is actually using it?
4. What impact does the index have on INSERT?
5. What impact does it have on UPDATE?
6. Could three independent indexes work?
7. Why might PostgreSQL use BitmapAnd?
8. Why might PostgreSQL still use Seq Scan?
9. Can the index remove the Sort?
10. What happens if FINISHED represents 90% of rows?
11. What happens if FINISHED represents 0.1%?
12. What changes with ORDER BY + LIMIT 20?
13. When could INCLUDE be useful?
14. When could a partial index be better?
15. How would you determine whether an existing index is unused?

<details>
<summary>Expected Reasoning</summary>

**1. Index candidate:**
```sql
CREATE INDEX idx_service_branch_status_date
ON service(branch_id, status, service_date DESC);
```
This is a candidate based on this query shape, not a universally correct index.

**2. Column order matters** because multicolumn B-tree indexes are searched from left to right. Equality predicates (`branch_id`, `status`) filter down the tree; the range predicate (`service_date`) bounds the scan. Column order determines which predicates can be used and in what order.

**3. Prove usage** by running `EXPLAIN (ANALYZE, BUFFERS)` and looking for "Index Scan" or "Index Only Scan" in the plan, plus seeing `rows=` values that match the filter selectivity.

**4. INSERT impact** - The index adds B-tree maintenance overhead. Each INSERT must find the correct key position and may cause page splits.

**5. UPDATE impact** - Updating an indexed column requires index modification. Updating non-indexed columns may use HOT updates to avoid index touch.

**6. Three independent indexes CAN work** because PostgreSQL can combine them using BitmapAnd: Bitmap Index Scan → Bitmap Heap Scan pattern. However, a well-designed composite index is usually cheaper.

**7. BitmapAnd** allows PostgreSQL to use multiple single-column indexes: each predicate filters rows independently, bitmaps are ANDed together, then the result fetches rows.

**8. Seq Scan might be chosen** when the predicate is not selective enough (e.g., 70% match) — sequential read is cheaper than random page access via index.

**9. Sort elimination** - The index provides rows in `service_date DESC` order, so the Sort node is unnecessary.

**10. 90% FINISHED** - The predicate would not be selective; Seq Scan would likely be chosen even with the index available.

**11. 0.1% FINISHED** - Index would be very effective; few heap pages need inspection.

**12. LIMIT 20** - Index allows early termination: stop after 20 matching rows instead of processing all.

**13. INCLUDE** useful when the query can satisfy all columns from the index (fewer columns), enabling Index Only Scan and avoiding heap access.

**14. Partial index** better when the query always filters `status = 'FINISHED'` and FINISHED is a large portion of the table—smaller index, faster maintenance.

**15. Check `pg_stat_user_indexes`** where `idx_scan = 0` with sufficient monitoring period, high index size, and no functional dependencies indicate an unused index.
</details>

---

## Files in This Lab

```
labs/02-database-index/
├── README.md           # This file - learning guide
├── schema.sql          # Table definition
├── seed.sql            # ~500,000 row realistic dataset
├── cleanup.sql         # Reset everything
├── queries/
│   ├── 01-baseline.sql
│   ├── 02-single-column-index.sql
│   ├── 03-composite-index.sql
│   ├── 04-column-order-experiment.sql
│   ├── 05-low-cardinality-selectivity.sql
│   ├── 06-order-by-limit.sql
│   ├── 07-covering-index.sql
│   ├── 08-write-cost.sql
│   ├── 09-storage-cost.sql
│   ├── 10-index-audit.sql
│   ├── 11-redundant-indexes.sql
│   ├── 12-partial-index.sql
│   ├── 13-functions-on-indexes.sql
│   ├── 14-statistics-and-analyze.sql
│   ├── 15-seqscan-is-correct.sql
│   ├── 16-production-safe-index.sql
│   └── 17-benchmark.sql
├── indexes/
│   └── create_indexes.sql
└── scripts/
    └── setup_lab.sh
```

---

## Key PostgreSQL Concepts

### Composite Index Column Order

For `WHERE a = ? AND b = ? AND c > ? ORDER BY c`:

```sql
CREATE INDEX idx ON table (a, b, c DESC);
```

- **Equality first**: Filter down index quickly
- **Range last**: Can use for filtering AND ordering
- **DESC**: Enables backward scan, eliminating Sort

### Leading Column Rule

For index on `(a, b, c)`:
- `WHERE a = ?` ✓ Uses index
- `WHERE a = ? AND b = ?` ✓ Uses index  
- `WHERE a = ? AND b = ? AND c > ?` ✓ Uses index
- `WHERE b = ?` ✗ Index not usable (no `a` constraint)

### Cardinality vs Selectivity

```text
cardinality = distinct values
selectivity = rows selected / total rows
```

Both matter! Low cardinality with high selectivity = index useful.

### SARGable Predicates

Predicates that allow PostgreSQL to use an index are **SARGable** (Search Argument Able). Wrapping indexed columns in functions (e.g., `DATE(col)`) makes them non-SARGable unless an expression index is used.

### The Planner Decisions

Indexes do not make decisions. The PostgreSQL planner chooses execution plans based on statistics in `pg_stats` (`ANALYZE`). A Seq Scan is often correct when a large percentage of rows match.

### Production-Safe Builds

Always use `CREATE INDEX CONCURRENTLY` on large production tables to avoid acquiring `ACCESS EXCLUSIVE` locks that block write operations.

---

## Running the Lab

```bash
# Setup
./scripts/setup_lab.sh

# Run all experiments
psql -d software_engineer_lab -f queries/01-baseline.sql
psql -d software_engineer_lab -f queries/02-single-column-index.sql
psql -d software_engineer_lab -f queries/03-composite-index.sql
psql -d software_engineer_lab -f queries/04-column-order-experiment.sql
psql -d software_engineer_lab -f queries/05-low-cardinality-selectivity.sql
psql -d software_engineer_lab -f queries/06-order-by-limit.sql
psql -d software_engineer_lab -f queries/07-covering-index.sql
psql -d software_engineer_lab -f queries/08-write-cost.sql
psql -d software_engineer_lab -f queries/09-storage-cost.sql
psql -d software_engineer_lab -f queries/10-index-audit.sql
psql -d software_engineer_lab -f queries/11-redundant-indexes.sql
psql -d software_engineer_lab -f queries/12-partial-index.sql
psql -d software_engineer_lab -f queries/13-functions-on-indexes.sql
psql -d software_engineer_lab -f queries/14-statistics-and-analyze.sql
psql -d software_engineer_lab -f queries/15-seqscan-is-correct.sql
psql -d software_engineer_lab -f queries/16-production-safe-index.sql
psql -d software_engineer_lab -f queries/17-benchmark.sql
```

---

## Theory Audit: Common Misconceptions Corrected

### ❌ Full Table Scan is always bad.
✅ **A sequential scan can be optimal when a query reads a large portion of a table.** Reading 70-90% of rows with index traversal plus random heap access (typically 10+ I/O per row) is almost always slower than a single sequential scan that reads pages consecutively.

### ❌ Never index low-cardinality columns.
✅ **Low cardinality alone does not determine index usefulness; predicate selectivity and data distribution matter.** An index on a 5-value column can be very useful if the predicate matches 0.1% of rows. An index on an infinite-unique column is useless if the predicate matches 99% of rows.

### ❌ Never create separate indexes when a composite index exists.
✅ **PostgreSQL can combine indexes with bitmap operations; composite indexes can still be better for important query patterns.** `BitmapAnd` allows the planner to intersect bitmaps from multiple indexes. Composite indexes win when the same prefix is used repeatedly and the index is smaller/faster than combining multiple bitmaps.

### ❌ Index columns are always evaluated strictly left to right and later columns become useless.
✅ **PostgreSQL multicolumn B-tree behavior allows the planner to use any prefix of the index key.** The leftmost column forms the primary partition. Subsequent columns allow finer filtering within each partition. The planner can use columns in any order as long as there is no column skip.

### ❌ If an index exists PostgreSQL will use it.
✅ **The planner chooses the estimated cheapest plan.** It compares all available indexes and the sequential scan, picking whichever it estimates has the lowest total cost based on `n_distinct`, `most_common_freqs`, and histogram data from `ANALYZE`.

### ❌ EXPLAIN cost is execution time.
✅ **Cost is an abstract unit based on CPU pages and I/O** used by the planner. `total_cost` vs `startup_cost` are relative estimates. `actual_time` (from `ANALYZE`) is wall-clock microseconds. Cost does NOT equal milliseconds.

### ❌ Index Only Scan means PostgreSQL never touches the heap.
✅ **It means PostgreSQL only needed index columns, but heap Fetches can still occur if the visibility map indicates potential unvacuumed tuples.** An `Index Only Scan` is only possible when the visibility map shows all tuples on the checked pages are visible. If `Heap Fetches > 0`, the visibility map was incomplete or the page contained recently updated/deleted rows.