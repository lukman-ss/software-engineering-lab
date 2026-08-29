# Lab 02 — Database Index

## Junior Thinking vs Senior Thinking

**Junior**:

```text
Query slow
→ add RAM
→ add CPU
→ create random index
```

**Senior**:

```text
identify slow query
→ reproduce
→ inspect EXPLAIN ANALYZE
→ inspect rows and buffers
→ inspect data distribution
→ form hypothesis
→ change index/query
→ measure again
→ inspect write/storage impact
→ monitor production workload
```

> **Do not optimize based on guesses. Optimize based on evidence.**

A senior engineer cares about:

- **latency** — actual query execution time
- **throughput** — queries per second the system can sustain
- **rows scanned** — how many rows the planner actually examined
- **buffers** — shared hits vs shared reads (cache effectiveness)
- **match fraction** — what fraction of rows each predicate filters
- **planner estimates** — how close `rows=` matches `Actual Rows`
- **write amplification** — cost of maintaining indexes on INSERT/UPDATE/DELETE
- **storage** — disk footprint of indexes
- **production locking** — whether index creation blocks writes
- **actual workload** — not just one query, but the full query mix

Not merely: **"Does an index exist?"**

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

The table is seeded with exactly 500,000 rows with realistic, skewed distributions via `generate_series()`:
- Branch 2 busiest (25%)
- Branch 3 moderate (20%)
- Branch 1, 4, 5 moderate (15% each)
- Branch 6 least busy (10%)
- FINISHED status: 70.0% of rows (high-volume)
- CANCELLED status: 20.0% of rows
- IN_PROGRESS status: 5.0% of rows
- WAITING status: 4.9% of rows
- PENDING_REFUND status: 0.1% of rows (rare case for selectivity experiments)

---

## Learning Progression

### 1. Setup

Run the setup script to create the schema and seed data:

```bash
./scripts/setup_lab.sh
```

Or manually:
```bash
psql -d se_lab -v ON_ERROR_STOP=1 -f labs/02-database-index/schema.sql
psql -d se_lab -v ON_ERROR_STOP=1 -f labs/02-database-index/seed.sql
```

Verify data volume:
```sql
SELECT COUNT(*) FROM service;
-- Expected: 500,000 rows
```

Check distributions:
```sql
SELECT status, COUNT(*) FROM service GROUP BY status;
SELECT branch_id, COUNT(*) FROM service GROUP BY branch_id;
```

### Baseline Verification

Verify the baseline has no dashboard-supporting secondary indexes:

```sql
SELECT indexname, indexdef
FROM pg_indexes
WHERE schemaname = 'public'
  AND tablename = 'service'
ORDER BY indexname;
```

Expected output: ONLY constraint-backed indexes (`service_pkey`, `service_invoice_no_key`).

No dashboard-supporting secondary indexes exist for (`branch_id`, `status`, `service_date`) in the baseline.

---

### 2. Cardinality vs Selectivity

These are **NOT** the same thing:

| Concept | Definition | Example |
|---------|------------|---------|
| **Cardinality** | Number of distinct values | `status` has 5 distinct values |
| **Match fraction** | `matching_rows / total_rows` | `status = 'FINISHED'` → 0.70 |

> **Convention used in this lab**: We use "match fraction" to refer to the numeric ratio (e.g., 0.70). We say a predicate is "highly selective" if the match fraction is small (e.g., 0.1%), meaning it filters out most rows. We say it is "not very selective" if the match fraction is large (e.g., 70%).

**Key insight**: Low cardinality does NOT mean index is useless.
- `status = 'FINISHED'` (70% match fraction) = not very selective = Seq Scan may win
- `status = 'PENDING_REFUND'` (0.1% match fraction) = highly selective = Index wins!

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
| `actual time` | **Milliseconds** from ANALYZE (not microseconds) |
| `rows` (estimate) | Planner's guess |
| `Actual Rows` | Reality check |

### Buffers

| Type | Meaning |
|------|---------|
| `shared read` | PostgreSQL had to **read the block into shared buffers**. The underlying OS may satisfy this from kernel page cache, so it does not necessarily mean physical disk I/O. |
| `shared hit` | Requested block was already present in PostgreSQL **shared buffers**. |

### Execution Plan Nodes

| Node | When You See It |
|------|-----------------|
| **Seq Scan** | No usable index, or index not helpful |
| **Index Scan** | Single index used |
| **Bitmap Heap Scan** | Heap rows fetched using a bitmap of tuple/page locations. The bitmap can come from: one Bitmap Index Scan; BitmapAnd over several indexes; or BitmapOr over several indexes. |
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

**Key insight for PostgreSQL 16**: For a B-tree index on `(a, b, c)`, constraints on leading columns determine how much of the index scan range can be bounded efficiently:
- `WHERE a = ?` → can efficiently constrain the index range
- `WHERE a = ? AND b = ?` → can constrain it further
- `WHERE a = ? AND b = ? AND c BETWEEN ...` → bounds a tight range

