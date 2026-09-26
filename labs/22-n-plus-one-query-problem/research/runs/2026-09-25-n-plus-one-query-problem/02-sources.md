## Source 1
Title: N+1 query problem with JPA and Hibernate
Publisher: Vlad Mihalcea
URL: https://vladmihalcea.com/n-plus-1-query-problem/
Published: March 17, 2020
Accessed: September 25, 2026
Source Tier: Tier 2
Relevance: Directly defines the N+1 problem, explains the lazy-loading trigger, demonstrates performance impact despite fast individual queries, and cautions against the unrestricted use of eager loading.

## Source 2
Title: Eloquent: Relationships | Laravel 11.x
Publisher: Laravel
URL: https://laravel.com/docs/11.x/eloquent-relationships#eager-loading
Published: Unknown
Accessed: September 25, 2026
Source Tier: Tier 1
Relevance: Official framework documentation detailing how eager loading acts as a solution to alleviate the N+1 problem by batching relation requests.

## Source 3
Title: Solving the N+1 Problem for GraphQL through Batching
Publisher: Shopify Engineering
URL: https://shopify.engineering/solving-the-n-1-problem-for-graphql-through-batching
Published: April 24, 2018
Accessed: September 25, 2026
Source Tier: Tier 1
Relevance: Examines how the N+1 problem extends beyond databases into network API interactions (Network N+1) and explains how to solve it using dataloaders/batching.

## Source 4
Title: Eager fetching is a code smell
Publisher: Vlad Mihalcea
URL: https://vladmihalcea.com/eager-fetching-is-a-code-smell/
Published: Unknown
Accessed: September 25, 2026
Source Tier: Tier 2
Relevance: Investigates the downside of naive N+1 fixes, primarily focusing on how over-fetching through eager loading bloats memory and strains data transfer.
