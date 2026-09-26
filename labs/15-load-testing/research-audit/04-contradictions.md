# Contradictions Audit

## Contradiction 1: Test Type Terminology Variations
Statement A: k6 states "no consensus even exists about the names of these test types" (surge, scale, stamina, limit testing).  
Location: `research/04-contradictions.md:3-9`, Source 1  
Statement B: Azure documentation adopts unified terms without mentioning industry variants.  
Location: `research/04-contradictions.md:7`, Source 11  
Type: SOURCE_CONFLICT  
Impact: MINOR  
Assessment: Non-conflicting nuance. Nomenclature in industry varies across tool vendors.

---

## Contradiction 2: Source 15 Index vs URL Mismatch
Statement A: Source 15 in `02-sources.md` is titled "A Collection of Best Practices for Production Services" with URL `https://sre.google/sre-book/service-best-practices/`.  
Location: `research/02-sources.md:145-148`  
Statement B: Evidence 2 and Finding 4 cite Source 15 as k6 Thresholds documentation with URL `https://grafana.com/docs/k6/latest/using-k6/thresholds/`.  
Location: `research/03-evidence.md:29`, `research/05-report.md:80, 99`  
Type: INTERNAL  
Impact: MEDIUM  
Assessment: Internal numbering slip in research files where k6 Thresholds was referenced as Source 15 instead of creating a distinct source entry.

---

## Contradiction 3: Cross-Lab Contamination (Spring DI)
Statement A: Source 16 lists "Introduction to the Spring IoC Container and Beans" and Contradiction 6 mentions "Service Locator vs DI Pattern".  
Location: `research/02-sources.md:155-162`, `research/04-contradictions.md:43-47`  
Statement B: The research plan and scope explicitly target load testing for the Booking Bengkel app.  
Location: `research/01-plan.md:3-15`  
Type: INTERNAL  
Impact: LOW  
Assessment: Stray residue from Lab 16. The researcher explicitly noted "Excluded from active evidence" and "Not applicable to load testing", so it does not corrupt conclusions.

---

## Contradiction 4: Breakpoint Testing in Elastic Cloud
Statement A: k6 advises turning off cloud elasticity during breakpoint testing to prevent infinite billing and masked saturation.  
Location: `research/04-contradictions.md:19-26`, Source 6  
Statement B: Azure presents breakpoint testing without elasticity warnings.  
Location: `research/04-contradictions.md:23`, Source 11  
Type: SOURCE_CONFLICT  
Impact: MINOR  
Assessment: k6 provides operational risk guidance, whereas Azure provides architectural definitions. Complementary, not contradictory.
