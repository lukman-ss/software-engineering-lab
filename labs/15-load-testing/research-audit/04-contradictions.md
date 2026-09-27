# Contradictions Audit: Lab 15 (Load Testing)

## Contradiction 1: Stress Test Percentage Increase (Fixed Rule vs Risk-Derived)
Statement A: Older tutorials and community discussions cite "50-100% above average" as default stress testing target.  
Location: `research/04-contradictions.md: Contradiction 2`  
Statement B: k6 official documentation explicitly states "there is no fixed percentage" and load must be determined by situational risk models.  
Location: `research/02-sources.md: Source 3`  
Type: SOURCE_CONFLICT  
Impact: Low. Tier 1 guidance (k6) is authoritative; the research correctly identifies the arbitrary percentage as an oversimplification.  
Assessment: RESOLVED. Research adopted the authoritative risk-derived stance.  

## Contradiction 2: Breakpoint Testing in Cloud Auto-Scaling Environments
Statement A: General load testing literature suggests running breakpoint testing until failure to find saturation capacity.  
Location: Azure Well-Architected Framework guidance (Source 11)  
Statement B: k6 documentation warns that breakpoint tests in elastic cloud environments simply scale infinitely and trigger massive cloud bills rather than finding application limits unless auto-scaling is disabled.  
Location: `research/02-sources.md: Source 6`  
Type: SOURCE_CONFLICT / OPERATIONAL_NUANCE  
Impact: Medium. Running breakpoint tests without pinning infrastructure capacity can invalidate tests and incur major cost.  
Assessment: RESOLVED. Documented as a critical operational nuance.  

## Contradiction 3: Test Environment (Controlled Production vs Isolated Staging)
Statement A: Microsoft Azure recommends running controlled performance tests in production during off-peak hours for true fidelity.  
Location: `research/02-sources.md: Source 11`  
Statement B: Common practice and lab warning advises avoiding live environments to prevent data corruption or user disruption, advocating for staging environments that mirror production.  
Location: `research/04-contradictions.md: Contradiction 5`  
Type: INTERNAL / METHODOLOGY_VARIATION  
Impact: Medium. Risk of testing in production without proper synthetic data segregation.  
Assessment: RESOLVED. Reconciled by establishing that staging-mirror is primary for destructive tests, with controlled production tests reserved for non-destructive canary load tests.  

## Contradiction 4: Residual Dependency Injection Artifact
Statement A: Source 16 ("Introduction to the Spring IoC Container and Beans") and Contradiction 6 ("Service Locator vs DI Pattern") are listed in load testing files.  
Location: `research/02-sources.md: Source 16`, `research/04-contradictions.md: Contradiction 6`  
Statement B: The lab topic is Load Testing, not Dependency Injection.  
Location: `research/01-plan.md`  
Type: INTERNAL  
Impact: Low (Clutter/noise). Both files explicitly note that it is excluded from active evidence, but retaining it indicates a copy-paste artifact from Lab 16.  
Assessment: NON-BLOCKING. Noted for cleanup.  