However, `WHERE b = ?` (skipping `a`) is not literally "impossible". PostgreSQL 16 can in principle use the index, but without a constraint on the leftmost column, it may need to scan a large or complete portion of the index. The planner will therefore often prefer Seq Scan or another index.

### Experiment 5: Selectivity Analysis (queries/05-low-cardinality-selectivity.sql)
Test `status = 'FINISHED'` (70% match) vs `status = 'PENDING_REFUND'` (0.1% match).
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
Create partial indexes for small subsets (e.g., `status = 'IN_PROGRESS'`). Compare index size and query performance against full indexes.

### Experiment 13: Functions on Indexed Columns (queries/13-functions-on-indexes.sql)
Compare `EXTRACT(YEAR FROM service_date)` (non-SARGable) with range predicates and expression indexes (`CREATE INDEX ON table((EXTRACT(YEAR FROM column)))`).

### Experiment 14: PostgreSQL Statistics & ANALYZE (queries/14-statistics-and-analyze.sql)
Inspect `pg_stats` (`n_distinct`, `most_common_vals`, `histogram_bounds`). See how running `ANALYZE` updates planner estimates.

### Experiment 15: When Seq Scan is Correct (queries/15-seqscan-is-correct.sql)
Destroy the misconception that Seq Scan = bad. See why reading 70% of the table sequentially beats index traversal + random heap access.

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
This is a workload-specific candidate, not a universally correct index.

**2. Column order matters** because leading equality predicates can tightly constrain the B-tree scan range. The range predicate then bounds the remaining portion. Column order affects how much of the index PostgreSQL must scan.

**3. Prove usage** by running `EXPLAIN (ANALYZE, BUFFERS)` and inspecting the actual plan tree to identify which index nodes participate (Index Scan, Index Only Scan, Bitmap Index Scan, Bitmap Heap Scan, BitmapAnd, or BitmapOr).

**4. INSERT impact** - The index adds B-tree maintenance overhead. Each INSERT must find the correct key position and may cause page splits.

**5. UPDATE impact** - Updating an indexed column requires index modification. Updating non-indexed columns may use HOT updates to avoid index touch.

**6. Three independent indexes CAN work** because PostgreSQL can combine them via bitmap operations. A composite index can be cheaper for this repeated query shape, but PostgreSQL can combine independent indexes through BitmapAnd. Measure the real workload.

**7. BitmapAnd** allows PostgreSQL to use multiple single-column indexes: each predicate filters rows independently, bitmaps are ANDed together, then the result fetches rows.

**8. Seq Scan might be chosen** when the predicate is not selective enough (e.g., 70% match) — sequential read is cheaper than random page access via index.

**9. Sort elimination** - The index provides rows in `service_date DESC` order. PostgreSQL B-tree indexes support backward scans, so an ASC-defined index can also satisfy DESC ordering via scan reversal.

**10. 90% FINISHED** - The predicate has a large match fraction (not highly selective); Seq Scan would likely be chosen even with the index available.

**11. 0.1% FINISHED** - The predicate has a very small match fraction (highly selective); index would be very effective; few heap pages need inspection.

**12. LIMIT 20** - Index allows early termination: stop after 20 matching rows instead of processing all.

**13. INCLUDE** is useful when all referenced columns are available from the index (fewer columns) and visibility map state allows heap visibility checks to be skipped, enabling Index Only Scan and avoiding heap access. Check "Heap Fetches" in the EXPLAIN ANALYZE output to verify.

**14. Partial index** better when the query always filters a small subset (e.g., `status = 'IN_PROGRESS'` ~5%). Partial index size advantage grows when the indexed subset is substantially smaller than the full table. Faster maintenance, less storage.

**15. Check `pg_stat_user_indexes`** where `idx_scan = 0` over a meaningful observation period is only a candidate signal. Before dropping, verify: statistics reset time, workload coverage, constraints/dependencies, rare jobs, maintenance/reporting traffic, index size, and write cost.
</details>

---

## Files in This Lab

```
labs/02-database-index/
├── README.md           # This file - learning guide
├── schema.sql          # Table definition
├── seed.sql            # Exactly 500,000 row realistic dataset
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
    └── verify_lab.sh
```

---

## Key PostgreSQL Concepts

### Composite Index Column Order

For `WHERE a = ? AND b = ? AND c > ? ORDER BY c DESC`:

```sql
CREATE INDEX idx ON table (a, b, c);
```

- **Equality first**: Filter down index quickly
- **Range last**: Can use for filtering AND ordering
- **Direction**: PostgreSQL B-tree indexes can be scanned **both forward and backward**. Both `(a, b, c)` and `(a, b, c DESC)` can support `ORDER BY c DESC` here. Explicit ASC/DESC definitions matter when you have **mixed ordering requirements** in a multicolumn ORDER BY (e.g., `ORDER BY x ASC, y DESC`).

### Leading Column Rule (PostgreSQL 16)

For an index on `(a, b, c)`:

Constraints on leading columns determine how much of the B-tree scan range can be bounded:

