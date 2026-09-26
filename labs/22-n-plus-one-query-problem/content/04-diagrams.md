# Diagrams

## N+1 Query Flow (Anti-Pattern)

```text
┌─────────────────────────────────────────────────────────────────┐
│                    N+1 Query Flow (4 queries)                    │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  1. App → DB: GetAllAuthors()                                    │
│       DB → App: Returns 3 authors (Alice, Bob, Charlie)        │
│                                                                  │
│  2. App → DB: GetPostsByAuthorID(1)  ← Loop iteration 1         │
│       DB → App: Returns posts by Alice                           │
│                                                                  │
│  3. App → DB: GetPostsByAuthorID(2)  ← Loop iteration 2         │
│       DB → App: Returns posts by Bob                           │
│                                                                  │
│  4. App → DB: GetPostsByAuthorID(3)  ← Loop iteration 3         │
│       DB → App: Returns posts by Charlie                        │
│                                                                  │
│  Total: 4 database round-trips                                   │
│  Pattern: 1 + N (N = number of parent records)                  │
└─────────────────────────────────────────────────────────────────┘
```

## Eager Loading Flow (Batched Solution)

```text
┌─────────────────────────────────────────────────────────────────┐
│                Eager Loading Flow (2 queries)                    │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  1. App → DB: GetAllAuthors()                                    │
│       DB → App: Returns 3 authors (Alice, Bob, Charlie)        │
│                                                                  │
│  2. App extracts IDs: [1, 2, 3]                                 │
│                                                                  │
│  3. App → DB: GetPostsByAuthorIDs([1, 2, 3])                     │
│       DB → App: Returns ALL related posts in one query         │
│                                                                  │
│  4. App: Groups posts by AuthorID in memory (map[int][]Post)    │
│                                                                  │
│  5. App: Builds AuthorWithPosts struct for each author          │
│                                                                  │
│  Total: 2 database round-trips                                    │
│  Query count is independent of N                                │
└─────────────────────────────────────────────────────────────────┘
```

## Architecture Overview

```text
┌──────────────────────────────────────────────────────────────────────┐
│                       Architecture Diagram                          │
├──────────────────────────────────────────────────────────────────────┤
│                                                                      │
│   ┌──────────┐       ┌──────────┐       ┌──────────┐                │
│   │   Demo   │──────▶│  Query   │──────▶│   Data   │                │
│   │  CLI     │       │Counter   │       │  Store   │                │
│   │          │       │Track     │       │(Mock DB) │                │
│   └──────────┘       └──────────┘       └────┬─────┘                │
│                                              │                       │
│                                              ▼                       │
│                                        ┌──────────────┐              │
│                                        │  Repository  │              │
│                                        │    Layer     │              │
│                                        └──────┬───────┘              │
│                                               │                      │
│                ┌──────────────────────────────┴────────────────────┐ │
│                │                                                   │ │
│   ┌────────────▼───────────────┐       ┌────────────────▼────────┐ │
│   │ GetAuthorsWithPostsNPlusOne│       │ GetAuthorsWithPostsEager│ │
│   │          (N+1)             │       │        (Eager)          │ │
│   │       4 queries            │       │      2 queries          │ │
│   └────────────────────────────┘       └─────────────────────────┘ │
│                                                                      │
└──────────────────────────────────────────────────────────────────────┘
```

## Query Comparison Matrix

```text
┌─────────────────────┬──────────────────┬──────────────────┐
│     Approach        │   Query Count  │ Time Complexity  │
├─────────────────────┼──────────────────┼──────────────────┤
│ N+1 (Naive)         │ 1 + N            │ O(N) queries     │
│ Eager Loading       │ 2                │ O(1) queries     │
└─────────────────────┴──────────────────┴──────────────────┘

Where N = number of parent records (authors in this case)
```

## Data Flow Comparison

```text
N+1 Approach (Iterative):
┌──────────┐     ┌──────────┐     ┌──────────┐
│ Get      │     │ Get      │     │ Get      │
│ Authors  │────▶│ Author 1 │────▶│ Author 2 │
│ (1 Q)    │     │ Posts    │     │ Posts    │
└──────────┘     │ (1 Q)    │     │ (1 Q)    │
                 └──────────┘     └──────────┘


Eager Loading Approach (Batched):
┌──────────┐     ┌──────────┐
│ Get      │     │ Get      │
│ Authors  │────▶│ All      │
│ (1 Q)    │     │ Posts    │
└──────────┘     │ (1 Q)    │
                 └────┬─────┘
                      ▼
               ┌─────────────┐
               │ Group In    │
               │ Memory      │
               └─────────────┘
```