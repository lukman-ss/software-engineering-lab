# Changes Made

## Revision 1
**Audit Issue:** HIGH — Source 8 URL returns HTTP 404 (`https://github.com/cep21/circuitbreaker`)

**Files Changed:**
- research/02-sources.md

**Action:**
- Corrected URL from `https://github.com/cep21/circuitbreaker` to `https://github.com/cep21/circuit`
- Verified the repository exists and is the correct cep21 circuit breaker implementation

**Verification:**
- Source checked: https://github.com/cep21/circuit returns 200 OK with valid content

**Status:** RESOLVED

---

## Revision 2
**Audit Issue:** HIGH — Source 10 URL returns HTTP 404 (`https://sre.google/workbook/load-shedding/`)

**Files Changed:**
- research/02-sources.md

**Action:**
- Replaced with correct SRE Workbook chapter URL: `https://sre.google/workbook/managing-load/`
- Updated title from "Load Shedding — Google SRE Workbook" to "Managing Load — Google SRE Workbook"
- Updated relevance to reflect actual chapter content (load shedding vs circuit breaker, Dressy case study, load balancing interaction)

**Verification:**
- Source checked: https://sre.google/workbook/managing-load/ returns 200 OK with load shedding content

**Status:** RESOLVED

---

## Revision 3
**Audit Issue:** HIGH — Source 12 URL returns HTTP 404 (`https://learn.microsoft.com/en-us/azure/architecture/patterns/fallback`)

**Files Changed:**
- research/02-sources.md

**Action:**
- Replaced with valid Azure pattern: `https://learn.microsoft.com/en-us/azure/architecture/patterns/queue-based-load-leveling`
- Updated title to "Queue-Based Load Leveling — Microsoft Azure Architecture Center"
- Updated relevance to support async flow decoupling claim

**Verification:**
- Source checked: https://learn.microsoft.com/en-us/azure/architecture/patterns/queue-based-load-leveling returns 200 OK with relevant content on decoupling via queues

**Status:** RESOLVED

---

## Revision 4
**Audit Issue:** HIGH — AWS Builder's Library source cited in evidence/report but missing from source inventory

**Files Changed:**
- research/02-sources.md
- research/03-evidence.md
- research/05-report.md

**Action:**
- Added Source 12: "Timeouts, retries, and backoff with jitter — AWS Builder's Library" with URL `https://aws.amazon.com/builders-library/timeouts-retries-and-backoff-with-jitter/`
- Updated Evidence 5 source citation to reference Source 12
- Updated Finding 3 source citations to reference Source 12 and Source 6

**Verification:**
- Source checked: https://aws.amazon.com/builders-library/timeouts-retries-and-backoff-with-jitter/ returns 200 OK

**Status:** RESOLVED

---

## Revision 5
**Audit Issue:** HIGH — Claim "In-memory circuit breakers only track state on a per-instance basis" lacks supporting citations

**Files Changed:**
- research/05-report.md (Limitations section)

**Action:**
- Rephrased claim to be more precise: "In-memory circuit breakers track state on a per-instance basis without shared coordination across a scaled-out fleet"
- Added citations: cep21/circuit (Source 8) for implementation context, Microsoft Azure Architecture Center (Source 3) for concurrency considerations in multi-instance deployments

**Verification:**
- cep21/circuit source confirms per-instance state management
- Azure Circuit Breaker pattern docs note "Concurrency: A large number of concurrent instances of an application can access the same circuit breaker" implying multi-instance coordination challenges

**Status:** RESOLVED

---

## Revision 6
**Audit Issue:** MEDIUM — Universal prescriptive claim "Asynchronous non-critical flows must be decoupled using queues + idempotency so third-party downtime never aborts core transaction persistence" is overgeneralized and lacks citations

**Files Changed:**
- research/10-final-research.md (Synthesis 3)

**Action:**
- Softened claim from universal "must" to conditional "can be decoupled"
- Added reference to Queue-Based Load Leveling pattern (Source 12)
- Clarified that this is a pattern application, not a universal law

**Verification:**
- Azure Queue-Based Load Leveling pattern confirms queues decouple task intake from service processing, allowing application to continue when service is unavailable
- Pattern documents mention idempotency for at-least-once delivery

**Status:** RESOLVED

---

## Revision 7
**Audit Issue:** MEDIUM — Minor framing inconsistency: "Areas of Disagreement" vs "Implementation Variances" for rolling-window vs consecutive failures

**Files Changed:**
- research/05-report.md (Areas of Disagreement section)
- research/03-core-concepts.md (Key Configuration Parameters section)

**Action:**
- Clarified the Areas of Disagreement statement to explicitly reference sources (Netflix Hystrix, sony/gobreaker) and note this lab uses consecutive failures for determinism
- Added "(Illustrative Examples, Not Recommendations)" header to configuration parameter table to clearly label numeric ranges as examples
- Updated contradiction framing to align with research/04-contradictions.md documentation that already classifies this as an implementation variance rather than a disagreement

**Verification:**
- research/04-contradictions.md already documents this as implementation options rather than contradictions
- Configuration ranges (5-20, 10-60s, etc.) are now clearly labeled as illustrative examples from implementations (Hystrix, gobreaker, cep21/circuit)

**Status:** RESOLVED

---

## Revision 8
**Audit Issue:** LOW — Stale references to old broken URLs in additional research files

**Files Changed:**
- research/01-research-plan.md (line 62: cep21/circuitbreaker → cep21/circuit)
- research/03-core-concepts.md (lines 46, 61: cep21/circuitbreaker → cep21/circuit)

**Action:**
- Updated all references from the broken repo name `cep21/circuitbreaker` to the correct `cep21/circuit` across research files for consistency with the corrected Source 8

**Verification:**
- All references now point to the valid repository

**Status:** RESOLVED