- `WHERE a = ?` → constrains index range efficiently
- `WHERE a = ? AND b = ?` → constrains further
- `WHERE a = ? AND b = ? AND c > ?` → bounds tight range
- `WHERE b = ?` → can use index in principle, but without a constraint on `a` the planner may need to scan a large or complete portion. Often prefers Seq Scan or another index.
- `WHERE b = ? AND c = ?` → same reasoning

PostgreSQL 16 can technically use any index structure, but performance depends on whether the leftmost column(s) provide useful constraints.

**PostgreSQL 18 compatibility note**: B-tree skip-scan optimization may change this behavior for certain cases. This repository targets PostgreSQL 16.

### Cardinality vs Match Fraction

```text
cardinality = distinct values
match_fraction = matching_rows / total_rows

A predicate is "selective" when match_fraction is small.
A predicate is "non-selective" when match_fraction is large.
```

Both matter! Low cardinality with a highly selective predicate = index useful.

### SARGable Predicates

Predicates that allow PostgreSQL to use an index are **SARGable** (Search Argument Able). Wrapping indexed columns in functions (e.g., `DATE(col)`) makes them non-SARGable unless an expression index is used.

### The Planner Decisions

Indexes do not make decisions. The PostgreSQL planner chooses execution plans based on statistics in `pg_stats` (`ANALYZE`). A Seq Scan is often correct when a large percentage of rows match.

### Production-Safe Builds

Use `CREATE INDEX CONCURRENTLY` on large production tables when you need to:
- Preserve concurrent write availability
- Fit between maintenance windows
- Avoid application downtime

Understand the trade-offs:
- Takes longer and requires more resources
- Cannot run in a transaction block
- Can leave an INVALID index if build fails
- Monitors old snapshots and can wait on long transactions

For smaller tables or during maintenance windows, a regular `CREATE INDEX` may be simpler and faster.

---

## Running the Lab

Requires a running PostgreSQL instance. Start infrastructure:

```bash
make infra-up
```

Then use Makefile targets:

```bash
# Setup: drops existing table, creates schema, seeds ~500k rows
make lab-02-setup

# Re-run seed only (preserves schema and indexes)
make lab-02-seed

# Run individual experiments
make lab-02-baseline    # Experiment 1: baseline query
make lab-02-indexes     # Create indexes
make lab-02-explain     # EXPLAIN ANALYZE experiments
make lab-02-benchmark   # Experiment 17: benchmark harness

# Or run SQL files directly via psql
psql -d se_lab -f labs/02-database-index/queries/03-composite-index.sql
psql -d se_lab -f labs/02-database-index/queries/17-benchmark.sql

# Automated structural verification
make lab-02-verify

# Cleanup: drops table and all indexes
make lab-02-clean
```

Configuration via environment variables (defaults shown):

```bash
DB_NAME=se_lab
DB_USER=postgres
DB_PASSWORD=postgres
DB_HOST=localhost
DB_PORT=5432
```

```bash
make lab-02-setup DB_HOST=your-host DB_USER=your-user
```

---

## Theory Audit: Common Misconceptions Corrected

### ❌ Full Table Scan is inherently bad.
✅ **A sequential scan can be optimal when a query reads a large portion of a table.** Reading a large percentage of rows with index traversal plus random heap access is often slower than a single sequential scan that reads pages consecutively.

### ❌ Do not index low-cardinality columns.
✅ **Low cardinality alone does not determine index usefulness; predicate selectivity and data distribution matter.** An index on a 5-value column can be very useful if the predicate match fraction is 0.1%. An index on an infinite-unique column is rarely useful if the predicate matches 99% of rows.

### ❌ Do not create separate indexes when a composite index exists.
✅ **PostgreSQL can combine indexes with bitmap operations.** `BitmapAnd` allows the planner to intersect bitmaps from multiple indexes. Composite indexes typically win when the same prefix is used repeatedly and the index is smaller/faster than combining multiple bitmaps, but separate indexes provide more flexibility for varied workloads.

### ❌ Index columns are evaluated strictly left to right and later columns become useless if one is skipped.
✅ **PostgreSQL multicolumn B-tree behavior allows the planner to use the index in principle, but skipping leading columns often results in scanning a large portion of the index.** The leftmost columns determine how much of the B-tree scan range can be bounded efficiently.

### ❌ If an index exists PostgreSQL will use it.
✅ **The planner chooses the estimated cheapest plan.** It compares available index paths and the sequential scan, picking whichever it estimates has the lowest total cost based on `n_distinct`, `most_common_freqs`, and histogram data from `ANALYZE`.

### ❌ EXPLAIN cost is execution time.
✅ **Cost is an abstract unit based on CPU pages and I/O** used by the planner. `total_cost` vs `startup_cost` are relative estimates. `actual_time` (from `EXPLAIN ANALYZE`) is reported in **milliseconds**. Cost does NOT equal milliseconds.

### ❌ Index Only Scan means PostgreSQL never touches the heap.
✅ **It means PostgreSQL only needed index columns to satisfy the query, but heap Fetches can still occur.** An `Index Only Scan` avoids heap access only when the visibility map indicates the tuples on the checked pages are visible to all transactions. If `Heap Fetches > 0`, the visibility map was incomplete or the page contained recently updated/deleted rows.