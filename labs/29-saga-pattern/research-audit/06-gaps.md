# 06 - Research Gap Analysis: Saga Pattern Research

## Gap 1: Verbatim Primary Extraction of Garcia-Molina & Salem (1987)
Type: WEAK_SOURCE  
Severity: LOW  
Location: `research/03-evidence.md:16-18`, `research/06-open-questions.md:9-10`  
Problem: The 1987 original ACM paper is a scanned PDF with LZW compression that prevented automated text scraping of verbatim sentences.  
Required Revision: Secondary citations (ACM Digital Library record, Cornell CS711 syllabus, Microsoft, Temporal) corroborate the authorship, publication year, and core definition. Verbatim citation is documented as a known limitation in `06-open-questions.md`.  
Can Be Approved Without Fix: YES  

---

## Gap 2: Single-Source Taxonomy for Step Classification (Pivot / Retryable)
Type: SCOPE_ERROR (Potential Overgeneralization)  
Severity: MEDIUM  
Location: `research/03-evidence.md:111-126`, `research/05-report.md:64-76`  
Problem: The explicit division of transactions into compensable, pivot, and retryable steps is derived almost exclusively from Microsoft Azure Architecture Center documentation and is not universally phrased this way across all microservices literature.  
Required Revision: Research report already flags this with MEDIUM confidence and notes that it is Microsoft's formalization. When drafting downstream content, present this taxonomy as a useful structured framework rather than a universal standard.  
Can Be Approved Without Fix: YES  

---

## Gap 3: Six Isolation Countermeasures Attribution
Type: WEAK_SOURCE  
Severity: MEDIUM  
Location: `research/03-evidence.md:202-216`, `research/05-report.md:129-142`  
Problem: The complete list of 6 isolation countermeasures (semantic lock, commutative updates, pessimistic view, reread values, version files, risk-based concurrency) is enumerated in public documentation only by Microsoft, while Microservices.io references the Manning textbook.  
Required Revision: Report explicitly documents this dependency and assigned MEDIUM confidence. Future content can reference semantic locking and commutative updates as primary patterns while treating the full 6-item list as Microsoft/Richardson's extended catalog.  
Can Be Approved Without Fix: YES  

---

## Gap 4: Canonical Protocol for Failed Compensations
Type: MISSING_CASE  
Severity: LOW  
Location: `research/04-contradictions.md:46-49`, `research/06-open-questions.md:4-5`  
Problem: Distributed systems literature lacks a universal automatic protocol for "compensation of a failed compensation," relying instead on retries, dead-letter queues, and operational intervention.  
Required Revision: Accurately cataloged as an inherent limitation of Sagas in `06-open-questions.md`. No revision needed to research files; content and implementation should reflect the necessity of retry backoff and operator alerting.  
Can Be Approved Without Fix: YES  
