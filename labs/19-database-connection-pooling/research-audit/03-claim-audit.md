# Claim Audit

## Claim 1: PostgreSQL max_connections default is 100 and acts as a hard slot limit
Location: `research/03-evidence.md:Evidence 1`, `research/05-report.md:Finding 1`
Evidence Provided: Direct quote from PostgreSQL documentation.
Source: PostgreSQL 18 Documentation
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Verified accurate.

## Claim 2: Database performance degrades past a "knee" due to contention and context switching
Location: `research/03-evidence.md:Evidence 3`, `research/05-report.md:Finding 4`
Evidence Provided: PostgreSQL Wiki description of disk thrashing, RAM usage, lock contention, context switches.
Source: PostgreSQL Wiki — Number Of Database Connections
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Verified accurate.

## Claim 3: Optimal pool sizing formula is ((core_count * 2) + effective_spindle_count)
Location: `research/03-evidence.md:Evidence 4`, `research/05-report.md:Finding 2`
Evidence Provided: Cited by PostgreSQL Wiki and HikariCP Wiki.
Source: PostgreSQL Wiki & HikariCP Wiki
Source Actually Supports Claim: YES
Classification: FACT / INTERPRETATION
Severity: LOW
Notes: Supported as a heuristic/starting point. Correctly noted that it was designed for HDDs and unverified for SSDs.

## Claim 4: Oracle Real-World Performance group demonstrated a 50x latency improvement by reducing pool size from 2048 to 96
Location: `research/03-evidence.md:Evidence 5`, `research/05-report.md:Finding 5`
Evidence Provided: Quoted from HikariCP Wiki referencing Oracle video demonstration.
Source: HikariCP Wiki — About Pool Sizing
Source Actually Supports Claim: YES
Classification: EXAMPLE
Severity: LOW
Notes: Secondary citation accurately noted in research. Confirmed directly in the HikariCP wiki text.

## Claim 5: HikariCP leakDetectionThreshold default is 0 with a minimum of 2000ms
Location: `research/03-evidence.md:Evidence 8`, `research/05-report.md:Finding 3`
Evidence Provided: HikariCP README snippet.
Source: HikariCP README
Source Actually Supports Claim: YES
Classification: IMPLEMENTATION-SPECIFIC
Severity: LOW
Notes: Accurate for HikariCP.

## Claim 6: PgBouncer transaction mode breaks session-level state
Location: `research/03-evidence.md:Evidence 10`, `research/05-report.md:Finding 8`
Evidence Provided: PgBouncer config documentation list of breaking features.
Source: PgBouncer Config Documentation
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Accurate.

## Claim 7: AWS RDS calculates PostgreSQL max_connections as LEAST(DBInstanceClassMemory/9531392, 5000)
Location: `research/03-evidence.md:Evidence 14`, `research/05-report.md:Finding 6`
Evidence Provided: AWS documentation table.
Source: AWS RDS User Guide — Connection Limits
Source Actually Supports Claim: YES
Classification: FACT / IMPLEMENTATION-SPECIFIC
Severity: LOW
Notes: Directly verified via AWS documentation.

## Claim 8: Azure PostgreSQL reserves 15 connections and recommends PgBouncer with 2-5x vCores
Location: `research/03-evidence.md:Evidence 12-13`, `research/05-report.md:Finding 6`
Evidence Provided: Azure Flexible Server Limits documentation.
Source: Microsoft Azure Documentation
Source Actually Supports Claim: YES
Classification: FACT / IMPLEMENTATION-SPECIFIC
Severity: LOW
Notes: Directly verified via Microsoft Learn documentation.

## Claim 9: Deadlock avoidance pool sizing formula is Tn * (Cm - 1) + 1
Location: `research/03-evidence.md:Evidence 17`, `research/05-report.md:Finding 11`
Evidence Provided: Quoted from HikariCP Wiki.
Source: HikariCP Wiki
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Accurate mathematical lower bound for resource deadlock avoidance.
