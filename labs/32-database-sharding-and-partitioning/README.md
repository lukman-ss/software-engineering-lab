# Database Sharding and Table Partitioning Lab

This lab provides runnable implementations and tests demonstrating database sharding, logical table partitioning, consistent hashing vs hash modulo routing, scatter-gather vs Global Secondary Index (GSI) queries, and distributed ID generation strategies.

## Overview of Components

1. **Logical Partitioning (`internal/partitioning`)**:
   - In-engine range partitioning.
   - Demonstrates partition pruning when filtering by partition range keys.
   - Fast partition dropping for bulk lifecycle operations.

2. **Physical Sharding (`internal/sharding`)**:
   - Independent shard instances (`Shard`).
   - `ModuloRouter`: Demonstrates naive modulo distribution and heavy data movement on resize.
   - `ConsistentHashRouter`: Ring-based hashing with virtual nodes minimizing data movement during cluster scale-out.
   - `Cluster`: Scatter-gather parallel execution and Global Secondary Index (GSI / Lookup Vindex) point-lookups.

3. **Distributed ID Generation (`internal/idgen`)**:
   - RFC 9562 compliant time-ordered UUIDv7 generator.
   - Vitess-style Sequence Block Allocator.

## Quick Start

### Run Tests
```bash
go test -v ./...
go test -race ./...
```

### Run Demo CLI
```bash
go run ./cmd/demo
```
