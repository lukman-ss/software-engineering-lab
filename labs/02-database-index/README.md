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
│   └── 10-index-audit.sql
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
```