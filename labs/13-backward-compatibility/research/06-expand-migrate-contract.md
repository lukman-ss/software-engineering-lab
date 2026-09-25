# The Expand -> Migrate -> Contract Pattern

## Overview
Also known as the **Parallel Change** pattern (formulated by Joshua Kerievsky and popularized by Danilo Sato / Martin Fowler). It allows breaking incompatible changes into non-breaking, incremental phases:

```text
EXPAND  -->  MIGRATE  -->  CONTRACT
```

## Phase 1: Expand
Augment the system so that it supports both the old format and the new format concurrently.
1. Add new columns or tables in the database with relaxed constraints (e.g., nullable).
2. Add new fields/endpoints to the application.
3. Deploy changes. Existing consumers notice no changes and do not break.

## Phase 2: Migrate
Gradually transition all data and consumers from the old contract to the new contract.
1. **Application Dual Writing**: Application writes data to both the old and new structures.
2. **Backfill**: Background script or worker copies all historical data from the old structure to the new structure.
3. **Dual / Fallback Reading**: Application begins reading from the new structure, with fallback to the old structure if empty.
4. **Consumer Migration**: Update external/internal clients, workers, and mobile apps to start consuming the new contract.
5. **Observability Verification**: Monitor metrics until all clients have completely switched and read/write requests on the old schema cease.

## Phase 3: Contract
Remove the old format once nothing relies on it.
1. Stop writing to the old database schema / fields.
2. Drop deprecated code paths and legacy adapter transformations.
3. Clean up the database: remove legacy columns/tables or drop deprecated endpoints.
4. Add desired hard constraints (e.g., `NOT NULL` on new columns) if required.
