# Contradictions

## Contradiction 1: Test Type Naming Consensus

SOURCE A: k6 documentation explicitly states "no consensus even exists about the names of these test types" and notes that stress tests may also be called "rush-hour, surge, or scale tests", soak tests may be called "endurance, constant high load, or stamina tests", and breakpoint tests may be called "capacity, point load, or limit testing" (Source 1).

SOURCE B: Microsoft Azure Well-Architected Framework documentation uses consistent naming: "Load testing," "Stress testing," "Spike testing," "Endurance/soak testing" with minimal mention of alternative names (Source 11).

ASSESSMENT: The k6 documentation acknowledges that naming varies by community and organization. Azure uses simplified common names for clarity. Both are internally consistent; the difference is granularity of naming variations documented. No factual disagreement; k6 is more transparent about naming inconsistency in the industry.

## Contradiction 2: Stress Test Load Increase Percentage

SOURCE A: k6 documentation states "Some testers might have default targets for stress tests—say an increase upon average load by 50 or 100 percent—there's no fixed percentage" and "The load simulated in a Stress test depends on the stressful situations that the system may be subject to" (Source 3).

SOURCE B: Some industry blogs and older JMeter tutorials reference the "50-100% above average" as a general rule of thumb without explicitly citing k6.

ASSESSMENT: k6 (Tier 1) is authoritative and explicitly states there is no fixed rule for how much above average a stress test should go. The "50-100%" figure found in various blogs is an oversimplification. The correct position per Tier 1 sources: load for stress testing should be derived from the system's own risk profile (e.g., expected peak events like payday, rush hour, end of workweek), not an arbitrary percentage.

## Contradiction 3: Breakpoint Testing in Elastic Environments

SOURCE A: k6 documentation explicitly warns: "Avoid breakpoint tests in elastic cloud environments. The elastic environment may grow as the test moves further, finding only the limit of your cloud account bill. If this test runs on a cloud environment, turning off elasticity on all the affected components is strongly recommended." (Source 6).

SOURCE B: Microsoft Azure Performance Testing documentation lists breakpoint testing as a valid approach for finding "maximum capacity" and "failure modes" (Source 11) but does not include the specific warning about elastic cloud environments.

ASSESSMENT: Both sources agree breakpoint testing finds system limits. k6 provides a more nuanced implementation warning about elasticity that Azure does not mention. This is likely because Azure's guidance is more general while k6's is tool-specific. k6's warning reflects a real operational risk: auto-scaling can mask the true system limit by simply adding resources, potentially leading to unbounded costs. No factual disagreement; k6 provides additional operational nuance.

## Contradiction 4: Threshold Evaluation Frequency in k6 Cloud

SOURCE A: k6 documentation notes: "When k6 runs in the cloud, thresholds are evaluated every 60 seconds. Therefore, the abortOnFail feature may be delayed by up to 60 seconds." (Source 15).

SOURCE B: The k6 local (CLI) execution mode evaluates thresholds more frequently.

ASSESSMENT: This difference is by design based on the execution environment. The cloud evaluation interval creates a potential gap for abortOnFail thresholds. This is not a contradiction between sources but a platform-dependent behavior. Engineers should be aware that cloud-based k6 tests may continue running up to 60 seconds past threshold failure.

## Contradiction 5: Performance Testing Environment Isolation

SOURCE A: Microsoft Azure documentation states for performance testing: "Run controlled production testing. Schedule tests during off-peak hours" and recommends production environments for the most realistic results (Source 11).

SOURCE B: The original lab specification advises: "Load test sebaiknya dilakukan pada lingkungan yang mendekati production" (load testing should be conducted in an environment close to production) and warns against using laptops as production stand-ins (original lab content).

ASSESSMENT: Azure's guidance is more nuanced — it distinguishes between staging/prod-like environments for different test types and includes the option of controlled production testing. The lab's advice to avoid laptop testing aligns with Azure's "mirror your production environment" recommendation. Both agree that non-production environments should mirror production as closely as possible. The apparent tension (test in staging vs. test in production) is resolved by Azure's framework: use staging for most tests, production only with control measures.

## Contradiction 6: Service Locator vs. DI Pattern (from related lab context)

SOURCES: These sources are from the Dependency Injection lab (labs/16), not directly load testing, but referenced in this lab's prior research. Martin Fowler notes that Service Locator "hides class dependencies" (Source from DI lab).

ASSESSMENT: Not applicable to load testing research. Excluded from scope.

## Summary

No material contradictions discovered in the load testing domain itself. The differences between sources are primarily:
1. Granularity of terminology documentation (k6 acknowledges more naming variants)
2. Operational nuances specific to tool implementations (k6's cloud threshold timing, elasticity warning)
3. Scope of coverage (Azure provides general guidance, k6 provides tool-specific details)

All Tier 1 sources agree on fundamental principles: test type definitions, metric importance, bottleneck identification methodology, and common pitfalls.