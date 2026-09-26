## Evidence 1
Claim: The N+1 query problem occurs when a data access framework executes N additional SQL statements to fetch data that could have been retrieved in a primary query.
Evidence: "The N+1 query problem happens when the data access framework executes N additional SQL statements to fetch the same data that could have been retrieved when executing the primary SQL query."
Source: N+1 query problem with JPA and Hibernate
URL: https://vladmihalcea.com/n-plus-1-query-problem/
Confidence: HIGH
Corroborated By: Source 3 (Shopify Engineering) which identifies a similar pattern for GraphQL ("server makes 1 round trip... then makes N round trips").
Notes: Serves as the foundational definition of the issue.

## Evidence 2
Claim: N+1 issues are difficult to detect during development because individual queries execute quickly enough to bypass slow query logs, yet their aggregate volume severely impacts overall response times.
Evidence: "...unlike the slow query log that can help you find slow-running queries, the N+1 issue won’t be spotted because each individual additional query runs sufficiently fast to not trigger the slow query log. The problem is executing a large number of additional queries that, overall, take sufficient time to slow down response time."
Source: N+1 query problem with JPA and Hibernate
URL: https://vladmihalcea.com/n-plus-1-query-problem/
Confidence: HIGH
Corroborated By: General APM documentation emphasizing total request query tracking over slow query logging.
Notes: Explains the common disconnect between local testing (where dataset is small) and production degradation.

## Evidence 3
Claim: Eager loading solves the relational N+1 query problem by loading necessary related models concurrent with the parent query, vastly reducing the total query count.
Evidence: "Eloquent can 'eager load' relationships at the time you query the parent model. Eager loading alleviates the 'N + 1' query problem."
Source: Eloquent: Relationships | Laravel 11.x
URL: https://laravel.com/docs/11.x/eloquent-relationships#eager-loading
Confidence: HIGH
Corroborated By: Source 1 (Vlad Mihalcea), advising the use of `JOIN FETCH` (eager fetching) to consolidate data retrieval.
Notes: Highlights the primary technique used to mitigate N+1 within ORMs.

## Evidence 4
Claim: Unrestricted eager loading creates severe memory bloat by fetching excessive amounts of unneeded data.
Evidence: "Using FetchType.EAGER either implicitly or explicitly for your JPA associations is a bad idea because you are going to fetch way more data that you need."
Source: N+1 query problem with JPA and Hibernate
URL: https://vladmihalcea.com/n-plus-1-query-problem/
Confidence: HIGH
Corroborated By: Source 4 (Vlad Mihalcea) explicitly defining eager fetching without bounds as an anti-pattern.
Notes: Warns against treating eager loading as a silver bullet without considering data volume.

## Evidence 5
Claim: The N+1 problem also occurs over network API boundaries, where executing multiple external HTTP requests introduces massive latency due to network overhead.
Evidence: "The n+1 problem means that the server executes multiple unnecessary round trips to datastores for nested data... The computing expenditure of these extra round trips are massive when applied to large requests."
Source: Solving the N+1 Problem for GraphQL through Batching
URL: https://shopify.engineering/solving-the-n-1-problem-for-graphql-through-batching
Confidence: HIGH
Corroborated By: REST architecture constraints requiring batch API endpoints.
Notes: Demonstrates the N+1 problem is architectural, not strictly limited to SQL databases.

## Evidence 6
Claim: Network-based N+1 issues can be resolved using batch loaders which group individual data promises and load them collectively.
Evidence: "GraphQL Batch allows applications to define batch loaders that specify how to group and load similar data... GraphQL Batch iterates through the grouped loads, uses their corresponding batch loader to load all the promises together, and replaces the promises with the loaded result."
Source: Solving the N+1 Problem for GraphQL through Batching
URL: https://shopify.engineering/solving-the-n-1-problem-for-graphql-through-batching
Confidence: HIGH
Corroborated By: Established DataLoader patterns popularized by Facebook.
Notes: Explains the batching/dataloader pattern as the solution to network N+1.
