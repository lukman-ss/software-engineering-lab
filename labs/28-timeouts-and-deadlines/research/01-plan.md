# Research Plan: Timeouts — Jangan Biarkan Satu Dependency Lambat Menahan Seluruh Sistem

## Research Topic
Timeouts in distributed systems: Understanding how to properly configure timeouts to prevent cascading failures while maintaining system reliability and correctness.

## Objective
To investigate best practices for implementing timeouts in distributed systems, including:
- Different types of timeouts (connect, read, request, statement, job)
- Timeout budgeting based on latency distributions
- Relationship between timeouts, retries, and circuit breakers
- Risks of improper timeout configuration
- Monitoring and observability for timeout-related issues

## Research Questions
1. What are the different types of timeouts in distributed systems and when should each be used?
2. How should timeout budgets be calculated based on service latency requirements and dependency latency distributions?
3. What are the dangers of setting timeouts too high or too low?
4. How do timeouts interact with retry mechanisms to prevent retry storms?
5. What are the correctness implications of timeouts (especially for idempotency)?
6. How should timeouts propagate through service chains (deadline propagation)?
7. What timeout mechanisms exist for databases, queues, and other non-HTTP dependencies?
8. What metrics should be monitored to detect timeout-related issues?

## Search Strategy
1. Search for official documentation from major tech companies on their timeout practices
2. Look for academic papers on distributed systems reliability and timeout strategies
3. Review industry best practices from reputable sources (Google SRE book, Amazon Builders' Library, etc.)
4. Examine open-source projects for timeout implementation patterns
5. Search for case studies of timeout-related outages

## Expected Primary Sources
- Google SRE Book (Site Reliability Engineering)
- Amazon Builders' Library articles
- Microsoft Azure Architecture Center
- CNCF papers on distributed systems
- Databases documentation (PostgreSQL, MySQL) on statement timeouts
- Message queue documentation (RabbitMQ, Apache Kafka) on job timeouts
- HTTP client library documentation (various languages)
- Academic papers on retry storms and cascade failures

## Risks / Unknowns
- Variability in timeout terminology across different technologies
- Lack of recent empirical data on timeout configuration effectiveness
- Proprietary implementations not publicly documented
- Evolution of best practices over time